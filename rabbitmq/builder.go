package rabbitmq

import (
	rab "github.com/michaelklishin/rabbit-hole"
)

//Builder builds a rabbit infrastructure
type Builder struct {
	client *BrokerClient
}

//UseVHost creates or use existing rabbit vhost
func (builder *Builder) UseVHost(vhost string) error {
	_, err := builder.client.api.PutVhost(vhost, rab.VhostSettings{Tracing: false})
	if err == nil {
		builder.client.config.VHost = vhost
		builder.client.connectoToAPI()
	}
	return err
}

//DeclareTopicExchange create a durable topic exchange
func (builder *Builder) DeclareTopicExchange(exchange string) error {
	_, err := builder.client.api.DeclareExchange(builder.client.config.VHost, exchange, rab.ExchangeSettings{
		Durable: true,
		Type:    "topic",
	})
	return err
}

//DeclareDirectExchange create a durable direct exchange
func (builder *Builder) DeclareDirectExchange(exchange string) error {
	_, err := builder.client.api.DeclareExchange(builder.client.config.VHost, exchange, rab.ExchangeSettings{
		Durable: true,
		Type:    "direct",
	})
	return err
}

//DeclareHeadersExchange create a durable headers exchange
func (builder *Builder) DeclareHeadersExchange(exchange string) error {
	_, err := builder.client.api.DeclareExchange(builder.client.config.VHost, exchange, rab.ExchangeSettings{
		Durable: true,
		Type:    "headers",
	})
	return err
}

//DeclareFanoutExchange create a durable fanout exchange
func (builder *Builder) DeclareFanoutExchange(exchange string) error {
	_, err := builder.client.api.DeclareExchange(builder.client.config.VHost, exchange, rab.ExchangeSettings{
		Durable: true,
		Type:    "fanout",
	})
	return err
}

//DeclareQueue creates a durable queue
func (builder *Builder) DeclareQueue(queue string) error {
	_, err := builder.client.api.DeclareQueue(builder.client.config.VHost, queue, rab.QueueSettings{
		Durable: true,
	})
	return err
}

//BindQueueToExchange binds a queue to a exchange
func (builder *Builder) BindQueueToExchange(queue, exchange, routingKey string) error {
	_, err := builder.client.api.DeclareBinding(builder.client.config.VHost, rab.BindingInfo{
		RoutingKey:      routingKey,
		Source:          exchange,
		Destination:     queue,
		DestinationType: "queue",
	})
	return err
}

//NewBuilder creates new broker builder
func NewBuilder(conn *BrokerClient) *Builder {
	builder := new(Builder)
	builder.client = conn
	return builder
}
