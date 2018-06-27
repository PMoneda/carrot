package carrot

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestShouldReturnAMQPURL(t *testing.T) {
	cnf := ConnectionConfig{
		Host:     "localhost",
		VHost:    "myvhost",
		Username: "guest",
		Password: "guest",
	}
	Convey("should return ampq host pattern", t, func() {
		So(cnf.GetAMQPURI(), ShouldEqual, "amqp://guest:guest@localhost:5672/myvhost")
	})

	Convey("should return current amqp port config", t, func() {
		cnf.AMQPPort = "9090"
		So(cnf.GetAMQPPort(), ShouldEqual, "9090")
	})

	Convey("should return api host pattern", t, func() {
		So(cnf.GetAPIURI(), ShouldEqual, "http://localhost:15672/")
	})

	Convey("should return current api port config", t, func() {
		cnf.APIPort = "9090"
		So(cnf.GetAPIPort(), ShouldEqual, "9090")
	})

}
