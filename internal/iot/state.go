package iot

import (
	"encoding/json"
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type StateReceiver struct{}

type State struct {
	Payload []byte
}

func (r *StateReceiver) stateHandler(stateChan chan<- *State) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		iotShadowState := new(IotShadowState)
		err := json.Unmarshal(msg.Payload(), iotShadowState)
		if err != nil {
			log.Println("[shadowHandler] error occurred in unmarshalling shadow payload", err)
			return
		}
		if iotShadowState.ShadowState != nil && iotShadowState.ShadowState.Reported != nil {
			stateUpdate, err := json.Marshal(iotShadowState.ShadowState.Reported)
			if err != nil {
				log.Println("[shadowHandler] error occurred in marshalling reported shadow payload")
				return
			}
			go r.broadCastState(stateChan, stateUpdate)
		}
		fmt.Printf("MSG: %s\n", msg.Payload())
	}
}

func (r *StateReceiver) broadCastState(stateChan chan<- *State, payload []byte) {
	stateChan <- &State{Payload: payload}
}
