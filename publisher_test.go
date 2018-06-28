package carrot

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRabbitPublish(t *testing.T) {
	config := ConnectionConfig{
		Host:     "localhost",
		Username: "guest",
		Password: "guest",
	}
	conn, _ := NewBrokerClient(&config)
	builder := NewBuilder(conn)
	exchange := "p_exchange"
	queue := "publisher_queue"
	builder.DeclareTopicExchange(exchange)
	builder.DeclareQueue(queue)
	builder.UseVHost("test_carrot_v1.1")
	builder.UpdateTopicPermission("guest", exchange)
	builder.BindQueueToExchange(queue, exchange, "*")

	Convey("should publish message to rabbitmq", t, func() {
		pub := NewPublisher(conn)

		Convey("should create vhost on rabbitmq", func() {

			errPublish := pub.Publish(exchange, "*", Message{ContentType: "application/json", Encoding: "utf-8", Data: []byte("hello")})
			ch, errChannel := pub.client.Channel()
			delivery, errConsume := ch.Consume(queue, "", false, false, false, false, nil)
			msg := <-delivery

			So(msg.Ack(false), ShouldBeNil)
			So(errPublish, ShouldBeNil)
			So(errChannel, ShouldBeNil)
			So(errConsume, ShouldBeNil)
			ch.Close()
		})

		Convey("should reconnect to a channel when it is closed", func() {
			pub.client.client.Close()
			err := pub.Publish(exchange, "*", Message{ContentType: "application/json", Encoding: "utf-8", Data: []byte("hello")})
			So(err, ShouldBeNil)
		})
	})
}
