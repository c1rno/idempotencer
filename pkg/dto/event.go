package dto

import (
	"bytes"
)

type Msg interface {
	Data() [][]byte
	String() string
}

func NewStringMsg(s ...string) Msg {
	d := make([][]byte, 0, len(s))
	for _, chunk := range s {
		d = append(d, []byte(chunk))
	}
	return raw(d)
}

func NewByteMsg(b ...[]byte) Msg {
	return raw(b)
}

type raw [][]byte

func (p raw) Data() [][]byte {
	return p
}

func (p raw) String() string {
	return string(bytes.Join(p, []byte(" ")))
}
