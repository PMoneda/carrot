package rabbitmq

import (
	"testing"

	"github.com/PMoneda/carrot"

	. "github.com/smartystreets/goconvey/convey"
)

func TestShouldConnectToRabbit(t *testing.T) {
	Convey("should connect to rabbitmq api", t, func() {
		connConfig := &carrot.ConnectionConfig{
			Host:     "localhost",
			Username: "guest",
			Password: "guest",
		}
		client, err := NewBrokerClient(connConfig)
		So(err, ShouldBeNil)
		f, err := client.api.ListVhosts()
		So(err, ShouldBeNil)
		So(len(f), ShouldBeGreaterThan, 0)
	})

	Convey("should not connect to rabbitmq api", t, func() {
		connConfig := &carrot.ConnectionConfig{
			Host:     "localhos",
			Username: "guest",
			Password: "guest",
		}
		_, err := NewBrokerClient(connConfig)
		So(err, ShouldNotBeNil)

	})

	Convey("should not create api client when url cannot be parsed", t, func() {
		connConfig := &carrot.ConnectionConfig{
			Host:     "localh*& os",
			Username: "guest",
			Password: "guest",
		}
		_, err := NewBrokerClient(connConfig)
		So(err, ShouldNotBeNil)
	})
}
