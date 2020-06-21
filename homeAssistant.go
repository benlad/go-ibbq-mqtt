package main

import (
	"encoding/json"
	"regexp"
	"strconv"

	log "github.com/mgutz/logxi/v1"
)

// The capitalized filed names are exported and thus end up in JSON
type AutoDiscoverConfigMessage struct {
	addr                string
	id                  string
	sensorNum           int
	Name                string `json:"name"`
	DeviceClass         string `json:"device_class"`
	StateTopic          string `json:"state_topic"`
	CommandTopic        string `json:"command_topic"`
	availabilityTopic   string `json:"availability_topic"`
	payloadAvailable    string `json:"payload_available"`
	payloadNotAvailable string `json:"payload_not_available"`
	UnitOfMeasurement   string `json:"unit_of_measurement"`
	UniqueId            string `json:"unique_id"`
	ValueTemplate       string `json:"value_template"`
}

type StateMessages map[string]int

type BatteryLevelStateMessage struct {
	BatteryLevel string `json:"battery"`
}

type RssiStateMessage struct {
	RssiLevel string `json:"rssi"`
}

type Availability bool

const (
	Online  Availability = true
	Offline Availability = false
)

func GetAvailabilityStr(availability Availability) string {
	var onlineStr string
	switch availability {
	case Online:
		onlineStr = "Online"
	case Offline:
		onlineStr = "Offline"
	}
	return onlineStr
}

func GetMessageObjectId(addr string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err.Error())
	}
	uniqueId := reg.ReplaceAllString(addr, "")
	return "inkbird_" + uniqueId
}

func GetMessageStateTopicTemperatures(addr string) string {
	return GetMessageObjectId(addr) + "/temperatures"
}

func GetMessageStateTopicDevice(addr string) string {
	return GetMessageObjectId(addr) + "/device"
}

func GetMessageStateTopicAvailability(addr string) string {
	return GetMessageObjectId(addr) + "/availability"
}

func GetMessageCommandTopicDevice(addr string) string {
	return GetMessageObjectId(addr) + "/command"
}

func NewTemperatureSensorConfigMessage(sensorNum int, addr string) AutoDiscoverConfigMessage {
	temperatureNumbered := "temperature" + strconv.Itoa(sensorNum)

	autoDiscoverConfigMessage := AutoDiscoverConfigMessage{}
	autoDiscoverConfigMessage.sensorNum = sensorNum
	autoDiscoverConfigMessage.addr = addr
	autoDiscoverConfigMessage.id = GetMessageObjectId(addr)
	autoDiscoverConfigMessage.Name = autoDiscoverConfigMessage.id + "_" + temperatureNumbered
	autoDiscoverConfigMessage.UniqueId = autoDiscoverConfigMessage.id + "_" + temperatureNumbered
	autoDiscoverConfigMessage.DeviceClass = "temperature"
	autoDiscoverConfigMessage.StateTopic = GetMessageStateTopicTemperatures(autoDiscoverConfigMessage.addr)
	autoDiscoverConfigMessage.availabilityTopic = GetMessageStateTopicAvailability(autoDiscoverConfigMessage.addr)
	autoDiscoverConfigMessage.payloadAvailable = GetAvailabilityStr(Online)
	autoDiscoverConfigMessage.payloadNotAvailable = GetAvailabilityStr(Offline)
	autoDiscoverConfigMessage.UnitOfMeasurement = "°C"
	autoDiscoverConfigMessage.ValueTemplate = "{{value_json.temperature" + strconv.Itoa(sensorNum) + "}}"
	return autoDiscoverConfigMessage
}

func (autoDiscoverConfigMessage *AutoDiscoverConfigMessage) SetConfigMessageSensorNumber(sensorNum int) {
	temperatureNumbered := "temperature" + strconv.Itoa(sensorNum)

	autoDiscoverConfigMessage.sensorNum = sensorNum
	autoDiscoverConfigMessage.Name = autoDiscoverConfigMessage.id + "_" + temperatureNumbered
	autoDiscoverConfigMessage.UniqueId = autoDiscoverConfigMessage.id + "_" + temperatureNumbered
	autoDiscoverConfigMessage.ValueTemplate = "{{value_json." + temperatureNumbered + "}}"
}

func NewTemperatureSensorBatteryConfigMessage(addr string) AutoDiscoverConfigMessage {
	autoDiscoverConfigMessage := AutoDiscoverConfigMessage{}
	autoDiscoverConfigMessage.addr = addr
	autoDiscoverConfigMessage.id = GetMessageObjectId(addr)
	autoDiscoverConfigMessage.Name = autoDiscoverConfigMessage.id + "_battery"
	autoDiscoverConfigMessage.UniqueId = autoDiscoverConfigMessage.id + "_battery"
	autoDiscoverConfigMessage.DeviceClass = "battery"
	autoDiscoverConfigMessage.StateTopic = GetMessageStateTopicDevice(autoDiscoverConfigMessage.addr)
	autoDiscoverConfigMessage.UnitOfMeasurement = "%"
	autoDiscoverConfigMessage.ValueTemplate = "{{value_json.battery}}"
	return autoDiscoverConfigMessage
}

func NewUnitsSwitchConfigMessage(addr string) AutoDiscoverConfigMessage {
	autoDiscoverConfigMessage := AutoDiscoverConfigMessage{}
	autoDiscoverConfigMessage.addr = addr
	autoDiscoverConfigMessage.id = GetMessageObjectId(addr)
	autoDiscoverConfigMessage.Name = autoDiscoverConfigMessage.id + "_switch"
	autoDiscoverConfigMessage.UniqueId = autoDiscoverConfigMessage.id + "_switch"
	autoDiscoverConfigMessage.StateTopic = GetMessageStateTopicDevice(autoDiscoverConfigMessage.addr)
	autoDiscoverConfigMessage.CommandTopic = GetMessageCommandTopicDevice(autoDiscoverConfigMessage.addr)

	autoDiscoverConfigMessage.ValueTemplate = "{{value_json.switch}}"
	return autoDiscoverConfigMessage
}

func (autoDiscoverConfigMessage *AutoDiscoverConfigMessage) toJson() string {
	j, _ := json.Marshal(autoDiscoverConfigMessage)

	return string(j)
}

func NewBatteryStateMessageJson(batteryLevel int) string {
	batteryLevelStateMessage := BatteryLevelStateMessage{BatteryLevel: strconv.Itoa(batteryLevel)}
	j, _ := json.Marshal(batteryLevelStateMessage)

	return string(j)
}

func NewRssiStateMessageJson(rssi int) string {
	rssiStateMessage := RssiStateMessage{RssiLevel: strconv.Itoa(rssi)}
	j, _ := json.Marshal(rssiStateMessage)

	return string(j)
}

func NewAvailabilityStateMessageJson(availability Availability) string {
	j, _ := json.Marshal(GetAvailabilityStr(availability))

	return string(j)
}

func NewStateMessages(temperatures []int) StateMessages {
	stateMessages := StateMessages{}
	for i, _ := range temperatures {
		stateMessages["temperature"+strconv.Itoa(i+1)] = temperatures[i]
	}
	return stateMessages
}

func (stateMessages StateMessages) toJson() string {
	j, _ := json.Marshal(stateMessages)

	return string(j)
}
