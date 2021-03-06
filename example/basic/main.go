package main

import (
	"fmt"
	"time"

	"github.com/PMoneda/carrot"
)

func main() {
	config := carrot.ConnectionConfig{
		Host:     "localhost",
		Username: "guest",
		Password: "guest",
	}
	conn, _ := carrot.NewBrokerClient(&config)
	builder := carrot.NewBuilder(conn)
	exchange := "p_exchange"
	queue := "publisher_queue"
	builder.UseVHost("test_carrot_v1.1")
	builder.DeclareTopicExchange(exchange)
	builder.DeclareQueue(queue)
	builder.UpdateTopicPermission("guest", exchange)
	builder.BindQueueToExchange(queue, exchange, "*")

	subConn, _ := carrot.NewBrokerClient(&config)
	subscriber := carrot.NewSubscriber(subConn)
	subscriber.SetMaxRetries(30)
	subscriber.Subscribe(carrot.SubscribeWorker{
		Queue: queue,
		Scale: 1,
		Handler: func(msg *carrot.MessageContext) error {
			return msg.Ack()
		},
	})
	pub := carrot.NewPublisher(conn)
	go func() {
		for {

			err := pub.Publish(exchange, "*", carrot.Message{ContentType: "application/json", Encoding: "utf-8", Data: []byte("hello")})
			if err != nil {
				fmt.Println(err)
			}
			time.Sleep(1 * time.Second)
		}
	}()
	ch := make(chan int)
	<-ch
}
