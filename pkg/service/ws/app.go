package ws

import (
	"context"
	"fmt"
	"hash/fnv"
	"runtime"
	"strings"
	"sync/atomic"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	UserKey     = "user:"
	UserKeyNode = "node:"
)

type WebSocketService interface {
	NewWebConn(ws *websocket.Conn, userID string, logger *zap.Logger) *WebConn
	HubRegister(webConn *WebConn)
	Publish(userIDs []string, message *WebSocketEvent, fallbackOfflineHandler func(userIDs []string))
	HubStop()
}

type wsApp struct {
	hostName string
	logger   *zap.Logger
	redisDB  redis.Cmdable

	handlers map[string]handlerFunc

	hubs                        []*Hub
	hubsStopCheckingForDeadlock chan bool
}

func NewApp(hostName string, logger *zap.Logger, redisDB redis.Cmdable) WebSocketService {
	a := &wsApp{
		hostName: hostName,
		logger:   logger,
		redisDB:  redisDB,
	}

	a.hubStart()
	a.initWebsocketRouter()

	return a
}

func (a *wsApp) SetStatusOnline(userID string, isOnline bool) {
	// store to redis
	userKey := UserKey + userID
	userKeyNode := UserKeyNode + a.hostName

	if isOnline {
		err := a.redisDB.HIncrBy(userKey, userKeyNode, 1).Err()
		if err != nil {
			a.logger.Warn("a.RedisDB.HSet", zap.Error(err))
		}

		err = a.redisDB.Expire(userKey, ConnMemberCacheTime).Err()
		if err != nil {
			a.logger.Warn("a.RedisDB.Expire", zap.Error(err))
		}

		return
	}

	totalSession := 0
	err := a.redisDB.HGet(userKey, userKeyNode).Scan(&totalSession)
	if err != nil {
		a.logger.Warn("a.redisDB.HSet", zap.Error(err))
		return
	}

	if totalSession == 1 {
		err = a.redisDB.HDel(userKey, userKeyNode).Err()
		if err != nil {
			a.logger.Warn("a.redisDB.HDel", zap.Error(err))
		}

		return
	}

	err = a.redisDB.HIncrBy(userKey, userKeyNode, -1).Err()
	if err != nil {
		a.logger.Warn("a.RedisDB.HSet", zap.Error(err))
	}
}

func (a *wsApp) newWebHub() *Hub {
	return &Hub{
		connectionCount: 0,
		logger:          a.logger,
		app:             a,
		connectionIndex: 0,
		register:        make(chan *WebConn, 1),
		unregister:      make(chan *WebConn, 1),
		broadcast:       make(chan *Broadcast, BroadcastQueueSize),
		stop:            make(chan struct{}),
		didStop:         make(chan struct{}),
		goroutineID:     0,
		ExplicitStop:    false,
	}
}

func (a *wsApp) TotalWebsocketConnections() int {
	count := int64(0)
	for _, hub := range a.hubs {
		count += atomic.LoadInt64(&hub.connectionCount)
	}

	return int(count)
}

func (a *wsApp) hubStart() {
	// Total number of hubs is twice the number of CPUs.
	numberOfHubs := runtime.NumCPU() * 2
	a.logger.Info("Starting websocket hubs", zap.Int("number_of_hubs", numberOfHubs))

	a.hubs = make([]*Hub, numberOfHubs)
	a.hubsStopCheckingForDeadlock = make(chan bool, 1)

	for i := 0; i < len(a.hubs); i++ {
		a.hubs[i] = a.newWebHub()
		a.hubs[i].connectionIndex = i
		a.hubs[i].Start()
	}

	go func() {
		ticker := time.NewTicker(DeadlockTicker)

		defer func() {
			ticker.Stop()
		}()

		for {
			select {
			case <-ticker.C:
				for _, hub := range a.hubs {
					if len(hub.broadcast) >= DeadlockWarn {
						a.logger.Error(
							"Hub processing might be deadlock with events in the buffer",
							zap.Int("hub", hub.connectionIndex),
							zap.Int("goroutine", hub.goroutineID),
							zap.Int("events", len(hub.broadcast)),
						)
						buf := make([]byte, 1<<16)
						runtime.Stack(buf, true)
						output := fmt.Sprintf("%s", buf)
						splits := strings.Split(output, "goroutine ")

						for _, part := range splits {
							if strings.Contains(part, fmt.Sprintf("%v", hub.goroutineID)) {
								a.logger.Error("Trace for possible deadlock goroutine", zap.String("trace", part))
							}
						}
					}
				}

			case <-a.hubsStopCheckingForDeadlock:
				return
			}
		}
	}()
}

