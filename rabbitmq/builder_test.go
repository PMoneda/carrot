package rabbitmq

import (
	"testing"

	"github.com/PMoneda/carrot"

	. "github.com/smartystreets/goconvey/convey"
)

func TestShouldBuildRabbitInfra(t *testing.T) {
	config := carrot.ConnectionConfig{
		Host:     "localhost",
		Username: "guest",
		Password: "guest",
	}

	Convey("should create vhost on rabbitmq", t, func() {
		conn, err := NewBrokerClient(&config)
		So(err, ShouldBeNil)
		builder := NewBuilder(conn)
		So(builder.UseVHost("test_carrot_v1"), ShouldBeNil)
		So(builder.client.config.VHost, ShouldEqual, "test_carrot_v1")

	})

	Convey("should create exchange on rabbitmq", t, func() {
		config.VHost = "test_carrot_v1"
		conn, err := NewBrokerClient(&config)
		So(err, ShouldBeNil)
		builder := NewBuilder(conn)
		So(builder.DeclareTopicExchange("my_topic_exchange"), ShouldBeNil)
		So(builder.DeclareDirectExchange("my_direct_exchange"), ShouldBeNil)
		So(builder.DeclareFanoutExchange("my_fanout_exchange"), ShouldBeNil)
		So(builder.DeclareHeadersExchange("my_headers_exchange"), ShouldBeNil)
	})

	Convey("should create queue on rabbitmq", t, func() {
		config.VHost = "test_carrot_v1"
		conn, err := NewBrokerClient(&config)
		So(err, ShouldBeNil)
		builder := NewBuilder(conn)
		So(builder.DeclareQueue("my_queue"), ShouldBeNil)
		info, err := builder.client.api.GetQueue(config.VHost, "my_queue")
		So(err, ShouldBeNil)
		So(info.Durable, ShouldBeTrue)
	})
}
