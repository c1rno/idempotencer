package logging

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
)

type Logger interface {
	Debug(string, map[string]interface{})
	Info(string, map[string]interface{})
	Error(string, map[string]interface{})
}

// TODO: redo?
// lvl = 0 is debug
// lvl = 1 is info
// lvl = 2 is error
// lvl = 3 is null
func NewLogger(lvl int) Logger {
	if stdLogger == nil {
		stdLogger = log.New(os.Stderr, "", log.LstdFlags)
	}
	return logger{
		lvl: lvl,
		l:   stdLogger,
	}
}

var stdLogger *log.Logger = nil

var _ Logger = logger{}

type logger struct {
	lvl int
	l   *log.Logger
}

func (l logger) Debug(msg string, data map[string]interface{}) {
	if l.lvl < 1 {
		l.print(msg, data)
	}
}

func (l logger) Info(msg string, data map[string]interface{}) {
	if l.lvl < 2 {
		l.print(msg, data)
	}
}

func (l logger) Error(msg string, data map[string]interface{}) {
	if l.lvl < 3 {
		l.print(msg, data)
	}
}

func (l logger) print(msg string, data map[string]interface{}) {
	if data != nil {
		buffer := &bytes.Buffer{}
		encoder := json.NewEncoder(buffer)
		encoder.SetEscapeHTML(false)
		_ = encoder.Encode(data)
		l.l.Printf("%s %s", msg, buffer.String())
	} else {
		l.l.Print(msg)
	}
}
