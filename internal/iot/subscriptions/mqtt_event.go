package subscriptions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
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
	productionOrderStageDeviceRepo  repository.ProductionOrderStageDeviceRepo
	deviceWorkingHistoryRepo        repository.DeviceWorkingHistoryRepo
	deviceProgressStatusHistoryRepo repository.DeviceProgressStatusHistoryRepo
	device                          model.Device
	logger                          *zap.Logger
	wsService                       ws.WebSocketService
	redisDB                         redis.Cmdable
}

func NewMQTTSubscription(
	config *configs.Config,
	db cockroach.Ext,
	productionOrderStageDeviceRepo repository.ProductionOrderStageDeviceRepo,
	deviceWorkingHistoryRepo repository.DeviceWorkingHistoryRepo,
	deviceProgressStatusHistoryRepo repository.DeviceProgressStatusHistoryRepo,
	logger *zap.Logger,
	wsService ws.WebSocketService,
	redisDB redis.Cmdable,
) *EventMQTTSubscription {
	return &EventMQTTSubscription{
		config:                          config,
		db:                              db,
		productionOrderStageDeviceRepo:  productionOrderStageDeviceRepo,
		deviceWorkingHistoryRepo:        deviceWorkingHistoryRepo,
		deviceProgressStatusHistoryRepo: deviceProgressStatusHistoryRepo,
		logger:                          logger,
		wsService:                       wsService,
		redisDB:                         redisDB,
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
	DeviceBroken    bool
	DeviceId        string
}

func parseIotData(jsonData []byte) (IotData, error) {
	var rawData map[string]interface{}
	if err := json.Unmarshal(jsonData, &rawData); err != nil {
		return IotData{}, err
	}
	data := IotData{}
	for key, value := range rawData {
		switch key {
		//case "OnOff":
		//	if v, ok := value.(float64); ok {
		//		data. = v == 1
		//	}
		case "MaMay":
			if v, ok := value.(string); ok {
				data.DeviceId = v
			}
		case "SoLuongSX":
			if v, ok := value.(float64); ok {
				data.Quantity = int64(v)
			}
		case "MayHu":
			if v, ok := value.(float64); ok {
				data.DeviceBroken = v == 1
			}
		case "MaLenhLamViec":
			switch v := value.(type) {
			case float64:
				data.WorkOrderID = strconv.FormatInt(int64(v), 10)
			case string:
				data.WorkOrderID = v
			}
		case "SanXuat":
			if v, ok := value.(float64); ok {
				data.StartProduction = v == 1
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

	//client.AddRoute("toppic1", func(client mqtt.Client, message mqtt.Message) {
	//	fmt.Println("============>>> Received message for topic 1: ", message.MessageID())
	//	message.Ack()
	//})

	client.AddRoute(topic, func(client mqtt.Client, message mqtt.Message) {
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
		if iotData.DeviceBroken == true {
			e := &model.Device{}
			err := cockroach.FindOne(ctx, e, "id = $1", iotData.DeviceId)
			if err != nil {
				fmt.Println("Lỗi không tìm thấy device: ", iotData.DeviceId, " ---> ", err)
				return
			}
			if e.Status == enum.CommonStatusDamage {
				return
			}
			updaterDevice := cockroach.NewUpdater(e.TableName(), model.DeviceFieldID, iotData.DeviceId)
			updaterDevice.Set(model.DeviceFieldStatus, enum.CommonStatusDamage)
			if err := cockroach.UpdateFields(ctx, updaterDevice); err != nil {
				fmt.Println("Lỗi update device: ", iotData.DeviceId, " ---> ", err)
			}
		}
		upsertDeviceWorkingHistory := func(ctx context.Context, orderStageDevice *model.ProductionOrderStageDevice, item IotData, dateStr string, now time.Time) error {
			redisKey := fmt.Sprintf("device_working_history:%s:%s", orderStageDevice.DeviceID, dateStr)

			deviceHistoryID, err := p.redisDB.Get(redisKey).Result()
			if errors.Is(err, redis.Nil) {
				id := idutil.ULIDNow()
				newHistory := &model.DeviceWorkingHistory{
					ID:                           id,
					ProductionOrderStageDeviceID: cockroach.String(orderStageDevice.ID),
					DeviceID:                     orderStageDevice.DeviceID,
					Date:                         dateStr,
					NumberOfPrintsPerDay:         item.Quantity,
					CreatedAt:                    now,
				}

				expiration := time.Until(time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location()))
				return cockroach.ExecInTx(ctx, func(tx context.Context) error {

					if err := p.deviceWorkingHistoryRepo.Insert(ctx, newHistory); err != nil {
						return err
					}

					return p.redisDB.Set(redisKey, id, expiration).Err()
				})

			} else if err != nil {
				return err
			}

			if item.Quantity == orderStageDevice.Quantity {
				return nil
			}

			history := model.DeviceWorkingHistory{}
			updater := cockroach.NewUpdater(history.TableName(), model.DeviceWorkingHistoryFieldID, deviceHistoryID)
			updater.Set(model.DeviceWorkingHistoryFieldQuantity, item.Quantity)
			updater.Set(model.DeviceWorkingHistoryFieldUpdatedAt, cockroach.Time(now))

			return cockroach.UpdateFields(ctx, updater)
		}

		upsertDeviceProcessHistory := func(ctx context.Context, orderStageDevice *model.ProductionOrderStageDevice, deviceStateStatus enum.ProductionOrderStageDeviceStatus, note string, now time.Time) error {
			modelData := &model.DeviceProgressStatusHistory{
				ID:                           idutil.ULIDNow(),
				ProductionOrderStageDeviceID: orderStageDevice.ID,
				DeviceID:                     orderStageDevice.DeviceID,
				ProcessStatus:                deviceStateStatus,
				CreatedBy:                    cockroach.String("HMI"),
				CreatedAt:                    now,
			}
			if deviceStateStatus == enum.ProductionOrderStageDeviceStatusFailed || deviceStateStatus == enum.ProductionOrderStageDeviceStatusPause {
				modelData.IsResolved = 0
				modelData.ErrorCode = cockroach.String(note)
				modelData.ErrorReason = cockroach.String(note)
				modelData.Description = cockroach.String(note)
			}
			return p.deviceProgressStatusHistoryRepo.Insert(ctx, modelData)
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
			hasUpdate := false
			mapping := map[int]string{
				1:  "PP1",
				2:  "CN3",
				3:  "PP4",
				4:  "CN2",
				5:  "CN1",
				6:  "VL1",
				7:  "MA1",
				8:  "NAT",
				9:  "MAU",
				10: "KH1",
			}
			note := ""
			if item.PauseReason > 10 || item.PauseReason < 1 {
				item.PauseReason = 10
			}
			note = mapping[item.PauseReason]
			tableStageDevice := model.ProductionOrderStageDevice{}
			orderStageDevice, err := p.productionOrderStageDeviceRepo.FindByID(ctx, iotData.WorkOrderID)
			fmt.Println("============>>> orderStageDevice: ", orderStageDevice)
			if err != nil {
				p.logger.Error("Error finding active order stage device", zap.Error(err))
				return err
			}
			if orderStageDevice == nil {
				fmt.Println("============>>> D: ", orderStageDevice)
				return nil
			}
			if orderStageDevice.DeviceID != item.DeviceId {
				fmt.Println("============>>> DeviceID mismatch: ", orderStageDevice.DeviceID, item.DeviceId)
				return nil
			}
			if orderStageDevice.ProcessStatus == enum.ProductionOrderStageDeviceStatusComplete {
				fmt.Println("============>>> E: ", orderStageDevice)
				return nil
			}
			deviceStateStatus := orderStageDevice.ProcessStatus
			fmt.Println("============>>> deviceStateStatus: ", deviceStateStatus)
			switch orderStageDevice.ProcessStatus {
			case enum.ProductionOrderStageDeviceStatusFailed:
				if !item.Pause || item.StartProduction {
					deviceStateStatus = enum.ProductionOrderStageDeviceStatusStart
				} else {
					return nil
				}
			case enum.ProductionOrderStageDeviceStatusNone:
				if item.StartProduction {
					deviceStateStatus = enum.ProductionOrderStageDeviceStatusStart
					fmt.Println("============>>> deviceStateStatus After 279: ", deviceStateStatus)
				}
				if item.Pause {
					return nil
				}
				item.Quantity = 0
			case enum.ProductionOrderStageDeviceStatusStart:
				if item.Pause {
					deviceStateStatus = enum.ProductionOrderStageDeviceStatusFailed
					fmt.Println("============>>> deviceStateStatus After 306: ", deviceStateStatus)
				} else if !item.StartProduction {
					deviceStateStatus = enum.ProductionOrderStageDeviceStatusComplete
					fmt.Println("============>>> deviceStateStatus After 317: ", deviceStateStatus)
				}
			case enum.ProductionOrderStageDeviceStatusPause:
				break
			default:
				p.logger.Warn("Unexpected device state status", zap.Any("deviceStateStatus", deviceStateStatus))
			}
			if orderStageDevice.ProcessStatus == enum.ProductionOrderStageDeviceStatusFailed {
				if item.Pause == false || item.StartProduction {
					deviceStateStatus = enum.ProductionOrderStageDeviceStatusStart
				} else {
					return nil
				}
			} else {
				if orderStageDevice.ProcessStatus == enum.ProductionOrderStageDeviceStatusNone {
					if item.StartProduction {
						deviceStateStatus = enum.ProductionOrderStageDeviceStatusStart
					}
					if item.Pause {
						return nil
					}
					item.Quantity = 0
				} else if orderStageDevice.ProcessStatus == enum.ProductionOrderStageDeviceStatusStart {
					if item.Pause {
						deviceStateStatus = enum.ProductionOrderStageDeviceStatusFailed
					} else if item.StartProduction == false {
						deviceStateStatus = enum.ProductionOrderStageDeviceStatusComplete
					}
					//} else if orderStageDevice.ProcessStatus == enum.ProductionOrderStageDeviceStatusFailed {
					//	if item.StartProduction {
					//		deviceStateStatus = enum.ProductionOrderStageDeviceStatusStart
					//		fmt.Println("============>>> deviceStateStatus After 322: ", deviceStateStatus)
					//	}
				} else {
					p.logger.Warn("Unexpected device state status", zap.Any("deviceStateStatus", deviceStateStatus))
					return nil
				}
			}
			updaterDeviceStage := cockroach.NewUpdater(tableStageDevice.TableName(), model.ProductionOrderStageDeviceFieldID, orderStageDevice.ID)
			fmt.Println("============>>> deviceStateStatus After 2: ", deviceStateStatus)
			if item.DefectQuantity > 0 {
				if orderStageDevice.Settings == nil {
					orderStageDevice.Settings = make(map[string]interface{})
				}

				val, ok := orderStageDevice.Settings["san_pham_loi"].(int64)
				if !ok || item.DefectQuantity != val {
					orderStageDevice.Settings["san_pham_loi"] = item.DefectQuantity
					updaterDeviceStage.Set(model.ProductionOrderStageDeviceFieldSettings, orderStageDevice.Settings)
					hasUpdate = true
				}
			}
			if item.Quantity > 0 && item.Quantity != orderStageDevice.Quantity {
				updaterDeviceStage.Set(model.ProductionOrderStageDeviceFieldQuantity, item.Quantity)
				hasUpdate = true
			}
			//err = p.productionOrderStageDeviceService.Edit(ctx, &production_order_stage_device.EditProductionOrderStageDeviceOpts{
			//	ID:                  orderStageDevice.ID,
			//	DeviceID:            orderStageDevice.DeviceID,
			//	UserID:              "HMI",
			//	ProcessStatus:       deviceStateStatus,
			//	Status:              orderStageDevice.Status,
			//	SanPhamLoi:          item.DefectQuantity,
			//	Quantity:            item.Quantity,
			//	Settings:            settings,
			//	Note:                note,
			//	EstimatedStartAt:    orderStageDevice.EstimatedStartAt,
			//	EstimatedCompleteAt: orderStageDevice.EstimatedCompleteAt,
			//})
			if deviceStateStatus != orderStageDevice.ProcessStatus {
				updaterDeviceStage.Set(model.ProductionOrderStageDeviceFieldProcessStatus, deviceStateStatus)
				hasUpdate = true
			}
			if deviceStateStatus == enum.ProductionOrderStageDeviceStatusStart && orderStageDevice.ProcessStatus == enum.ProductionOrderStageDeviceStatusNone {
				updaterDeviceStage.Set(model.ProductionOrderStageDeviceFieldStartAt, now)
				hasUpdate = true
			}
			if deviceStateStatus == enum.ProductionOrderStageDeviceStatusComplete && orderStageDevice.ProcessStatus == enum.ProductionOrderStageDeviceStatusStart {
				updaterDeviceStage.Set(model.ProductionOrderStageDeviceFieldCompleteAt, now)
				hasUpdate = true
			}

			if hasUpdate {
				if err := cockroach.UpdateFields(ctx, updaterDeviceStage); err != nil {
					p.logger.Error("Error tableStageDevice", zap.Error(err))
					return err
				}
			}
			if deviceStateStatus != orderStageDevice.ProcessStatus {
				if err := upsertDeviceProcessHistory(ctx, orderStageDevice, deviceStateStatus, note, now); err != nil {
					p.logger.Error("Error upsertDeviceProcessHistory", zap.Error(err))
					return err
				}
			}

			if err := upsertDeviceWorkingHistory(ctx, orderStageDevice, item, dateStr, now); err != nil {
				p.logger.Error("Error upsertDeviceWorkingHistory", zap.Error(err))
				return err
			}
			if err := insertEventLog(ctx, orderStageDevice, item, dateStr, now, jsonStr); err != nil {
				p.logger.Error("Error insertEventLog", zap.Error(err))
				return err
			}
			return nil
		}
		err = processIotData(ctx, iotData, dateStr, now, string(jsonStr))
		if err != nil {
			return
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
