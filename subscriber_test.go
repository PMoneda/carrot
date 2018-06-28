package carrot

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestShouldSubscribeOnRabbit(t *testing.T) {
	config := ConnectionConfig{
		Host:     "localhost",
		Username: "guest",
		Password: "guest",
	}
	conn, _ := NewBrokerClient(&config)
	builder := NewBuilder(conn)
	exchange := "p_exchange"
	queue := "publisher_queue"
	builder.UseVHost("test_carrot_v1.1")
	builder.DeclareTopicExchange(exchange)
	builder.DeclareQueue(queue)
	builder.DeclareTopicExchange("error_ex")
	builder.DeclareQueue("error_q")
	builder.UpdateTopicPermission("guest", exchange)
	builder.UpdateTopicPermission("guest", "error_ex")
	builder.BindQueueToExchange(queue, exchange, "*")
	builder.BindQueueToExchange("error_q", "error_ex", "*")

	pub := NewPublisher(conn)

	Convey("should test subscriber", t, func() {

		Convey("should bind to a queue", func() {
			wait := make(chan bool)
			subConn, _ := NewBrokerClient(&config)
			subscriber := NewSubscriber(subConn)
			subscriber.Subscribe(SubscribeWorker{
				Queue: queue,
				Scale: 1,
				Handler: func(msg *MessageContext) error {
					ok := string(msg.Message.Data) == "hello"
					if !ok {
						t.Fail()
					}
					msg.Ack()
					wait <- ok
					return nil
				},
			})

			pub.Publish(exchange, "*", Message{ContentType: "application/json", Encoding: "utf-8", Data: []byte("hello")})
			<-wait
			subscriber.client.client.Close()
		})

		Convey("should redirect message to another queue", func() {
			wait := make(chan bool)
			subConn, _ := NewBrokerClient(&config)
			subscriber := NewSubscriber(subConn)
			subscriber.Subscribe(SubscribeWorker{
				Queue: queue,
				Scale: 1,
				Handler: func(msg *MessageContext) error {
					err := msg.RedirectTo("error_ex", "*")
					if err != nil {
						t.Fail()
					}
					msg.Ack()
					return nil
				},
			})

			subscriber.Subscribe(SubscribeWorker{
				Queue: "error_q",
				Scale: 1,
				Handler: func(msg *MessageContext) error {
					msg.Ack()
					time.Sleep(10 * time.Millisecond)
					wait <- true
					return nil
				},
			})
			pub.Publish(exchange, "*", Message{ContentType: "application/json", Encoding: "utf-8", Data: []byte("hello")})
			<-wait
			subscriber.client.client.Close()
		})

		Convey("should nack message", func() {
			wait := make(chan bool)
			subConn, _ := NewBrokerClient(&config)
			subscriber := NewSubscriber(subConn)
			subscriber.Subscribe(SubscribeWorker{
				Queue: queue,
				Scale: 1,
				Handler: func(msg *MessageContext) error {
					if err := msg.Nack(true); err != nil {
						t.Fail()
					}
					wait <- true
					return nil
				},
			})
			pub.Publish(exchange, "*", Message{ContentType: "application/json", Encoding: "utf-8", Data: []byte("hello")})
			<-wait
		})
	})
}
