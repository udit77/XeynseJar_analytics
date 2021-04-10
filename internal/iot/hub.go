package iot

import (
	"encoding/json"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/xeynse/XeynseJar_analytics/internal/entity"
	"github.com/xeynse/XeynseJar_analytics/internal/repo/jar"
)

type Resource interface {
	SubscribeToJarStateTopic() error
	Run()
}

const (
	JarSubscribeTopic = "xeynse/smarJar/analytics/stats"
)

type IOTHub struct {
	mqtt            mqtt.Client
	JarStateUpdater chan *State
	jarDBRepo       jar.Repo
}

func New(db jar.Repo) (Resource, error) {
	iotHub := &IOTHub{
		JarStateUpdater: make(chan *State),
		jarDBRepo:       db,
	}
	receiver := new(StateReceiver)
	iotHub.mqtt = InitializeMqTT(receiver.stateHandler(iotHub.JarStateUpdater))
	err := iotHub.SubscribeToJarStateTopic()
	return iotHub, err
}

func (hub *IOTHub) SubscribeToJarStateTopic() error {
	//TODO set mqtt wait timeout
	if token := hub.mqtt.Subscribe(JarSubscribeTopic, 0, nil); token.Wait() && token.Error() != nil {
		log.Printf("failed to create subscription for: %v, %v\n", JarSubscribeTopic, token.Error())
		return token.Error()
	}
	return nil
}

func (hub *IOTHub) Run() {
	for {
		select {
		case state := <-hub.JarStateUpdater:
			jarState := make(map[string]entity.JarState)
			err := json.Unmarshal(state.Payload, &jarState)
			if err != nil {
				log.Println("[Run] error occurred in unmarshalling jar payload :", err)
				return
			}
			err = hub.jarDBRepo.InsertJarStateData(jarState)
			if err != nil {
				log.Println("[Run] error occurred in inserting jar stat to database :", err)
				return
			}
		}
	}
}
