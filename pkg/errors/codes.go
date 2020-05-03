package errors

var errorsMap = map[int]errMapVal{
	UnknownError:           {msg: "Unknown error", fatal: true},
	ConfigReadingError:     {msg: "Can't read config"},
	ConfigUnmarshallError:  {msg: "Building config error", fatal: true},
	MetricsSrvFail:         {msg: "Metrics server failed"},
	MetricsSrvShutdownFail: {msg: "Metrics server shutdown failed"},
}

const (
	UnknownError int = iota // is't ZERO, default value, means we don't set `code` into BaseError
	ConfigReadingError
	ConfigUnmarshallError
	MetricsSrvFail
	MetricsSrvShutdownFail
)

type errMapVal struct {
	msg   string
	fatal bool
}
