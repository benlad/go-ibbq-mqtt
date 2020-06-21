package main

import (
	"fmt"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// MqttClient is the public interface
type MqttClient interface {
	Init()
	Pub(topic string, payload interface{})
	PubRaw(topic string, payload interface{})
	SubRaw(topic string, messageHandler mqtt.MessageHandler)
}

// mqttClient implements the MqttClient interface, encapsulating the paho.mqtt.golang module.
type mqttClient struct {
	client mqtt.Client
}

func NewMqttClient() MqttClient {
	c := &mqttClient{}

	return c
}

func (m *mqttClient) Init() {
	logger.Info("Connecting to mqtt broker", "broker", os.Getenv("MQTT_SERVER"))

	opts := mqtt.NewClientOptions().AddBroker(os.Getenv("MQTT_SERVER")).SetClientID("go-ibbq-mqtt")
	opts.SetKeepAlive(2 * time.Second)
	opts.SetPingTimeout(1 * time.Second)

	m.client = mqtt.NewClient(opts)

	if token := m.client.Connect(); token.Wait() && token.Error() != nil {
		logger.Fatal("Error connecting to mqtt", "err", token.Error())
	}

	logger.Info("Connected to mqtt broker", "broker", os.Getenv("MQTT_SERVER"))
}

func (m *mqttClient) Pub(topic string, payload interface{}) {
	t := getTopic(topic)

	logger.Info("Publishing to mqtt", "topic", t)
	statustoken := m.client.Publish(t, 0, false, payload)

	statustoken.Wait()
	if statustoken.Error() != nil {
		logger.Error("Error publishing to mqtt", "err", statustoken.Error())
	}
}

func (m *mqttClient) PubRaw(topic string, payload interface{}) {
	t := topic

	logger.Info("Publishing to mqtt", "topic", t)
	statustoken := m.client.Publish(t, 0, true, payload)

	statustoken.Wait()
	if statustoken.Error() != nil {
		logger.Error("Error publishing to mqtt", "err", statustoken.Error())
	}
}

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("MSG: %s\n", msg.Payload())
	logger.Info("*************************************** ********** *****  *** * Received MQTT messate payload", msg.Payload())
}

func (m *mqttClient) SubRaw(topic string, messageHandler mqtt.MessageHandler) {
	t := topic

	logger.Info("Subscribing to mqtt", "topic", t)
	statustoken := m.client.Subscribe(t, 0, messageHandler)

	if statustoken.Error() != nil {
		logger.Error("Error subscribing to mqtt", "err", statustoken.Error())
	}
}

func getTopic(topic string) string {
	return fmt.Sprintf("%s/%s", os.Getenv("MQTT_TOPIC"), topic)
}
