package iot

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type StateReceiver struct{}

type State struct {
	Payload []byte
}

func (r *StateReceiver) stateHandler(stateChan chan<- *State) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		r.broadCastState(stateChan, msg.Payload())
		fmt.Printf("TOPIC: %s\n", msg.Topic())
		fmt.Printf("MSG: %s\n", msg.Payload())
	}
}

func (r *StateReceiver) broadCastState(stateChan chan<- *State, payload []byte) {
	stateChan <- &State{Payload: payload}
}
