package subscriptions

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nats-io/stan.go"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/constants"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/nats"
)

const (
	MaxRetryPushNotify = 10
)

type NotificationSubscription struct {
	db         cockroach.Ext
	busFactory nats.BusFactory
	logger     *zap.Logger
}

func NewNotificationSubscription(
	db cockroach.Ext,
	busFactory nats.BusFactory,
	logger *zap.Logger,
) *NotificationSubscription {
	return &NotificationSubscription{
		db:         db,
		busFactory: busFactory,
		logger:     logger,
	}
}

func (n *NotificationSubscription) Subscribe() ([]stan.Subscription, error) {
	bus := n.busFactory.GetConn()

	subPushNotification, err := bus.QueueSubscribe(
		constants.SubjectPushNotification,
		nats.ToQueueName(constants.SubjectPushNotification),
		n.handleNotificationEvent,
		stan.DurableName("d_"+constants.SubjectPushNotification),
		stan.SetManualAckMode(),
	)
	if err != nil {
		return nil, fmt.Errorf("bus.QueueSubscribe: %w", err)
	}

	return []stan.Subscription{
		subPushNotification,
	}, nil
}

func (n *NotificationSubscription) handleNotificationEvent(msg *stan.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ctx = cockroach.ContextWithDB(ctx, n.db)

	if msg.RedeliveryCount >= MaxRetryPushNotify {
		defer nats.LogError(msg.Ack(), "msg.Ack", n.logger)
	}

	ctx = ctxzap.ToContext(ctx, n.logger)

	req := &nats.PushNotifyEvent{}
	err := json.Unmarshal(msg.Data, req)
	if err != nil {
		n.logger.Error("err json.Unmarshal", zap.Error(err))
		nats.LogError(msg.Ack(), "msg.Ack", n.logger)
		return
	}

	if err != nil {
		n.logger.Error("err n.notificationService.PushNotification", zap.Error(err))
		nats.LogError(msg.Ack(), "msg.Ack", n.logger)
		return
	}

	nats.LogError(msg.Ack(), "msg.Ack", n.logger)
}
