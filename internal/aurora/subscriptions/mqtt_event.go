package subscriptions

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/idutil"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/model"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/repository"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/enum"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/service/ws"
)

type EventMQTTSubscription struct {
	db cockroach.Ext
	//busFactory nats.BusFactory
	productionOrderStageDeviceRepo repository.ProductionOrderStageDeviceRepo
	deviceWorkingHistoryRepo       repository.DeviceWorkingHistoryRepo
	logger                         *zap.Logger
	wsService                      ws.WebSocketService
}

func NewMQTTSubscription(
	db cockroach.Ext,
	productionOrderStageDeviceRepo repository.ProductionOrderStageDeviceRepo,
	deviceWorkingHistoryRepo repository.DeviceWorkingHistoryRepo,
	logger *zap.Logger,
	wsService ws.WebSocketService,
) *EventMQTTSubscription {
	return &EventMQTTSubscription{
		db:                             db,
		productionOrderStageDeviceRepo: productionOrderStageDeviceRepo,
		deviceWorkingHistoryRepo:       deviceWorkingHistoryRepo,
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
	var broker = "146.196.65.9"
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

		// data is {"d":[{"tag":"B_PR04:Counter","value":413.00},{"tag":"B_PR03:Counter","value":391.00},{"tag":"B_PR03:TG_in","value":1.00},{"tag":"B_PR03:ON_OFF","value":1},{"tag":"B_PR04:TG_in","value":11.00},{"tag":"B_PR04:ON_OFF","value":1},{"tag":"B_PR03:SL_1Ngay","value":0.00},{"tag":"B_PR03:Counter_Day","value":393.00},{"tag":"B_PR04:SL_1Ngay","value":0.00},{"tag":"B_PR04:Counter_Day","value":431.00}],"ts":"2023-12-08T07:13:34+0000"}
		now := time.Now()
		dateStr := now.Format("2006-01-02")
		fmt.Println(iotData.D, iotData.Ts)
		for _, item := range iotData.D {
			// find device in production order stage device
			s := strings.Split(item.Tag, ":")
			if len(s) != 2 {
				continue
			}
			deviceID := s[0]
			action := s[1]

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
			if len(orderStageDevices) == 1 {
				device := orderStageDevices[0]
				activeStageID = device.ProductionOrderStageID
				if device.ProcessStatus == enum.ProductionOrderStageDeviceStatusStart {
					if int64(item.Value) > 0 && int64(item.Value) > device.Quantity {
						// update first device
						_ = p.productionOrderStageDeviceRepo.Update(ctx, &model.ProductionOrderStageDevice{
							ID:                     device.ID,
							ProductionOrderStageID: device.ProductionOrderStageID,
							DeviceID:               device.DeviceID,
							Quantity:               int64(item.Value),
							ProcessStatus:          device.ProcessStatus,
							Status:                 device.Status,
							Settings:               device.Settings,
							Note:                   device.Note,
							CreatedAt:              device.CreatedAt,
							UpdatedAt:              device.UpdatedAt,
							Responsible:            device.Responsible,
							EstimatedCompleteAt:    device.EstimatedCompleteAt,
							AssignedQuantity:       device.AssignedQuantity,
						})
					}

					// upsert to device_working_history
					if action == "SL_in_Ngay" {
						counterByDay := item.Value
						p.logger.Info("SL_in_Ngay", zap.String("deviceID", deviceID), zap.Int("counterByDay", int(counterByDay)))

						deviceWorkingHistories, err := p.deviceWorkingHistoryRepo.Search(ctx, &repository.SearchDeviceWorkingHistoryOpts{
							DeviceID: deviceID,
							Date:     dateStr,
							Limit:    1,
						})
						if err != nil {
							// todo check error is not_found or not
						}

						if len(deviceWorkingHistories) == 0 {
							p.deviceWorkingHistoryRepo.Insert(ctx, &model.DeviceWorkingHistory{
								ID:                           idutil.ULIDNow(),
								ProductionOrderStageDeviceID: device.ID,
								DeviceID:                     deviceID,
								Date:                         dateStr,
								Quantity:                     int64(counterByDay),
								WorkingTime:                  0,
								CreatedAt:                    now,
							})
						} else {
							deviceWorkingHistory := deviceWorkingHistories[0].DeviceWorkingHistory
							deviceWorkingHistory.Quantity = int64(counterByDay)
							err := p.deviceWorkingHistoryRepo.Update(ctx, deviceWorkingHistory)
							if err != nil {
								// todo nothing
								p.logger.Error(" p.deviceWorkingHistoryRepo.Update error", zap.Error(err))
							}
						}
					}
					// upsert to device_working_history
					if action == "TG_in_Ngay" {
						TG_in_1Ngay := item.Value
						p.logger.Info("TG_in_Ngay", zap.String("deviceID", deviceID), zap.Int("TG_in_Ngay", int(TG_in_1Ngay)))

						deviceWorkingHistories, err := p.deviceWorkingHistoryRepo.Search(ctx, &repository.SearchDeviceWorkingHistoryOpts{
							DeviceID: deviceID,
							Date:     dateStr,
							Limit:    1,
						})
						if err != nil {
							// todo check error is not_found or not
						}

						if len(deviceWorkingHistories) == 0 {
							p.deviceWorkingHistoryRepo.Insert(ctx, &model.DeviceWorkingHistory{
								ID:                           idutil.ULIDNow(),
								ProductionOrderStageDeviceID: device.ID,
								DeviceID:                     deviceID,
								Date:                         dateStr,
								WorkingTime:                  int64(TG_in_1Ngay),
								CreatedAt:                    now,
							})
						} else {
							deviceWorkingHistory := deviceWorkingHistories[0].DeviceWorkingHistory
							deviceWorkingHistory.WorkingTime = int64(TG_in_1Ngay)
							err := p.deviceWorkingHistoryRepo.Update(ctx, deviceWorkingHistory)
							if err != nil {
								// todo nothing
								p.logger.Error(" p.deviceWorkingHistoryRepo.Update error", zap.Error(err))
							}
						}
					}
				}
			}
			// insert event log
			_ = p.productionOrderStageDeviceRepo.InsertEventLog(ctx, &model.EventLog{
				ID:        time.Now().UnixNano(),
				DeviceID:  deviceID,
				StageID:   cockroach.String(activeStageID),
				Quantity:  item.Value,
				Msg:       cockroach.String(string(message.Payload())),
				Date:      cockroach.String(dateStr),
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
