package subscriptions

import (
	"context"
	"encoding/json"
	"fmt"
	"mmlabel.gitlab.com/mm-printing-backend/internal/aurora/service/production_order_stage_device"
	"strconv"
	"time"

	"mmlabel.gitlab.com/mm-printing-backend/internal/iot/configs"
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
	config *configs.Config
	db     cockroach.Ext
	//busFactory nats.BusFactory
	productionOrderStageDeviceRepo    repository.ProductionOrderStageDeviceRepo
	deviceWorkingHistoryRepo          repository.DeviceWorkingHistoryRepo
	productionOrderStageDeviceService production_order_stage_device.Service
	logger                            *zap.Logger
	wsService                         ws.WebSocketService
}

func NewMQTTSubscription(
	config *configs.Config,
	db cockroach.Ext,
	productionOrderStageDeviceRepo repository.ProductionOrderStageDeviceRepo,
	deviceWorkingHistoryRepo repository.DeviceWorkingHistoryRepo,
	productionOrderStageDeviceService production_order_stage_device.Service,
	logger *zap.Logger,
	wsService ws.WebSocketService,
) *EventMQTTSubscription {
	return &EventMQTTSubscription{
		config:                            config,
		db:                                db,
		productionOrderStageDeviceRepo:    productionOrderStageDeviceRepo,
		deviceWorkingHistoryRepo:          deviceWorkingHistoryRepo,
		productionOrderStageDeviceService: productionOrderStageDeviceService,
		logger:                            logger,
		wsService:                         wsService,
	}
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("==============================================>>>  Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("==============================================>>>>>Connected======>>>>>")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
	panic("Connect lost at " + time.Now().String())
}

type IotParseData struct {
	DeviceID   string
	SL_in_Ngay int64
	TG_in_Ngay int64
	SL_in_LSX  int64
	TG_in_LSX  int64
}

type IotData struct {
	Quantity        int64
	WorkOrderID     string
	TestProduction  bool
	StartProduction bool
	StopPO          bool
	Pause           bool
	Setup           bool
	PauseReason     int
	DefectQuantity  int64
	PrintPerDay     int64
}

func parseIotData(jsonData []byte) (IotData, error) {
	var rawData map[string]interface{}
	if err := json.Unmarshal(jsonData, &rawData); err != nil {
		return IotData{}, err
	}
	data := IotData{}
	for key, value := range rawData {
		switch key {
		case "OnOff":
			if v, ok := value.(float64); ok {
				data.StartProduction = v == 1
			}
		case "SetUp":
			if v, ok := value.(float64); ok {
				data.Setup = v == 1
			}
		case "SoLuongSX":
			if v, ok := value.(float64); ok {
				data.Quantity = int64(v)
			}
		case "MaLenhLamViec":
			switch v := value.(type) {
			case float64:
				data.WorkOrderID = strconv.FormatInt(int64(v), 10)
			case string:
				data.WorkOrderID = v
			}
		case "SXThu":
			if v, ok := value.(float64); ok {
				data.TestProduction = v == 1
			}
		case "SX":
			if v, ok := value.(float64); ok {
				data.StartProduction = v == 1
			}
		case "NgungPO":
			if v, ok := value.(float64); ok {
				data.StopPO = v == 1
			}
		case "TamDung":
			if v, ok := value.(float64); ok {
				data.Pause = v == 1
			}
		case "TamDungLyDo":
			if v, ok := value.(float64); ok {
				data.PauseReason = int(v)
			}
		case "SoLuongLoi":
			if v, ok := value.(float64); ok {
				data.DefectQuantity = int64(v)
			}
		case "SoLuongInNgay":
			if v, ok := value.(float64); ok {
				data.PrintPerDay = int64(v)
			}
		}
	}

	return data, nil
}

func (p *EventMQTTSubscription) Subscribe() error {
	// get from config
	b, _ := json.Marshal(p.config.MQTT)
	fmt.Println("config.MQTT", string(b))
	var broker = p.config.MQTT.Host
	var port = p.config.MQTT.Port //31883
	fmt.Println(fmt.Sprintf("tcp://%s:%s", broker, port))
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%s", broker, port))
	//opts.SetClientID("go_mqtt_client")
	opts.SetUsername(p.config.MQTT.Username)
	opts.SetPassword(p.config.MQTT.Password)
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

		now := time.Now()
		dateStr := now.Format("2006-01-02")

		//Minh
		jsonStr := message.Payload()
		iotData, err := parseIotData([]byte(jsonStr))
		if err != nil {
			fmt.Println("Lỗi:", err)
			return
		}
		upsertDeviceWorkingHistory := func(ctx context.Context, orderStageDevice *model.ProductionOrderStageDevice, item IotData, dateStr string, now time.Time) error {
			deviceWorkingHistories, err := p.deviceWorkingHistoryRepo.Search(ctx, &repository.SearchDeviceWorkingHistoryOpts{
				DeviceID: orderStageDevice.DeviceID,
				Date:     dateStr,
				Limit:    1,
			})
			if err != nil {
				return err
			}

			if len(deviceWorkingHistories) == 0 {
				// Chèn mới
				newHistory := &model.DeviceWorkingHistory{
					ID:                           idutil.ULIDNow(),
					ProductionOrderStageDeviceID: cockroach.String(orderStageDevice.ID),
					DeviceID:                     orderStageDevice.DeviceID,
					Date:                         dateStr,
					NumberOfPrintsPerDay:         item.Quantity,
					CreatedAt:                    now,
				}
				return p.deviceWorkingHistoryRepo.Insert(ctx, newHistory)
			}

			// Cập nhật
			existingHistory := deviceWorkingHistories[0].DeviceWorkingHistory
			existingHistory.Quantity = item.Quantity
			existingHistory.UpdatedAt = cockroach.Time(now)
			return p.deviceWorkingHistoryRepo.Update(ctx, existingHistory)
		}

		insertEventLog := func(ctx context.Context, orderStageDevice *model.ProductionOrderStageDevice, item IotData, dateStr string, now time.Time, payload string) error {
			eventLog := &model.EventLog{
				ID:          time.Now().UnixNano(),
				DeviceID:    orderStageDevice.DeviceID,
				StageID:     cockroach.String(orderStageDevice.ProductionOrderStageID),
				StageStatus: nil, // todo: xác định trạng thái stage
				Quantity:    float64(item.Quantity),
				Msg:         cockroach.String(payload),
				Date:        cockroach.String(dateStr),
				CreatedAt:   now,
			}
			return p.productionOrderStageDeviceRepo.InsertEventLog(ctx, eventLog)
		}
		processIotData := func(ctx context.Context, item IotData, dateStr string, now time.Time, jsonStr string) error {
			//table := model.ProductionOrderStageDevice{}
			//tableProductProgress := model.DeviceProgressStatusHistory{}
			tableDevice := model.Device{}
			orderStageDevice, err := p.productionOrderStageDeviceRepo.FindByID(ctx, iotData.WorkOrderID)
			if err != nil {
				p.logger.Error("Error finding active order stage device", zap.Error(err))
				return err
			}
			if orderStageDevice == nil {
				return nil
			}
			settings := &production_order_stage_device.Settings{}
			deviceStateStatus := orderStageDevice.ProcessStatus
			switch orderStageDevice.ProcessStatus {
			case enum.ProductionOrderStageDeviceStatusNone:
			case enum.ProductionOrderStageDeviceStatusCompleteTestProduce:
				if item.StartProduction {
					deviceStateStatus = enum.ProductionOrderStageDeviceStatusStart
				} else if item.TestProduction {
					deviceStateStatus = enum.ProductionOrderStageDeviceStatusTestProduce
				} else if item.Setup {
					deviceStateStatus = enum.ProductionOrderStageDeviceStatusSetup
				}
			case enum.ProductionOrderStageDeviceStatusTestProduce:
				if item.TestProduction == false {
					deviceStateStatus = enum.ProductionOrderStageDeviceStatusCompleteTestProduce
				}
			case enum.ProductionOrderStageDeviceStatusSetup:
				if item.StartProduction {
					deviceStateStatus = enum.ProductionOrderStageDeviceStatusStart
				}
				if item.Setup == false {
					deviceStateStatus = enum.ProductionOrderStageDeviceStatusNone
				}
			case enum.ProductionOrderStageDeviceStatusStart:
				if item.Pause {
					deviceStateStatus = enum.ProductionOrderStageDeviceStatusFailed
					orderStageDevice.Note = cockroach.String(strconv.Itoa(item.PauseReason))
					settings = &production_order_stage_device.Settings{
						DefectiveError: strconv.Itoa(item.PauseReason),
						Description:    strconv.Itoa(item.PauseReason),
					}
				} else if item.StopPO {
					deviceStateStatus = enum.ProductionOrderStageDeviceStatusPause
					orderStageDevice.Note = cockroach.String(strconv.Itoa(item.PauseReason))
					settings = &production_order_stage_device.Settings{
						DefectiveError: "StopPO",
						Description:    "StopPO",
					}
				} else if item.StartProduction == false {
					deviceStateStatus = enum.ProductionOrderStageDeviceStatusComplete
				}
			case enum.ProductionOrderStageDeviceStatusFailed:
				if item.StartProduction {
					deviceStateStatus = enum.ProductionOrderStageDeviceStatusStart
				}

			default:
				p.logger.Warn("Unexpected device state status", zap.Any("deviceStateStatus", deviceStateStatus))
			}

			err = p.productionOrderStageDeviceService.Edit(ctx, &production_order_stage_device.EditProductionOrderStageDeviceOpts{
				ID:            orderStageDevice.ID,
				DeviceID:      orderStageDevice.DeviceID,
				UserID:        "HMI",
				ProcessStatus: deviceStateStatus,
				Status:        orderStageDevice.Status,
				SanPhamLoi:    item.DefectQuantity,
				Quantity:      item.Quantity,
				Settings:      settings,
			})
			if err != nil {
				p.logger.Error("Error updating production order stage device", zap.Error(err))
				return err
			}

			if err := upsertDeviceWorkingHistory(ctx, orderStageDevice, item, dateStr, now); err != nil {
				p.logger.Error("Error upserting device working history", zap.Error(err))
				return err
			}
			if err := insertEventLog(ctx, orderStageDevice, item, dateStr, now, jsonStr); err != nil {
				p.logger.Error("Error upserting device working history", zap.Error(err))
				return err
			}

			if item.Pause {
				if item.PauseReason == 1 {
					updaterDevice := cockroach.NewUpdater(tableDevice.TableName(), model.DeviceFieldID, orderStageDevice.DeviceID)
					updaterDevice.Set(model.DeviceFieldStatus, enum.CommonStatusDamage)
					if err := cockroach.UpdateFields(ctx, updaterDevice); err != nil {
						p.logger.Error("Error upserting device working history", zap.Error(err))
						return err
					}
				}
			}

			return nil
		}

		err = processIotData(ctx, iotData, dateStr, now, string(jsonStr))
		if err != nil {
			return
		}
		//End Minh

		// iotData.D is {"d":[{"tag":"B_PR03:SL_in_LSX","value":259.00},{"tag":"B_PR03:ON_OFF","value":1},{"tag":"B_PR03:SL_in_Ngay","value":143.00},{"tag":"B_PR03:TG_in_Ngay","value":32.00},{"tag":"B_PR03:SL_in_LSX_1","value":0.00},{"tag":"B_PR03:TG_in_LSX_1","value":0.00},{"tag":"B_PR03:SL_in_LSX_2","value":0.00},{"tag":"B_PR03:TG_in_LSX_2","value":0.00},{"tag":"B_PR03:SL_in_LSX_3","value":0.00},{"tag":"B_PR03:TG_in_LSX_3","value":0.00},{"tag":"B_PR03:SL_in_LSX_4","value":0.00},{"tag":"B_PR03:TG_in_LSX_4","value":0.00},{"tag":"B_PR03:TG_in_LSX","value":32.00},{"tag":"B_PR04:SL_in_LSX","value":0.00},{"tag":"B_PR04:ON_OFF","value":1},{"tag":"B_PR04:SL_in_Ngay","value":0.00},{"tag":"B_PR04:TG_in_Ngay","value":0.00},{"tag":"B_PR04:TG_in_LSX","value":0.00},{"tag":"B_PR04:SL_in_LSX_1","value":0.00},{"tag":"B_PR04:TG_in_LSX_1","value":0.00},{"tag":"B_PR04:SL_in_LSX_2","value":0.00},{"tag":"B_PR04:TG_in_LSX_2","value":0.00},{"tag":"B_PR04:SL_in_LSX_3","value":0.00},{"tag":"B_PR04:TG_in_LSX_3","value":0.00},{"tag":"B_PR04:SL_in_LSX_4","value":0.00},{"tag":"B_PR04:TG_in_LSX_4","value":0.00}],"ts":"2023-12-27T01:19:15+0000"}
		// group by by device id
		//
		//mappingData := make(map[string]*IotParseData, 0)
		//for _, item := range iotData.D {
		//
		//	d := strings.Split(item.Tag, ":")
		//	if len(d) != 2 {
		//		continue
		//	}
		//	deviceID := d[0]
		//	action := d[1]
		//	if _, ok := mappingData[deviceID]; !ok {
		//		mappingData[deviceID] = &IotParseData{
		//			DeviceID:   deviceID,
		//			SL_in_Ngay: 0,
		//			TG_in_Ngay: 0,
		//			SL_in_LSX:  0,
		//			TG_in_LSX:  0,
		//		}
		//	}
		//	if action == "SL_in_Ngay" {
		//		mappingData[deviceID].SL_in_Ngay = int64(item.Value)
		//	}
		//	if action == "TG_in_Ngay" {
		//		mappingData[deviceID].TG_in_Ngay = int64(item.Value)
		//	}
		//	if action == "SL_in_LSX" {
		//		mappingData[deviceID].SL_in_LSX = int64(item.Value)
		//	}
		//	if action == "TG_in_LSX" {
		//		mappingData[deviceID].TG_in_LSX = int64(item.Value)
		//	}
		//}
		//fmt.Println(iotData.D, iotData.Ts)
		//for deviceID, item := range mappingData {
		//	// find device in production order stage device
		//	p.logger.Info("IOT+PARSE+DATA",
		//		zap.Any("deviceID", deviceID),
		//		zap.Any("mappingData", item))
		//
		//	orderStageDevices, err := p.productionOrderStageDeviceRepo.Search(ctx, &repository.SearchProductionOrderStageDevicesOpts{
		//		DeviceIDs:                    []string{deviceID},
		//		ProductionOrderStageStatuses: []enum.ProductionOrderStageStatus{enum.ProductionOrderStageStatusProductionStart},
		//		Limit:                        1,
		//		Offset:                       0,
		//	})
		//
		//	if err != nil {
		//		// todo nothing
		//		p.logger.Error(" p.productionOrderStageDeviceRepo.Search error", zap.Error(err))
		//	}
		//	activeStageID := ""
		//	orderStageDeviceID := ""
		//	if len(orderStageDevices) >= 1 && orderStageDevices[0].ProcessStatus == enum.ProductionOrderStageDeviceStatusStart {
		//		device := orderStageDevices[0]
		//		orderStageDeviceID = device.ID
		//		activeStageID = device.ProductionOrderStageID
		//		if item.SL_in_LSX > 0 && item.SL_in_LSX > device.Quantity {
		//			// update first device
		//			_ = p.productionOrderStageDeviceRepo.Update(ctx, &model.ProductionOrderStageDevice{
		//				ID:                     device.ID,
		//				ProductionOrderStageID: device.ProductionOrderStageID,
		//				DeviceID:               device.DeviceID,
		//				Quantity:               item.SL_in_LSX,
		//				ProcessStatus:          device.ProcessStatus,
		//				Status:                 device.Status,
		//				Settings:               device.Settings,
		//				Note:                   device.Note,
		//				CreatedAt:              device.CreatedAt,
		//				UpdatedAt:              device.UpdatedAt,
		//				Responsible:            device.Responsible,
		//				EstimatedCompleteAt:    device.EstimatedCompleteAt,
		//				AssignedQuantity:       device.AssignedQuantity,
		//			})
		//		}
		//	}
		//	// upsert to device_working_history
		//	deviceWorkingHistories, err := p.deviceWorkingHistoryRepo.Search(ctx, &repository.SearchDeviceWorkingHistoryOpts{
		//		DeviceID: deviceID,
		//		Date:     dateStr,
		//		Limit:    1,
		//	})
		//	if err != nil {
		//		// todo nothing
		//		p.logger.Error(" p.deviceWorkingHistoryRepo.Search error", zap.Error(err))
		//	}
		//	if len(deviceWorkingHistories) == 0 {
		//		p.deviceWorkingHistoryRepo.Insert(ctx, &model.DeviceWorkingHistory{
		//			ID:                           idutil.ULIDNow(),
		//			ProductionOrderStageDeviceID: cockroach.String(orderStageDeviceID),
		//			DeviceID:                     deviceID,
		//			Date:                         dateStr,
		//			Quantity:                     item.SL_in_Ngay, // todo remove this field in db
		//			WorkingTime:                  item.TG_in_Ngay, // todo remove this field in db
		//			NumberOfPrintsPerDay:         item.SL_in_Ngay,
		//			PrintingTimePerDay:           item.TG_in_Ngay,
		//			PoQuantity:                   item.SL_in_LSX,
		//			PoWorkingTime:                item.TG_in_LSX,
		//			CreatedAt:                    now,
		//		})
		//	} else {
		//		// update
		//		deviceWorkingHistory := deviceWorkingHistories[0].DeviceWorkingHistory
		//
		//		deviceWorkingHistory.NumberOfPrintsPerDay = item.SL_in_Ngay
		//		deviceWorkingHistory.PrintingTimePerDay = item.TG_in_Ngay
		//		deviceWorkingHistory.Quantity = item.SL_in_Ngay    // todo remove this field in db
		//		deviceWorkingHistory.WorkingTime = item.TG_in_Ngay // todo remove this field in db
		//		deviceWorkingHistory.PoWorkingTime = item.TG_in_LSX
		//		deviceWorkingHistory.PoQuantity = item.SL_in_LSX
		//
		//		deviceWorkingHistory.UpdatedAt = cockroach.Time(now)
		//
		//		err := p.deviceWorkingHistoryRepo.Update(ctx, deviceWorkingHistory)
		//		if err != nil {
		//			// todo nothing
		//			p.logger.Error(" p.deviceWorkingHistoryRepo.Update error", zap.Error(err))
		//		}
		//	}
		//
		//	fmt.Println("activeStageID", activeStageID)
		//	// insert event log
		//	_ = p.productionOrderStageDeviceRepo.InsertEventLog(ctx, &model.EventLog{
		//		ID:          time.Now().UnixNano(),
		//		DeviceID:    deviceID,
		//		StageID:     cockroach.String(activeStageID),
		//		StageStatus: nil, // todo check stage status
		//		Quantity:    float64(item.SL_in_Ngay),
		//		Msg:         cockroach.String(string(message.Payload())),
		//		Date:        cockroach.String(dateStr),
		//		CreatedAt:   now,
		//	})
		//}

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
