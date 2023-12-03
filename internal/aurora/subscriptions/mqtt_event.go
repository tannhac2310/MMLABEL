package subscriptions

import (
	"context"
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/service/ws"
	"strings"
	"time"
)

type EventMQTTSubscription struct {
	db cockroach.Ext
	//busFactory nats.BusFactory
	productionOrderStageDeviceRepo repository.ProductionOrderStageDeviceRepo
	logger                         *zap.Logger
	wsService                      ws.WebSocketService
}

func NewMQTTSubscription(
	db cockroach.Ext,
	productionOrderStageDeviceRepo repository.ProductionOrderStageDeviceRepo,
	logger *zap.Logger,
	wsService ws.WebSocketService,
) *EventMQTTSubscription {
	return &EventMQTTSubscription{
		db:                             db,
		productionOrderStageDeviceRepo: productionOrderStageDeviceRepo,
		logger:                         logger,
		wsService:                      wsService,
	}
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("==============================================>>>  Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("==============================================>>>>>Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

func (p *EventMQTTSubscription) Subscribe() error {
	var broker = "146.196.65.17"
	var port = 31883
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	//opts.SetClientID("go_mqtt_client")
	opts.SetUsername("user1")
	opts.SetPassword("123")
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	topic := "toppic"
	// Subscribe to the desired topic
	if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Printf("Error subscribing to topic '%s': %v\n", topic, token.Error())

		return nil
	}
	client.AddRoute("toppic1", func(client mqtt.Client, message mqtt.Message) {
		fmt.Println("============>>> Received message for topic 1: ", message.MessageID())
		message.Ack()
	})

	client.AddRoute("toppic", func(client mqtt.Client, message mqtt.Message) {
		fmt.Println("============>>> Received message for topic ====", string(message.Payload()))
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		ctx = ctxzap.ToContext(ctx, p.logger)
		ctx = cockroach.ContextWithDB(ctx, p.db)

		// parse message.Payload() to data struct
		var iotData MyStruct
		err := json.Unmarshal(message.Payload(), &iotData)
		if err != nil {
			fmt.Println("Error parsing JSON:", err)
			return
		}

		fmt.Println(iotData.D, iotData.Ts)
		for _, item := range iotData.D {
			// find device in production order stage device
			deviceID := strings.Replace(item.Tag, ":Counter", "", -1)
			orderStageDevices, err := p.productionOrderStageDeviceRepo.Search(ctx, &repository.SearchProductionOrderStageDevicesOpts{
				DeviceID:                   deviceID,
				ProductionOrderStageStatus: enum.ProductionOrderStageStatusProductionStart,
				Limit:                      1,
				Offset:                     0,
			})

			if err != nil {
				fmt.Println("============>>> Received message for topic ====", err)
			}
			activeStageID := ""
			if len(orderStageDevices) == 1 && int64(item.Value) > 0 && int64(item.Value) > orderStageDevices[0].Quantity {
				// update first device
				device := orderStageDevices[0]
				activeStageID = device.ProductionOrderStageID
				_ = p.productionOrderStageDeviceRepo.Update(ctx, &model.ProductionOrderStageDevice{
					ID:                     device.ID,
					ProductionOrderStageID: device.ProductionOrderStageID,
					DeviceID:               device.DeviceID,
					Quantity:               int64(item.Value),
					ProcessStatus:          device.ProcessStatus,
					Status:                 device.Status,
					Responsible:            device.Responsible,
					Settings:               device.Settings,
					Note:                   device.Note,
					CreatedAt:              device.CreatedAt,
					UpdatedAt:              device.UpdatedAt,
				})
			}

			now := time.Now()
			date := now.Format("2006-01-02")
			// insert event log
			_ = p.productionOrderStageDeviceRepo.InsertEventLog(ctx, &model.EventLog{
				ID:        time.Now().UnixNano(),
				DeviceID:  deviceID,
				StageID:   cockroach.String(activeStageID),
				Quantity:  item.Value,
				Msg:       cockroach.String(string(message.Payload())),
				Date:      cockroach.String(date),
				CreatedAt: now,
			})
		}

		message.Ack()
	})
	fmt.Printf("Subscribed to topic '%s'\n", topic)
	return nil
}

// create struct to parse {"d":[{"tag":"B_PR04:CounterTag","value":38.00}],"ts":"2023-11-17T04:13:55+0000"}

type MyStruct struct {
	D  []DataItem `json:"d"`
	Ts string     `json:"ts"`
}

type DataItem struct {
	Tag   string  `json:"tag"`
	Value float64 `json:"value"`
}
