package queue

import (
	"bytes"
	"encoding/json"
	stdErrors "errors"

	"github.com/c1rno/idempotencer/pkg/dto"
	"github.com/c1rno/idempotencer/pkg/errors"
	"github.com/c1rno/idempotencer/pkg/helpers"
	"github.com/c1rno/idempotencer/pkg/logging"
	"go.nanomsg.org/mangos/v3"
	"go.nanomsg.org/mangos/v3/protocol"
	"go.nanomsg.org/mangos/v3/protocol/req"
	_ "go.nanomsg.org/mangos/v3/transport/all"
)

const (
	msgSize      = 1024
	identityHead = "identity"
)

func NewMangosClient(conf ClientConfig, logger logging.Logger) Client {
	if conf.Identity == "" {
		conf.Identity = helpers.UniqIdentity()
	}
	return &mangosClient{
		conf: conf,
		log:  logger,
	}
}

type mangosClient struct {
	sock mangos.Socket
	conf ClientConfig
	log  logging.Logger
}

func (m *mangosClient) Connect() errors.Error {
	m.log.Debug("Connecting", m.conf.ToLoggerCtx())
	defer m.log.Info("Connected", m.conf.ToLoggerCtx())

	if m.conf.Socket == "" {
		return helpers.NewErrWithLog(
			m.log,
			errors.InvalidConfiguration,
			stdErrors.New("destination socket not defined"),
		)
	}

	var err error
	if m.sock, err = req.NewSocket(); err != nil {
		return helpers.NewErrWithLog(m.log, errors.NewRouterSocketCreationFail, err)
	}
	if err = m.sock.SetOption(protocol.OptionBestEffort, true); err != nil {
		return helpers.NewErrWithLog(m.log, errors.UnknownError, err)
	}
	if err = m.sock.Dial(proto + m.conf.Socket); err != nil {
		return helpers.NewErrWithLog(m.log, errors.ConnectSocketFail, err)
	}
	return nil
}

func (m *mangosClient) Disconnect() errors.Error {
	m.log.Debug("Disconnecting", m.conf.ToLoggerCtx())
	defer m.log.Info("Disconnected", m.conf.ToLoggerCtx())

	err := m.sock.Close()
	if err != nil {
		return helpers.NewErrWithLog(m.log, errors.CloseSocketFail, err)
	}
	return nil
}

func (m *mangosClient) Pull() (dto.Msg, errors.Error) {
	rcv, err := m.sock.RecvMsg()
	if err != nil {
		return nil, errors.NewError(errors.PullSocketNotReadyError, err)
	}
	msg := dto.NewByteMsg(rcv.Body)
	rcv.Free()
	logCtx := m.conf.ToLoggerCtx()
	logCtx["data"] = msg.String()
	m.log.Debug("Received", logCtx)
	return msg, nil
}

func (m *mangosClient) Push(msg dto.Msg) errors.Error {
	toSend := mangos.NewMessage(msgSize)
	toSend.Body = append(toSend.Body, bytes.Join(msg.Data(), []byte(""))...)
	head, _ := json.Marshal(map[string]string{
		identityHead: m.conf.Identity,
	})
	toSend.Header = append(toSend.Header, head...)
	if err := m.sock.SendMsg(toSend); err != nil {
		return helpers.NewErrWithLog(m.log, errors.PushSocketError, err)
	}
	toSend.Free()
	logCtx := m.conf.ToLoggerCtx()
	logCtx["data"] = msg.String()
	m.log.Debug("Send", logCtx)
	return nil
}