func (a *wsApp) HubStop() {
	a.logger.Info("stopping websocket hub connections")

	select {
	case a.hubsStopCheckingForDeadlock <- true:
	default:
		a.logger.Warn("We appear to have already sent the stop checking for deadlocks command")
	}

	for _, hub := range a.hubs {
		hub.Stop()
	}

	a.hubs = []*Hub{}
}

func (a *wsApp) GetHubForUserID(userID string) *Hub {
	if len(a.hubs) == 0 {
		return nil
	}

	hash := fnv.New32a()
	_, _ = hash.Write([]byte(userID))
	index := hash.Sum32() % uint32(len(a.hubs))
	return a.hubs[index]
}

func (a *wsApp) HubRegister(webConn *WebConn) {
	hub := a.GetHubForUserID(webConn.UserID)
	if hub != nil {
		hub.Register(webConn)
	}
}

func (a *wsApp) HubUnregister(webConn *WebConn) {
	hub := a.GetHubForUserID(webConn.UserID)
	if hub != nil {
		hub.Unregister(webConn)
	}
}

func (a *wsApp) Publish(userIDs []string, message *WebSocketEvent, fallbackOfflineHandler func(userIDs []string)) {
	a.PublishSkipClusterSend(userIDs, message)

	// handle send to cluster
	var (
		offlineUserIDs             []string
		onlineInAnotherNodeUserIDs []string
		mapNode                    = make(map[string]bool)
	)
	checkOnline := func(userID string) {
		userKey := UserKey + userID

		// check user is online
		onlineNodes := a.redisDB.HKeys(userKey).Val()
		if len(onlineNodes) == 0 {
			offlineUserIDs = append(offlineUserIDs, userID)
			return
		}

		for _, node := range onlineNodes {
			// check key is `node:` prefix
			if strings.Contains(node, UserKeyNode) {
				nodeName := node[len(UserKeyNode):]
				if nodeName == a.hostName {
					continue
				}

				mapNode[nodeName] = true
				onlineInAnotherNodeUserIDs = append(onlineInAnotherNodeUserIDs, userID)
			}
		}
	}

	for _, userID := range userIDs {
		checkOnline(userID)
	}

	if len(offlineUserIDs) > 0 && fallbackOfflineHandler != nil {
		fallbackOfflineHandler(offlineUserIDs)
	}

	if len(onlineInAnotherNodeUserIDs) == 0 {
		return
	}

	// TODO: send to cluster for user online in another node
	// need to finalize which message broker we will use
}

func (a *wsApp) PublishSkipClusterSend(userIDs []string, message *WebSocketEvent) {
	for _, userID := range userIDs {
		hub := a.GetHubForUserID(userID)
		if hub != nil {
			hub.Broadcast(&Broadcast{
				UserID: userID,
				Data:   message,
			})
		}
	}
}

func (a *wsApp) NewWebConn(ws *websocket.Conn, userID string, logger *zap.Logger) *WebConn {
	ctx, cancel := context.WithTimeout(context.Background(), ConnMemberCacheTime)

	wc := &WebConn{
		ctx:          ctx,
		cancelFunc:   cancel,
		App:          a,
		Send:         make(chan WebSocketMessage, SendQueueSize),
		WebSocket:    ws,
		UserID:       userID,
		logger:       logger,
		endWritePump: make(chan struct{}),
		pumpFinished: make(chan struct{}),
	}

	return wc
}
