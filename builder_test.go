package carrot

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestShouldBuildRabbitInfra(t *testing.T) {
	config := ConnectionConfig{
		Host:     "localhost",
		Username: "guest",
		Password: "guest",
	}
	conn, _ := NewBrokerClient(&config)
	builder := NewBuilder(conn)
	builder.client.api.DeleteVhost("test_carrot_v1")
	Convey("should build RabbitMq Infra", t, func() {

		Convey("should create vhost on rabbitmq", func() {
			So(builder.UseVHost("test_carrot_v1"), ShouldBeNil)
			So(builder.client.config.VHost, ShouldEqual, "test_carrot_v1")
		})

		Convey("should create exchange on rabbitmq", func() {
			config.VHost = "test_carrot_v1"
			conn, err := NewBrokerClient(&config)
			So(err, ShouldBeNil)
			builder := NewBuilder(conn)
			So(builder.DeclareTopicExchange("my_topic_exchange"), ShouldBeNil)
			So(builder.DeclareDirectExchange("my_direct_exchange"), ShouldBeNil)
			So(builder.DeclareFanoutExchange("my_fanout_exchange"), ShouldBeNil)
			So(builder.DeclareHeadersExchange("my_headers_exchange"), ShouldBeNil)
		})

		Convey("should create queue on rabbitmq", func() {
			config.VHost = "test_carrot_v1"
			conn, err := NewBrokerClient(&config)
			So(err, ShouldBeNil)
			builder := NewBuilder(conn)
			So(builder.DeclareQueue("my_queue"), ShouldBeNil)
			info, err := builder.client.api.GetQueue(config.VHost, "my_queue")
			So(err, ShouldBeNil)
			So(info.Durable, ShouldBeTrue)
		})

		Convey("should bind queue to exchange on rabbitmq", func() {
			config.VHost = "test_carrot_v1"
			conn, err := NewBrokerClient(&config)
			So(err, ShouldBeNil)
			builder := NewBuilder(conn)
			So(builder.DeclareQueue("my_queue"), ShouldBeNil)
			So(builder.DeclareTopicExchange("my_topic_exchange"), ShouldBeNil)
			So(builder.BindQueueToExchange("my_queue", "my_topic_exchange", ".*"), ShouldBeNil)
			info, err := builder.client.api.ListQueueBindings(config.VHost, "my_queue")
			So(err, ShouldBeNil)
			So(len(info), ShouldBeGreaterThan, 1)
		})

		Convey("should create topic permission on rabbitmq", func() {
			config.VHost = "test_carrot_v1"
			conn, err := NewBrokerClient(&config)
			So(err, ShouldBeNil)
			builder := NewBuilder(conn)
			So(builder.DeclareTopicExchange("my_topic_exchange"), ShouldBeNil)
			So(builder.UpdateTopicPermission("guest", "my_topic_exchange"), ShouldBeNil)
			list, err := builder.client.api.ListTopicPermissionsOf("guest")
			So(err, ShouldBeNil)
			So(len(list), ShouldBeGreaterThan, 0)
		})
	})

}
