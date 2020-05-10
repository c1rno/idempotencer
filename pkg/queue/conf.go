package queue

type BrokerConfig struct {
	InSocket, OutSocket string
}

func (b BrokerConfig) ToLoggerCtx() map[string]interface{} {
	return map[string]interface{}{
		"upstream":   b.InSocket,
		"downstream": b.OutSocket,
	}
}

type ClientConfig struct {
	Socket   string
	Identity string
}

func (c ClientConfig) ToLoggerCtx() map[string]interface{} {
	return map[string]interface{}{
		"identity": c.Identity,
		"destination": c.Socket,
	}
}
