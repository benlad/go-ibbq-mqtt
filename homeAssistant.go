package main

import (
	"encoding/json"
	"regexp"
	"strconv"

	log "github.com/mgutz/logxi/v1"
)

// The capitalized filed names are exported and thus end up in JSON
type AutoDiscoverSensorConfigMessage struct {
	addr                string
	id                  string
	sensorNum           int
	Name                string `json:"name"`
	DeviceClass         string `json:"device_class"`
	StateTopic          string `json:"state_topic"`
	availabilityTopic   string `json:"availability_topic"`
	payloadAvailable    string `json:"payload_available"`
	payloadNotAvailable string `json:"payload_not_available"`
	UnitOfMeasurement   string `json:"unit_of_measurement"`
	UniqueId            string `json:"unique_id"`
	ValueTemplate       string `json:"value_template"`
}

type AutoDiscoverSwitchConfigMessage struct {
	addr                string
	id                  string
	sensorNum           int
	Name                string `json:"name"`
	StateTopic          string `json:"state_topic"`
	CommandTopic        string `json:"command_topic"`
	PayloadOn           string `json:"payload_on"`
	PayloadOff          string `json:"payload_off"`
	availabilityTopic   string `json:"availability_topic"`
	payloadAvailable    string `json:"payload_available"`
	payloadNotAvailable string `json:"payload_not_available"`
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

type TemperatureUnits string

const (
	Celsius    TemperatureUnits = "C"
	Fahrenheit TemperatureUnits = "F"
)

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

func GetMessageUnitsCommandTopicDevice(addr string) string {
	return GetMessageObjectId(addr) + "/units"
}

func NewTemperatureSensorConfigMessage(sensorNum int, addr string) AutoDiscoverSensorConfigMessage {
	temperatureNumbered := "temperature" + strconv.Itoa(sensorNum)

	autoDiscoverSensorConfigMessage := AutoDiscoverSensorConfigMessage{}
	autoDiscoverSensorConfigMessage.sensorNum = sensorNum
	autoDiscoverSensorConfigMessage.addr = addr
	autoDiscoverSensorConfigMessage.id = GetMessageObjectId(addr)
	autoDiscoverSensorConfigMessage.Name = autoDiscoverSensorConfigMessage.id + "_" + temperatureNumbered
	autoDiscoverSensorConfigMessage.UniqueId = autoDiscoverSensorConfigMessage.id + "_" + temperatureNumbered
	autoDiscoverSensorConfigMessage.DeviceClass = "temperature"
	autoDiscoverSensorConfigMessage.StateTopic = GetMessageStateTopicTemperatures(autoDiscoverSensorConfigMessage.addr)
	autoDiscoverSensorConfigMessage.availabilityTopic = GetMessageStateTopicAvailability(autoDiscoverSensorConfigMessage.addr)
	autoDiscoverSensorConfigMessage.payloadAvailable = GetAvailabilityStr(Online)
	autoDiscoverSensorConfigMessage.payloadNotAvailable = GetAvailabilityStr(Offline)
	autoDiscoverSensorConfigMessage.UnitOfMeasurement = "°C"
	autoDiscoverSensorConfigMessage.ValueTemplate = "{{value_json.temperature" + strconv.Itoa(sensorNum) + "}}"
	return autoDiscoverSensorConfigMessage
}

func (autoDiscoverSensorConfigMessage *AutoDiscoverSensorConfigMessage) SetConfigMessageSensorNumber(sensorNum int) {
	temperatureNumbered := "temperature" + strconv.Itoa(sensorNum)

	autoDiscoverSensorConfigMessage.sensorNum = sensorNum
	autoDiscoverSensorConfigMessage.Name = autoDiscoverSensorConfigMessage.id + "_" + temperatureNumbered
	autoDiscoverSensorConfigMessage.UniqueId = autoDiscoverSensorConfigMessage.id + "_" + temperatureNumbered
	autoDiscoverSensorConfigMessage.ValueTemplate = "{{value_json." + temperatureNumbered + "}}"
}

func NewTemperatureSensorBatteryConfigMessage(addr string) AutoDiscoverSensorConfigMessage {
	autoDiscoverSensorConfigMessage := AutoDiscoverSensorConfigMessage{}
	autoDiscoverSensorConfigMessage.addr = addr
	autoDiscoverSensorConfigMessage.id = GetMessageObjectId(addr)
	autoDiscoverSensorConfigMessage.Name = autoDiscoverSensorConfigMessage.id + "_battery"
	autoDiscoverSensorConfigMessage.UniqueId = autoDiscoverSensorConfigMessage.id + "_battery"
	autoDiscoverSensorConfigMessage.DeviceClass = "battery"
	autoDiscoverSensorConfigMessage.StateTopic = GetMessageStateTopicDevice(autoDiscoverSensorConfigMessage.addr)
	autoDiscoverSensorConfigMessage.UnitOfMeasurement = "%"
	autoDiscoverSensorConfigMessage.ValueTemplate = "{{value_json.battery}}"
	return autoDiscoverSensorConfigMessage
}

func NewUnitsSwitchConfigMessage(addr string) AutoDiscoverSwitchConfigMessage {
	autoDiscoverSwitchConfigMessage := AutoDiscoverSwitchConfigMessage{}
	autoDiscoverSwitchConfigMessage.addr = addr
	autoDiscoverSwitchConfigMessage.id = GetMessageObjectId(addr)
	autoDiscoverSwitchConfigMessage.Name = autoDiscoverSwitchConfigMessage.id + "_switch"
	autoDiscoverSwitchConfigMessage.StateTopic = GetMessageUnitsCommandTopicDevice(autoDiscoverSwitchConfigMessage.addr)
	autoDiscoverSwitchConfigMessage.CommandTopic = GetMessageUnitsCommandTopicDevice(autoDiscoverSwitchConfigMessage.addr)
	autoDiscoverSwitchConfigMessage.PayloadOn = "F"
	autoDiscoverSwitchConfigMessage.PayloadOff = "C"

	return autoDiscoverSwitchConfigMessage
}

func (autoDiscoverSensorConfigMessage *AutoDiscoverSensorConfigMessage) toJson() string {
	j, _ := json.Marshal(autoDiscoverSensorConfigMessage)

	return string(j)
}

func (autoDiscoverSwitchConfigMessage *AutoDiscoverSwitchConfigMessage) toJson() string {
	j, _ := json.Marshal(autoDiscoverSwitchConfigMessage)

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

func NewUnitsSwitchCommandMessageJson(temperatureUnits TemperatureUnits) string {
	j, _ := json.Marshal(temperatureUnits)

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
