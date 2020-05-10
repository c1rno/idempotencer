package queue

import (
	stdErrors "errors"

	"github.com/c1rno/idempotencer/pkg/dto"
	"github.com/c1rno/idempotencer/pkg/errors"
	"github.com/c1rno/idempotencer/pkg/helpers"
	"github.com/c1rno/idempotencer/pkg/logging"
	zmq "github.com/pebbe/zmq4"
)

func NewClient(conf ClientConfig, logger logging.Logger) Client {
	return &client{
		conf: conf,
		log:  logger,
	}
}

type Client interface {
	Connect() errors.Error
	Disconnect() errors.Error
	Pull() (dto.Msg, errors.Error)
	Push(dto.Msg) errors.Error
}

type client struct {
	sock *zmq.Socket
	conf ClientConfig
	log  logging.Logger
}

func (c *client) Disconnect() errors.Error {
	c.log.Debug("Disconnecting", c.conf.ToLoggerCtx())
	defer c.log.Info("Disconnected", c.conf.ToLoggerCtx())

	err := c.sock.Close()
	if err != nil {
		return helpers.NewErrWithLog(c.log, errors.CloseSocketFail, err)
	}
	return nil
}

func (c *client) Connect() errors.Error {
	c.log.Debug("Connecting", c.conf.ToLoggerCtx())
	defer c.log.Info("Connected", c.conf.ToLoggerCtx())

	if c.conf.Socket == "" {
		return helpers.NewErrWithLog(
			c.log,
			errors.InvalidConfiguration,
			stdErrors.New("destination socket not defined"),
		)
	}

	var err error
	if c.sock, err = zmq.NewSocket(zmq.REQ); err != nil {
		return helpers.NewErrWithLog(c.log, errors.NewRouterSocketCreationFail, err)
	}
	if err = c.sock.Connect(proto + c.conf.Socket); err != nil {
		return helpers.NewErrWithLog(c.log, errors.ConnectSocketFail, err)
	}
	return nil
}

func (c *client) Push(s dto.Msg) (err errors.Error) {
	defer func() {
		if r := recover(); r != nil {
			if err == nil {
				err = errors.NewError(errors.PushSocketError, nil)
			}
			c.log.Error(err.Error(), map[string]interface{}{
				"panic": r,
			})
		}
	}()
	if _, serr := c.sock.SendMessageDontwait(s.Data()); serr != nil {
		err = helpers.NewErrWithLog(c.log, errors.PushSocketError, serr)
		return
	}
	logCtx := c.conf.ToLoggerCtx()
	logCtx["data"] = s.String()
	c.log.Debug("Send", logCtx)
	return
}

func (c *client) Pull() (d dto.Msg, err errors.Error) {
	defer func() {
		if r := recover(); r != nil {
			if err == nil {
				err = errors.NewError(errors.PullSocketError, nil)
			}
			c.log.Error(err.Error(), map[string]interface{}{
				"panic": r,
			})
		}
	}()
	msg, serr := c.sock.RecvMessage(zmq.DONTWAIT)
	if serr != nil {
		if serr.Error() == unavailable {
			err = errors.NewError(errors.PullSocketNotReadyError, serr)
			return
		}
		err = helpers.NewErrWithLog(c.log, errors.PullSocketError, serr)
		return
	}
	ret := dto.NewStringMsg(msg...)
	logCtx := c.conf.ToLoggerCtx()
	logCtx["data"] = ret.String()
	c.log.Debug("Received", logCtx)
	return ret, err
}
