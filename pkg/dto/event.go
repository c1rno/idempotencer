package dto

type Msg interface {
	// over one frame sending it?
	Data() []string
}

func NewRawMsg(s ...string) Msg {
	return raw(s)
}

type raw []string

func (p raw) Data() []string {
	return p
}
