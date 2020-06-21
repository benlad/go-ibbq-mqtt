/*
   Copyright 2018 the original author or authors

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

   Build and run...
   go build
   sudo setcap cap_net_admin,cap_net_raw+eip go-ibbq-mqtt
   HA_AUTO_DISCOVERY=TRUE LOGXI=* ./go-ibbq-mqtt
*/
package main

import (
	"context"
	"strconv"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/go-ble/ble"
	"github.com/joho/godotenv"
	log "github.com/mgutz/logxi/v1"
	"github.com/sworisbreathing/go-ibbq/v2"
)

var logger = log.New("main")
var mc = NewMqttClient()
var bbq ibbq.Ibbq
var batteryLevelConfigMessage AutoDiscoverConfigMessage
var tempSensorConfigMessage AutoDiscoverConfigMessage

var unitsChangeMessageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	payloadStr := string(msg.Payload())
	if payloadStr == "C" {
		logger.Info("Configuring for °C")
		bbq.ConfigureTemperatureCelsius()
	} else if payloadStr == "F" {
		logger.Info("Configuring for °F")
		bbq.ConfigureTemperatureFahrenheit()
	} else {
		logger.Warn("Configuuration of invalid units")
	}
}

func temperatureReceived(temperatures []int) {
	logger.Info("Received temperature data", "temperatures", temperatures)

	stateMessages := NewStateMessages(temperatures)
	//	rssiMessage := NewRssiStateMessageJson(rssi)

	if getEnvBool("HA_AUTO_DISCOVERY") {
		mc.PubRaw(tempSensorConfigMessage.StateTopic, stateMessages.toJson())
		//		mc.PubRaw("inky", rssiMessage)
	} else {
		t := &temperature{temperatures}
		mc.Pub("temperatures", t.toJson())
	}
}

func batteryLevelReceived(level int) {
	logger.Info("Received battery data", "batteryPct", strconv.Itoa(level))

	batteryLevelStateMessage := NewBatteryStateMessageJson(level)

	if getEnvBool("HA_AUTO_DISCOVERY") {
		mc.PubRaw(batteryLevelConfigMessage.StateTopic, batteryLevelStateMessage)
	} else {
		b := &batteryLevel{level}
		mc.Pub("batterylevel", b.toJson())
	}
}

func statusUpdated(ibbqStatus ibbq.Status) {
	logger.Info("Status updated", "status", ibbqStatus)

	s := &status{string(ibbqStatus)}
	mc.Pub("status", s.toJson())
}

func disconnectedHandler(cancel func(), done chan struct{}) func() {
	return func() {
		logger.Info("Disconnected")
		cancel()
		close(done)
	}
}

func configureEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", "err", err)
	}
}

type debug_bbq struct {
	Addr string
}

func (bbq *debug_bbq) GetAddr() string {
	return bbq.Addr
}

func initializeiBbq(ctx context.Context, cancel context.CancelFunc, done chan struct{}) {
	logger.Debug("instantiating ibbq structs")
	//bbq := debug_bbq{Addr: "fg:11:ab:22:cd:33"}

	var err error
	var config ibbq.Configuration
	logger.Debug("instantiated ibbq structs")

	if config, err = ibbq.NewConfiguration(60*time.Second, 5*time.Minute); err != nil {
		logger.Fatal("Error creating configuration", "err", err)
	}

	logger.Info("Connecting to device")
	if bbq, err = ibbq.NewIbbq(ctx, config, disconnectedHandler(cancel, done), temperatureReceived, batteryLevelReceived, statusUpdated); err != nil {
		logger.Fatal("Error creating iBBQ", "err", err)
	} else {
		// Do home assistant discovery
		if getEnvBool("HA_AUTO_DISCOVERY") {

		}
	}

	if err = bbq.Connect(); err != nil {
		logger.Fatal("Error connecting to device", "err", err)
	} else {

		if getEnvBool("HA_AUTO_DISCOVERY") {
			logger.Info("Publish for Home Assistant MQTT auto discovery", "status")

			batteryLevelConfigMessage = NewTemperatureSensorBatteryConfigMessage(bbq.GetAddr())
			mc.PubRaw("homeassistant/sensor/"+GetMessageObjectId(bbq.GetAddr())+"/battery/config", batteryLevelConfigMessage.toJson())

			tempSensorConfigMessage = NewTemperatureSensorConfigMessage(1, bbq.GetAddr())
			mc.PubRaw("homeassistant/sensor/"+GetMessageObjectId(bbq.GetAddr())+"/temperature1/config", tempSensorConfigMessage.toJson())
			tempSensorConfigMessage.SetConfigMessageSensorNumber(2)
			mc.PubRaw("homeassistant/sensor/"+GetMessageObjectId(bbq.GetAddr())+"/temperature2/config", tempSensorConfigMessage.toJson())
			tempSensorConfigMessage.SetConfigMessageSensorNumber(3)
			mc.PubRaw("homeassistant/sensor/"+GetMessageObjectId(bbq.GetAddr())+"/temperature3/config", tempSensorConfigMessage.toJson())
			tempSensorConfigMessage.SetConfigMessageSensorNumber(4)
			mc.PubRaw("homeassistant/sensor/"+GetMessageObjectId(bbq.GetAddr())+"/temperature4/config", tempSensorConfigMessage.toJson())
			tempSensorConfigMessage.SetConfigMessageSensorNumber(5)
			mc.PubRaw("homeassistant/sensor/"+GetMessageObjectId(bbq.GetAddr())+"/temperature5/config", tempSensorConfigMessage.toJson())
			tempSensorConfigMessage.SetConfigMessageSensorNumber(6)
			mc.PubRaw("homeassistant/sensor/"+GetMessageObjectId(bbq.GetAddr())+"/temperature6/config", tempSensorConfigMessage.toJson())

			availabilityStateMessage := NewAvailabilityStateMessageJson(Online)
			mc.PubRaw(GetMessageStateTopicAvailability(bbq.GetAddr()), availabilityStateMessage)

			//			switchConfigMessage := NewUnitsSwitchConfigMessage(bbq.GetAddr())
			mc.SubRaw(GetMessageObjectId(bbq.GetAddr())+"/units", unitsChangeMessageHandler)
			// inkbird_f8300232744d/units
		}
	}
	logger.Info("Connected to device")
}

func main() {
	logger.Info(`
	_____ ____        _  ____  ____  ____        _      ____  _____  _____ 
	/  __//  _ \      / \/  _ \/  _ \/  _ \      / \__/|/  _ \/__ __\/__ __\
	| |  _| / \|_____ | || | //| | //| / \|_____ | |\/||| / \|  / \    / \  
	| |_//| \_/|\____\| || |_\\| |_\\| \_\|\____\| |  ||| \_\|  | |    | |  
	\____\\____/      \_/\____/\____/\____\      \_/  \|\____\  \_/    \_/  
																	
`)
	configureEnv()

	logger.Debug("initializing context")
	ctx1, cancel := context.WithCancel(context.Background())
	defer cancel()
	registerInterruptHandler(cancel, ctx1)
	ctx := ble.WithSigHandler(ctx1, cancel)
	logger.Debug("context initialized")

	mc.Init()

	done := make(chan struct{})
	initializeiBbq(ctx, cancel, done)

	<-ctx.Done()
	<-done
	logger.Info("Exiting")
}
