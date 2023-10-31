package subscriptions

import (
	"fmt"
	"github.com/nats-io/stan.go"

	"go.uber.org/zap"

	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/service/message"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/constants"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/nats"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/service/ws"
)

type ZaloEventSubscription struct {
	db             cockroach.Ext
	busFactory     nats.BusFactory
	logger         *zap.Logger
	wsService      ws.WebSocketService
	messageService message.Service
}

func NewZaloSubscription(
	db cockroach.Ext,
	busFactory nats.BusFactory,
	logger *zap.Logger,
	wsService ws.WebSocketService,
	messageService message.Service,
) *ZaloEventSubscription {
	return &ZaloEventSubscription{
		db:             db,
		busFactory:     busFactory,
		logger:         logger,
		wsService:      wsService,
		messageService: messageService,
	}
}

func (p *ZaloEventSubscription) Subscribe() ([]stan.Subscription, error) {
	bus := p.busFactory.GetConn()

	subZaloEvent, err := bus.QueueSubscribe(
		constants.SubjectZaloEvent,
		nats.ToQueueName(constants.SubjectZaloEvent),
		p.handleZaloEvent,
		stan.DurableName("d_"+constants.SubjectZaloEvent),
		stan.SetManualAckMode(),
	)
	if err != nil {
		return nil, fmt.Errorf("bus.QueueSubscribe: %w", err)
	}

	return []stan.Subscription{
		subZaloEvent,
	}, nil
}
