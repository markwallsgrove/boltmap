package blitzortung

import (
	"fmt"
	"math/rand"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const topic = "blitzortung/1.1/#"

type Client struct {
	mqttClient mqtt.Client
	strikes    chan Strike
}

func NewClient() *Client {
	return &Client{
		strikes: make(chan Strike, 256),
	}
}

func (c *Client) Connect(broker string) error {
	opts := mqtt.NewClientOptions().
		AddBroker(broker).
		SetClientID(uniqueClientID()).
		SetCleanSession(true).
		SetAutoReconnect(false)

	opts.SetOnConnectHandler(func(client mqtt.Client) {
		client.Subscribe(topic, 0, func(_ mqtt.Client, msg mqtt.Message) {
			s, err := ParseStrike(msg.Payload())
			if err != nil {
				return
			}
			select {
			case c.strikes <- s:
			default:
			}
		})
	})

	c.mqttClient = mqtt.NewClient(opts)
	token := c.mqttClient.Connect()
	token.Wait()
	if err := token.Error(); err != nil {
		return fmt.Errorf("connect to %s: %w", broker, err)
	}
	return nil
}

func (c *Client) Strikes() <-chan Strike {
	return c.strikes
}

func uniqueClientID() string {
	return fmt.Sprintf("boltmap-%d", rand.Int63())
}

func (c *Client) Disconnect() {
	if c.mqttClient != nil && c.mqttClient.IsConnected() {
		c.mqttClient.Disconnect(250)
	}
}
