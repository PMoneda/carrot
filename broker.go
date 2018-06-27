package carrot

type Broker interface {
}

type BrokerConfig interface {
}

type BrokerConnection interface {
	Connect(url, username, password string) error

	Setup(config BrokerConfig) error
}
