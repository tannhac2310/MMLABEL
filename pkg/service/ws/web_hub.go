package ws

import (
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
)

const (
	BroadcastQueueSize = 4096
	DeadlockTicker     = 15 * time.Second                // check every 15 seconds
	DeadlockWarn       = (BroadcastQueueSize * 99) / 100 // number of buffered messages before printing stack trace

	SessionCacheSize    = 35000
	ConnMemberCacheTime = 30 * time.Minute
)

type Broadcast struct {
	UserID string
	Data   *WebSocketEvent
}

type Hub struct {
	connectionCount int64
	logger          *zap.Logger
	app             *wsApp
	connectionIndex int
	register        chan *WebConn
	unregister      chan *WebConn
	broadcast       chan *Broadcast
	stop            chan struct{}
	didStop         chan struct{}
	goroutineID     int
	ExplicitStop    bool
}

func (h *Hub) Register(webConn *WebConn) {
	select {
	case h.register <- webConn:
	case <-h.didStop:
	}
}

func (h *Hub) Unregister(webConn *WebConn) {
	select {
	case h.unregister <- webConn:
	case <-h.stop:
	}
}

func (h *Hub) Broadcast(message *Broadcast) {
	if h != nil && h.broadcast != nil && message != nil {
		select {
		case h.broadcast <- message:
		case <-h.didStop:
		}
	}
}

func getGoroutineID() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		id = -1
	}
	return id
}

func (h *Hub) Stop() {
	close(h.stop)
	<-h.didStop
}

func (h *Hub) Start() {
	var doStart func()
	var doRecoverableStart func()
	var doRecover func()

	doStart = func() {
		h.goroutineID = getGoroutineID()
		h.logger.Debug("Hub for index is starting with goroutine", zap.Int("index", h.connectionIndex), zap.Int("goroutine", h.goroutineID))

		connections := newHubConnectionIndex()

		for {
			select {
			case webCon := <-h.register:
				connections.Add(webCon)
				atomic.StoreInt64(&h.connectionCount, int64(len(connections.All())))

				h.app.SetStatusOnline(webCon.UserID, true)
			case webCon := <-h.unregister:
				connections.Remove(webCon)
				atomic.StoreInt64(&h.connectionCount, int64(len(connections.All())))

				h.app.SetStatusOnline(webCon.UserID, false)
			case msg := <-h.broadcast:
				candidates := connections.All()
				if msg.UserID != "" {
					candidates = connections.ForUser(msg.UserID)
				}

				data := msg.Data.PrecomputeJSON()
				for _, webCon := range candidates {
					select {
					case webCon.Send <- data:
					default:
						webCon.logger.Error("webhub.broadcast: cannot send, closing websocket for user")
						close(webCon.Send)
						connections.Remove(webCon)
						atomic.StoreInt64(&h.connectionCount, int64(len(connections.All())))

						h.app.SetStatusOnline(webCon.UserID, false)
					}
				}
			case <-h.stop:
				userIDs := make(map[string]bool)

				for _, webCon := range connections.All() {
					userIDs[webCon.UserID] = true
					webCon.Close()
				}

				for userID := range userIDs {
					h.app.SetStatusOnline(userID, false)
				}

				h.ExplicitStop = true
				close(h.didStop)

				return
			}
		}
	}

	doRecoverableStart = func() {
		defer doRecover()
		doStart()
	}

	doRecover = func() {
		if !h.ExplicitStop {
			if r := recover(); r != nil {
				h.logger.Error("Recovering from Hub panic.", zap.Any("panic", r))
			} else {
				h.logger.Error("Webhub stopped unexpectedly. Recovering.")
			}

			h.logger.Error(string(debug.Stack()))

			go doRecoverableStart()
		}
	}

	go doRecoverableStart()
}

type hubConnectionIndexIndexes struct {
	connections         int
	connectionsByUserID int
}

// hubConnectionIndex provides fast addition, removal, and iteration of web connections.
type hubConnectionIndex struct {
	connections         []*WebConn
	connectionsByUserID map[string][]*WebConn
	connectionIndexes   map[*WebConn]*hubConnectionIndexIndexes
}

func newHubConnectionIndex() *hubConnectionIndex {
	return &hubConnectionIndex{
		connections:         make([]*WebConn, 0, SessionCacheSize),
		connectionsByUserID: make(map[string][]*WebConn),
		connectionIndexes:   make(map[*WebConn]*hubConnectionIndexIndexes),
	}
}

func (i *hubConnectionIndex) Add(wc *WebConn) {
	i.connections = append(i.connections, wc)
	i.connectionsByUserID[wc.UserID] = append(i.connectionsByUserID[wc.UserID], wc)
	i.connectionIndexes[wc] = &hubConnectionIndexIndexes{
		connections:         len(i.connections) - 1,
		connectionsByUserID: len(i.connectionsByUserID[wc.UserID]) - 1,
	}
}

func (i *hubConnectionIndex) Remove(wc *WebConn) {
	indexes, ok := i.connectionIndexes[wc]
	if !ok {
		return
	}

	last := i.connections[len(i.connections)-1]
	i.connections[indexes.connections] = last
	i.connections = i.connections[:len(i.connections)-1]
	i.connectionIndexes[last].connections = indexes.connections

	userConnections := i.connectionsByUserID[wc.UserID]
	last = userConnections[len(userConnections)-1]
	userConnections[indexes.connectionsByUserID] = last
	i.connectionsByUserID[wc.UserID] = userConnections[:len(userConnections)-1]
	i.connectionIndexes[last].connectionsByUserID = indexes.connectionsByUserID

	delete(i.connectionIndexes, wc)
}

func (i *hubConnectionIndex) ForUser(id string) []*WebConn {
	return i.connectionsByUserID[id]
}

func (i *hubConnectionIndex) All() []*WebConn {
	return i.connections
}
