package rabbitmq

import (
	"github.com/PMoneda/carrot"
	rab "github.com/michaelklishin/rabbit-hole"
	"github.com/streadway/amqp"
)

//BrokerClient is a struct to manager api and ampq connection
type BrokerClient struct {
	api    *rab.Client
	client *amqp.Connection
	config carrot.ConnectionConfig
}

func (broker *BrokerClient) connectoToAmqp() (err error) {
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

//NewBrokerClient creates a new rabbit broker client
func NewBrokerClient(config carrot.ConnectionConfig) (client *BrokerClient, err error) {
	client = new(BrokerClient)
	client.config = config
	err = client.connectoToAPI()
	if err != nil {
		return
	}
	err = client.connectoToAmqp()
	return
}
