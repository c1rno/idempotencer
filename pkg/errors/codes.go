package errors

var errorsMap = map[int]errMapVal{
	UnknownError:                {msg: "Unknown error", fatal: true},
	ConfigReadingError:          {msg: "Can't read config"},
	ConfigUnmarshallError:       {msg: "Building config error", fatal: true},
	MetricsSrvFail:              {msg: "Metrics server failed"},
	MetricsSrvShutdownFail:      {msg: "Metrics server shutdown failed"},
	NewRouterSocketCreationFail: {msg: "Can't create new ROUTER socket"},
	BindSocketFail:              {msg: "Can't bind socket"},
	CloseSocketFail:             {msg: "Can't close socket"},
	NewReqSocketCreationFail:    {msg: "Can't create new REQ socket"},
	ConnectSocketFail:           {msg: "Can't connect socket"},
	PullSocketError:             {msg: "Pull socket error", fatal: true},
	PushSocketError:             {msg: "Push socket error"},
	LoadBalancerError:           {msg: "LoadBalancer start error"},
	InvalidConfiguration:        {msg: "Required configuration is missed or invalid"},
	PullSocketNotReadyError:     {msg: "Pull socket error, try later"},
}

const (
	UnknownError int = iota // is't ZERO, default value, means we don't set `code` into BaseError
	ConfigReadingError
	ConfigUnmarshallError
	MetricsSrvFail
	MetricsSrvShutdownFail
	NewRouterSocketCreationFail
	BindSocketFail
	CloseSocketFail
	NewReqSocketCreationFail
	ConnectSocketFail
	PullSocketError
	PushSocketError
	LoadBalancerError
	InvalidConfiguration
	PullSocketNotReadyError
)

type errMapVal struct {
	msg   string
	fatal bool
}
