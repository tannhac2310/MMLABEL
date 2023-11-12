package subscriptions

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"

	"mmlabel.gitlab.com/mm-printing-backend/pkg/database/cockroach"
	"mmlabel.gitlab.com/mm-printing-backend/pkg/service/ws"
)

type EventMQTTSubscription struct {
	db cockroach.Ext
	//busFactory nats.BusFactory
	logger    *zap.Logger
	wsService ws.WebSocketService
}

func NewMQTTSubscription(
	db cockroach.Ext,
	//busFactory nats.BusFactory,
	logger *zap.Logger,
	wsService ws.WebSocketService,
) *EventMQTTSubscription {
	return &EventMQTTSubscription{
		db: db,
		//busFactory: busFactory,
		logger:    logger,
		wsService: wsService,
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
	opts.SetClientID("go_mqtt_client")
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
		fmt.Println("============>>> Received message for topic", string(message.Payload()))
		message.Ack()
	})
	fmt.Printf("Subscribed to topic '%s'\n", topic)
	return nil
}
