package carrot

import (
	"fmt"

	rab "github.com/michaelklishin/rabbit-hole"
	"github.com/streadway/amqp"
)

//BrokerClient is a struct to manager api and ampq connection
type BrokerClient struct {
	api     *rab.Client
	client  *amqp.Connection
	config  *ConnectionConfig
	channel *amqp.Channel
}

func (broker *BrokerClient) connectoToAmqp() (err error) {
	if broker.client != nil {
		broker.client.Close()
	}
	broker.client, err = amqp.Dial(broker.config.GetAMQPURI())
	return
}

func (broker *BrokerClient) connectoToAPI() (err error) {
	broker.api, err = rab.NewClient(broker.config.GetAPIURI(), broker.config.Username, broker.config.Password)
	if err != nil {
		return
	}
	return
}

//Channel return amqp channel with reconnect capabilities
func (broker *BrokerClient) Channel() (ch *amqp.Channel, err error) {
	if broker.channel != nil {
		//TODO refactor to use NotifyClose
		errClose := broker.channel.Close()
		if errClose != nil {
			fmt.Println(errClose)
		}
	}
	ch, err = broker.client.Channel()
	if err == nil {
		broker.channel = ch
	}
	return
}

//NewBrokerClient creates a new rabbit broker client
func NewBrokerClient(config *ConnectionConfig) (client *BrokerClient, err error) {
	client = new(BrokerClient)
	client.config = config
	err = client.connectoToAPI()
	if err != nil {
		return
	}
	err = client.connectoToAmqp()
	return
}
