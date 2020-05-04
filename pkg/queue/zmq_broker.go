package queue

import (
	"fmt"

	"github.com/c1rno/idempotencer/pkg/errors"
	"github.com/c1rno/idempotencer/pkg/helpers"
	"github.com/c1rno/idempotencer/pkg/logging"
	zmq "github.com/pebbe/zmq4"
)

const (
	proto            = "tcp://"
	READY = "\001" //  Signals worker is ready
)

func NewBroker(conf BrokerConfig, logger logging.Logger) Broker {
	return &simpleLoadBalancer{
		conf:              conf,
		log:               logger,
		downstreamClients: []string{},
	}
}

// Remember, zmq sockets are not thread safe!
type Broker interface {
	// Blocking operation
	Start() errors.Error
	Stop() errors.Error
}

type simpleLoadBalancer struct {
	upstreamSock, downstreamSock *zmq.Socket
	reactor                      *zmq.Reactor
	conf                         BrokerConfig
	log                          logging.Logger
	downstreamClients            []string
}

func (b *simpleLoadBalancer) Stop() errors.Error {
	b.log.Debug("Stopping broker", b.conf.ToLoggerCtx())
	defer b.log.Info("Broker stopped", b.conf.ToLoggerCtx())

	b.reactor.RemoveSocket(b.upstreamSock)
	b.reactor.RemoveSocket(b.downstreamSock)

	err1 := b.upstreamSock.Close()
	err2 := b.downstreamSock.Close()
	if err1 != nil {
		return helpers.NewErrWithLog(b.log, errors.CloseSocketFail, err1)
	}
	if err2 != nil {
		return helpers.NewErrWithLog(b.log, errors.CloseSocketFail, err2)
	}
	return nil
}

func (b *simpleLoadBalancer) Start() errors.Error {
	b.log.Debug("Starting broker", b.conf.ToLoggerCtx())

	if b.conf.InSocket == "" || b.conf.OutSocket == "" {
		return helpers.NewErrWithLog(b.log, errors.InvalidConfiguration, nil)
	}

	var err error
	if b.upstreamSock, err = zmq.NewSocket(zmq.ROUTER); err != nil {
		return helpers.NewErrWithLog(b.log, errors.NewRouterSocketCreationFail, err)
	}
	if b.downstreamSock, err = zmq.NewSocket(zmq.ROUTER); err != nil {
		return helpers.NewErrWithLog(b.log, errors.NewRouterSocketCreationFail, err)
	}
	if err = b.upstreamSock.Bind(proto + b.conf.InSocket); err != nil {
		return helpers.NewErrWithLog(b.log, errors.BindSocketFail, err)
	}
	if err = b.downstreamSock.Bind(proto + b.conf.OutSocket); err != nil {
		return helpers.NewErrWithLog(b.log, errors.BindSocketFail, err)
	}
	b.reactor = zmq.NewReactor()
	b.reactor.AddSocket(
		b.downstreamSock,
		zmq.POLLIN,
		func(e zmq.State) error { return b.handleDownstream() },
	)
	b.log.Info("Broker started", b.conf.ToLoggerCtx())
	if err = b.reactor.Run(-1); err != nil {
		return helpers.NewErrWithLog(b.log, errors.ReactorError, err)
	}
	return nil
}

func (b *simpleLoadBalancer) handleDownstream() error {
	//  Use worker identity for load-balancing
	msg, err := b.downstreamSock.RecvMessage(0)
	if err != nil {
		return err
	}
	b.log.Debug(fmt.Sprintf("Received from downstream: %v", msg), nil)
	identity, msg := unwrap(msg)
	b.downstreamClients = append(b.downstreamClients, identity)
	b.log.Debug(fmt.Sprintf("Add identity: %v", identity), nil)

	//  Enable reader on frontend if we went from 0 to 1 workers
	if len(b.downstreamClients) == 1 {
		b.log.Debug("Start processing upstream", nil)
		b.reactor.AddSocket(
			b.upstreamSock,
			zmq.POLLIN,
			func(e zmq.State) error { return b.handleUpstream() },
		)
	}

	//  Forward message to client if it's not a READY
	if msg[0] != READY {
		if _, err = b.upstreamSock.SendMessage(msg); err != nil {
			helpers.NewErrWithLog(b.log, errors.PushSocketError, err)
		}
		b.log.Debug(fmt.Sprintf("Send to upstream: %v", msg), nil)
	}

	return nil
}

func (b *simpleLoadBalancer) handleUpstream() error {
	//  Get client request, route to first available worker
	msg, err := b.upstreamSock.RecvMessage(0)
	if err != nil {
		return err
	}
	b.log.Debug(fmt.Sprintf("Received from upstream: %v", msg), nil)
	if _, err = b.downstreamSock.SendMessage(b.downstreamClients[0], "", msg); err != nil {
		helpers.NewErrWithLog(b.log, errors.PushSocketError, err)
	}
	b.log.Debug(fmt.Sprintf("Send to downstream: {%v, %v, %v}", b.downstreamClients[0], "", msg), nil)
	b.log.Debug(fmt.Sprintf("Remove identity: %v", b.downstreamClients[0]), nil)
	b.downstreamClients = b.downstreamClients[1:]

	//  Cancel reader on frontend if we went from 1 to 0 workers
	if len(b.downstreamClients) == 0 {
		b.log.Debug("Stop processing upstream", nil)
		b.reactor.RemoveSocket(b.upstreamSock)
	}
	return nil
}

//  Pops frame off front of message and returns it as 'head'
//  If next frame is empty, pops that empty frame.
//  Return remaining frames of message as 'tail'
func unwrap(msg []string) (head string, tail []string) {
	head = msg[0]
	if len(msg) > 1 && msg[1] == "" {
		tail = msg[2:]
	} else {
		tail = msg[1:]
	}
	return
}
