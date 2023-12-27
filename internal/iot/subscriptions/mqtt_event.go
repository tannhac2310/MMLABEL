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

type IotParseData struct {
	DeviceID   string
	SL_in_Ngay int64
	TG_in_Ngay int64
	SL_in_LSX  int64
	TG_in_LSX  int64
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

		now := time.Now()
		dateStr := now.Format("2006-01-02")
		// iotData.D is {"d":[{"tag":"B_PR03:SL_in_LSX","value":259.00},{"tag":"B_PR03:ON_OFF","value":1},{"tag":"B_PR03:SL_in_Ngay","value":143.00},{"tag":"B_PR03:TG_in_Ngay","value":32.00},{"tag":"B_PR03:SL_in_LSX_1","value":0.00},{"tag":"B_PR03:TG_in_LSX_1","value":0.00},{"tag":"B_PR03:SL_in_LSX_2","value":0.00},{"tag":"B_PR03:TG_in_LSX_2","value":0.00},{"tag":"B_PR03:SL_in_LSX_3","value":0.00},{"tag":"B_PR03:TG_in_LSX_3","value":0.00},{"tag":"B_PR03:SL_in_LSX_4","value":0.00},{"tag":"B_PR03:TG_in_LSX_4","value":0.00},{"tag":"B_PR03:TG_in_LSX","value":32.00},{"tag":"B_PR04:SL_in_LSX","value":0.00},{"tag":"B_PR04:ON_OFF","value":1},{"tag":"B_PR04:SL_in_Ngay","value":0.00},{"tag":"B_PR04:TG_in_Ngay","value":0.00},{"tag":"B_PR04:TG_in_LSX","value":0.00},{"tag":"B_PR04:SL_in_LSX_1","value":0.00},{"tag":"B_PR04:TG_in_LSX_1","value":0.00},{"tag":"B_PR04:SL_in_LSX_2","value":0.00},{"tag":"B_PR04:TG_in_LSX_2","value":0.00},{"tag":"B_PR04:SL_in_LSX_3","value":0.00},{"tag":"B_PR04:TG_in_LSX_3","value":0.00},{"tag":"B_PR04:SL_in_LSX_4","value":0.00},{"tag":"B_PR04:TG_in_LSX_4","value":0.00}],"ts":"2023-12-27T01:19:15+0000"}
		// group by by device id
		//
		mappingData := make(map[string]*IotParseData, 0)
		for _, item := range iotData.D {

			d := strings.Split(item.Tag, ":")
			if len(d) != 2 {
				continue
			}
			deviceID := d[0]
			action := d[1]
			if _, ok := mappingData[deviceID]; !ok {
				mappingData[deviceID] = &IotParseData{
					DeviceID:   deviceID,
					SL_in_Ngay: 0,
					TG_in_Ngay: 0,
					SL_in_LSX:  0,
					TG_in_LSX:  0,
				}
			}
			if action == "SL_in_Ngay" {
				mappingData[deviceID].SL_in_Ngay = int64(item.Value)
			}
			if action == "TG_in_Ngay" {
				mappingData[deviceID].TG_in_Ngay = int64(item.Value)
			}
			if action == "SL_in_LSX" {
				mappingData[deviceID].SL_in_LSX = int64(item.Value)
			}
			if action == "TG_in_LSX" {
				mappingData[deviceID].TG_in_LSX = int64(item.Value)
			}
		}
		fmt.Println(iotData.D, iotData.Ts)
		for deviceID, item := range mappingData {
			// find device in production order stage device

			p.logger.Info("IOT+PARSE+DATA",
				zap.Any("deviceID", deviceID),
				zap.Any("mappingData", item))

			orderStageDevices, err := p.productionOrderStageDeviceRepo.Search(ctx, &repository.SearchProductionOrderStageDevicesOpts{
				DeviceID:                   deviceID,
				ProductionOrderStageStatus: enum.ProductionOrderStageStatusProductionStart,
				Limit:                      1,
				Offset:                     0,
			})

			if err != nil {
				// todo nothing
				p.logger.Error(" p.productionOrderStageDeviceRepo.Search error", zap.Error(err))
			}
			activeStageID := ""
			orderStageDeviceID := ""
			if len(orderStageDevices) >= 1 && orderStageDevices[0].ProcessStatus == enum.ProductionOrderStageDeviceStatusStart {
				device := orderStageDevices[0]
				orderStageDeviceID = device.ID
				activeStageID = device.ProductionOrderStageID
				if item.SL_in_LSX > 0 && item.SL_in_LSX > device.Quantity {
					// update first device
					_ = p.productionOrderStageDeviceRepo.Update(ctx, &model.ProductionOrderStageDevice{
						ID:                     device.ID,
						ProductionOrderStageID: device.ProductionOrderStageID,
						DeviceID:               device.DeviceID,
						Quantity:               item.SL_in_LSX,
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
			}
			// upsert to device_working_history
			deviceWorkingHistories, err := p.deviceWorkingHistoryRepo.Search(ctx, &repository.SearchDeviceWorkingHistoryOpts{
				DeviceID: deviceID,
				Date:     dateStr,
				Limit:    1,
			})
			if err != nil {
				// todo nothing
				p.logger.Error(" p.deviceWorkingHistoryRepo.Search error", zap.Error(err))
			}
			if len(deviceWorkingHistories) == 0 {
				p.deviceWorkingHistoryRepo.Insert(ctx, &model.DeviceWorkingHistory{
					ID:                           idutil.ULIDNow(),
					ProductionOrderStageDeviceID: cockroach.String(orderStageDeviceID),
					DeviceID:                     deviceID,
					Date:                         dateStr,
					Quantity:                     item.SL_in_Ngay, // todo remove this field in db
					WorkingTime:                  item.TG_in_Ngay, // todo remove this field in db
					NumberOfPrintsPerDay:         item.SL_in_Ngay,
					PrintingTimePerDay:           item.TG_in_Ngay,
					PoQuantity:                   item.SL_in_LSX,
					PoWorkingTime:                item.TG_in_LSX,
					CreatedAt:                    now,
				})
			} else {
				// update
				deviceWorkingHistory := deviceWorkingHistories[0].DeviceWorkingHistory

				deviceWorkingHistory.NumberOfPrintsPerDay = item.SL_in_Ngay
				deviceWorkingHistory.PrintingTimePerDay = item.TG_in_Ngay
				deviceWorkingHistory.Quantity = item.SL_in_Ngay    // todo remove this field in db
				deviceWorkingHistory.WorkingTime = item.TG_in_Ngay // todo remove this field in db
				deviceWorkingHistory.PoWorkingTime = item.TG_in_LSX
				deviceWorkingHistory.PoQuantity = item.SL_in_LSX

				deviceWorkingHistory.UpdatedAt = cockroach.Time(now)

				err := p.deviceWorkingHistoryRepo.Update(ctx, deviceWorkingHistory)
				if err != nil {
					// todo nothing
					p.logger.Error(" p.deviceWorkingHistoryRepo.Update error", zap.Error(err))
				}
			}

			fmt.Println("activeStageID", activeStageID)
			// insert event log
			_ = p.productionOrderStageDeviceRepo.InsertEventLog(ctx, &model.EventLog{
				ID:          time.Now().UnixNano(),
				DeviceID:    deviceID,
				StageID:     cockroach.String(activeStageID),
				StageStatus: nil, // todo check stage status
				Quantity:    float64(item.SL_in_Ngay),
				Msg:         cockroach.String(string(message.Payload())),
				Date:        cockroach.String(dateStr),
				CreatedAt:   now,
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
