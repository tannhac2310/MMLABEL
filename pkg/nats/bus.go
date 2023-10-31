package nats

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/nats-io/stan.go"
	"go.uber.org/zap"
)

var Bus BusFactory

// BusFactory implemented by busFactoryImpl
type BusFactory interface {
	GetConn() stan.Conn
	RegisterCallbackConnected(funcs []Callback)
}

type busFactoryImpl struct {
	sync.RWMutex
	natsClusterID      string
	hostName           string
	busConn            stan.Conn
	logger             *zap.Logger
	afterConnectedFunc []Callback
}

type Callback func() error

// NewBusFactory creates new BusFactory
func NewBusFactory(zapLogger *zap.Logger, natsClusterID, natsAddr, hostName string) (BusFactory, error) {
	if natsClusterID == "" {
		return nil, errors.New("missing natsClusterID")
	}

	if natsAddr == "" {
		return nil, errors.New("missing natsAddr")
	}

	if hostName == "" {
		return nil, errors.New("missing hostName")
	}

	b := &busFactoryImpl{
		natsClusterID: natsClusterID,
		hostName:      hostName,
		busConn:       nil,
		logger:        zapLogger,
	}

	urls, err := URLs(natsAddr)
	if err != nil {
		return nil, err
	}

	var sopts []stan.Option
	sopts = append(sopts,
		stan.NatsURL(strings.Join(urls, ",")),
		stan.Pings(3, 5),
		stan.PubAckWait(3*time.Second),
		stan.SetConnectionLostHandler(func(_ stan.Conn, e error) {
			b.logger.Error("lost connection, will attempt to recreate it", zap.Error(e))
			b.connect(sopts...)
		}))

	if err != nil {
		return nil, err
	}

	b.connect(sopts...)

	if Bus == nil {
		Bus = b
	}

	return b, nil
}

func (rcv *busFactoryImpl) GetConn() stan.Conn {
	rcv.RLock()
	sc := rcv.busConn
	rcv.RUnlock()
	return sc
}

func (rcv *busFactoryImpl) connect(opts ...stan.Option) {
	sc, err := stan.Connect(rcv.natsClusterID, rcv.hostName, opts...)
	count := 0
	for err != nil {
		count++
		rcv.logger.Info("connecting...", zap.Int("attempt", count))

		sc, err = stan.Connect(rcv.natsClusterID, rcv.hostName, opts...)
		if err != nil {
			rcv.logger.Error("err connect", zap.Error(err), zap.Int("attempt", count))
		}

		time.Sleep(2 * time.Second)
	}

	rcv.Lock()
	rcv.busConn = sc
	rcv.Unlock()

	for _, f := range rcv.afterConnectedFunc {
		err := f()
		if err != nil {
			rcv.logger.Error("err run callback after connected", zap.Error(err))
		}
	}
}

func (rcv *busFactoryImpl) RegisterCallbackConnected(funcs []Callback) {
	rcv.afterConnectedFunc = funcs
}

func URLs(natsAddr string) ([]string, error) {
	natsIPs, err := net.LookupIP(natsAddr)
	if err != nil {
		return nil, err
	}

	natsURLs := make([]string, 0, len(natsIPs))
	for _, ip := range natsIPs {
		natsURLs = append(natsURLs, fmt.Sprintf("nats://%s:4222", ip.String()))
	}

	return natsURLs, nil
}

//TODO: store msg can not push
func handlePushMsgFail(ctx context.Context, msg interface{}, err error) error {
	return err
}

func PublishEvt(ctx context.Context, busFactory BusFactory, event string, msg interface{}) error {
	bus := busFactory.GetConn()
	var msgID string
	data, _ := json.Marshal(msg)

	msgID, err := bus.PublishAsync(event, data, func(s string, err error) {
		if err != nil {
			ctxzap.Extract(ctx).Error("PublishAsync Ack failed", zap.String("msg-id", s), zap.String("callback-error", err.Error()))
		}
	})
	if err != nil {
		return handlePushMsgFail(ctx, msg, fmt.Errorf("PublishAsync rcv.BusFactory.PublishAsync failed, msgID: %s, %w", msgID, err))
	}

	return err
}

func LogError(err error, msg string, logger *zap.Logger) {
	if err != nil {
		logger.Error(msg, zap.Error(err))
	}
}

func ToQueueName(channelName string) string {
	return "queue_" + channelName
}

func ToDurableName(channelName string) string {
	return "durable_" + channelName
}
