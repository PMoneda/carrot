package rabbitmq

import "github.com/streadway/amqp"

type Publisher struct {
	client  *BrokerClient
	channel *amqp.Channel
}

//Message encapsulate some data configuration
type Message struct {
	Data        []byte
	ContentType string
	Encoding    string
	Headers     map[string]interface{}
}

//Publish a message to exchange in routingkey
func (pub *Publisher) Publish(exchange, routingKey string, message Message) error {

	err := pub.channel.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			Headers:         message.Headers,
			ContentType:     message.ContentType,
			ContentEncoding: message.Encoding,
			Body:            message.Data,
			DeliveryMode:    amqp.Persistent,
			Priority:        0,
		},
	)
	for err != nil {
		pub.channel, err = pub.client.Channel()
		err = pub.channel.Publish(
			exchange,
			routingKey,
			false,
			false,
			amqp.Publishing{
				Headers:         message.Headers,
				ContentType:     message.ContentType,
				ContentEncoding: message.Encoding,
				Body:            message.Data,
				DeliveryMode:    amqp.Persistent,
				Priority:        0,
			},
		)
	}
	return err
}

//NewPublisher creates a new broker publisher
func NewPublisher(client *BrokerClient) (*Publisher, error) {
	pub := new(Publisher)
	pub.client = client
	ch, err := pub.client.client.Channel()
	pub.channel = ch
	return pub, err
}
