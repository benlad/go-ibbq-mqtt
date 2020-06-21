package main

import (
	"encoding/json"
	"regexp"
	"strconv"

	log "github.com/mgutz/logxi/v1"
)

// The capitalized filed names are exported and thus end up in JSON
type AutoDiscoverConfigMessage struct {
	id                string
	sensorNum         int
	Name              string `json:"name"`
	DeviceClass       string `json:"device_class"`
	StateTopic        string `json:"state_topic"`
	UnitOfMeasurement string `json:"unit_of_measurement"`
	UniqueId          string `json:"unique_id"`
	ValueTemplate     string `json:"value_template"`
}

type StateMessage struct {
	Temperature float64 `json:"temperature"`
}
type StateMessages map[string]float64

func GetMessageObjectId(id string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err.Error())
	}
	uniqueId := reg.ReplaceAllString(id, "")
	return "inkbird_" + uniqueId
}

func GetMessageStateTopic(id string) string {
	return GetMessageObjectId(id) + "/state"
}

func NewTemperatureSensorConfigMessage(sensorNum int, id string) AutoDiscoverConfigMessage {
	temperatureNumbered := "temperature" + strconv.Itoa(sensorNum)

	autoDiscoverConfigMessage := AutoDiscoverConfigMessage{}
	autoDiscoverConfigMessage.sensorNum = sensorNum
	autoDiscoverConfigMessage.id = GetMessageObjectId(id)
	autoDiscoverConfigMessage.Name = autoDiscoverConfigMessage.id + "_" + temperatureNumbered
	autoDiscoverConfigMessage.DeviceClass = "temperature"
	autoDiscoverConfigMessage.StateTopic = autoDiscoverConfigMessage.id
	autoDiscoverConfigMessage.UnitOfMeasurement = "°C"
	autoDiscoverConfigMessage.UniqueId = autoDiscoverConfigMessage.id + "_" + temperatureNumbered
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

func (autoDiscoverConfigMessage *AutoDiscoverConfigMessage) toJson() string {
	j, _ := json.Marshal(autoDiscoverConfigMessage)

	return string(j)
}

// func NewStateMessages(temperatures []float64) StateMessages {
// 	//	stateMessage := StateMessage{}
// 	stateMessages := StateMessages{}
// 	for i, _ := range temperatures {
// 		//		stateMessage.Temperature = temperatures[i]
// 		stateMessages = append(stateMessages, StateMessage{Temperature: temperatures[i]})
// 	}
// 	return stateMessages
// }

func NewStateMessages(temperatures []float64) StateMessages {
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
