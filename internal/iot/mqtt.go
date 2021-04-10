package iot

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/xeynse/XeynseJar_analytics/internal/config"
	"github.com/xeynse/XeynseJar_analytics/internal/util"
)

func InitializeMqTT(defaultHandler mqtt.MessageHandler) mqtt.Client {
	tlsConfig, err := NewTLSConfig()
	if err != nil {
		log.Fatalf("failed to create TLS configuration: %v", err)
	}

	//Should not be more than 23
	clientID := fmt.Sprintf("%v-%v", "xeynseJar", util.GetRandom(12))

	opts := mqtt.NewClientOptions()
	opts.AddBroker(config.GetConfig().MQTT.Broker)
	opts.SetClientID(clientID)
	opts.SetTLSConfig(tlsConfig)
	opts.SetCleanSession(false)
	opts.SetMaxReconnectInterval(1 * time.Second)
	opts.SetDefaultPublishHandler(defaultHandler)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("failed to create connection: %v", token.Error())
	}

	return c
}

func NewTLSConfig() (tslConfig *tls.Config, err error) {
	certificatePool := x509.NewCertPool()
	pemCerts, err := ioutil.ReadFile(config.GetConfig().Certs.Root)
	if err != nil {
		return
	}
	certificatePool.AppendCertsFromPEM(pemCerts)
	cert, err := tls.LoadX509KeyPair(config.GetConfig().Certs.Cert, config.GetConfig().Certs.Private)
	if err != nil {
		return
	}
	tslConfig = &tls.Config{
		RootCAs:      certificatePool,
		ClientAuth:   tls.NoClientCert,
		ClientCAs:    nil,
		Certificates: []tls.Certificate{cert},
	}
	return
}
