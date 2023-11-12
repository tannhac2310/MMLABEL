package subscriptions

import (
	"context"
	"fmt"
	"github.com/nats-io/stan.go"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
)

func (p *EventMQTTSubscription) handleZaloEvent(msg *stan.Msg) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ctx = ctxzap.ToContext(ctx, p.logger)
	ctx = cockroach.ContextWithDB(ctx, p.db)

	err2 := msg.Ack()
	if err2 != nil {
		return
	}
	fmt.Println("==done==")
}
