package queue

type BrokerConfig struct {
	InSocket, OutSocket string
}

type ConsumerConfig struct {
	Socket string
}

type ProducerConfig struct {
	Socket string
}
