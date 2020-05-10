package queue

import (
	"context"
	"sync"

	"github.com/c1rno/idempotencer/pkg/errors"
	"github.com/c1rno/idempotencer/pkg/helpers"
	"github.com/c1rno/idempotencer/pkg/logging"
	"go.nanomsg.org/mangos/v3"
	"go.nanomsg.org/mangos/v3/protocol"
	"go.nanomsg.org/mangos/v3/protocol/rep"
	_ "go.nanomsg.org/mangos/v3/transport/all"
)

const (
	defaultWorkersNum = 10
)

func NewMangosBroker(conf BrokerConfig, logger logging.Logger) Broker {
	ctx, done := context.WithCancel(context.Background())
	return &mangosLoadBalancer{
		conf:      conf,
		log:       logger,
		wg:        sync.WaitGroup{},
		ctx:       ctx,
		done:      done,
		upConns:   map[uint32]*conn{},
		downConns: map[uint32]*conn{},
	}
}

type mangosLoadBalancer struct {
	upstreamSock, downstreamSock mangos.Socket
	upConns, downConns           map[uint32]*conn
	conf                         BrokerConfig
	log                          logging.Logger
	wg                           sync.WaitGroup
	ctx                          context.Context
	done                         context.CancelFunc
}

func (m *mangosLoadBalancer) Start() errors.Error {
	m.log.Debug("Starting broker", m.conf.ToLoggerCtx())

	if m.conf.InSocket == "" || m.conf.OutSocket == "" {
		return helpers.NewErrWithLog(m.log, errors.InvalidConfiguration, nil)
	}

	var err error
	if m.upstreamSock, err = rep.NewSocket(); err != nil {
		return helpers.NewErrWithLog(m.log, errors.NewRouterSocketCreationFail, err)
	}
	if err = m.upstreamSock.SetOption(protocol.OptionBestEffort, true); err != nil {
		return helpers.NewErrWithLog(m.log, errors.UnknownError, err)
	}
	if m.downstreamSock, err = rep.NewSocket(); err != nil {
		return helpers.NewErrWithLog(m.log, errors.NewRouterSocketCreationFail, err)
	}
	if err = m.downstreamSock.SetOption(protocol.OptionBestEffort, true); err != nil {
		return helpers.NewErrWithLog(m.log, errors.UnknownError, err)
	}

	if err = m.downstreamSock.Listen(proto + m.conf.OutSocket); err != nil {
		return helpers.NewErrWithLog(m.log, errors.BindSocketFail, err)
	}
	if err = m.upstreamSock.Listen(proto + m.conf.InSocket); err != nil {
		return helpers.NewErrWithLog(m.log, errors.BindSocketFail, err)
	}

	if err = m.startRouting(); err != nil {
		return helpers.NewErrWithLog(m.log, errors.LoadBalancerError, err)
	}
	m.log.Info("Broker started", m.conf.ToLoggerCtx())
	<-m.ctx.Done()
	return nil
}

func (m *mangosLoadBalancer) Stop() errors.Error {
	m.log.Debug("Stopping broker", m.conf.ToLoggerCtx())
	defer m.log.Info("Broker stopped", m.conf.ToLoggerCtx())

	m.done()
	m.wg.Wait()
	err1 := m.upstreamSock.Close()
	err2 := m.downstreamSock.Close()
	if err1 != nil {
		return helpers.NewErrWithLog(m.log, errors.CloseSocketFail, err1)
	}
	if err2 != nil {
		return helpers.NewErrWithLog(m.log, errors.CloseSocketFail, err2)
	}
	return nil
}

func (m *mangosLoadBalancer) startRouting() errors.Error {
	m.downstreamSock.SetPipeEventHook(func(e mangos.PipeEvent, p mangos.Pipe) {
		m.log.Debug("Downstream hook fired", map[string]interface{}{
			"event":   e,
			"id":      p.ID(),
			"address": p.Address(),
		})
		if e == mangos.PipeEventAttached {
			c, _ := newConn(m.ctx, m.downstreamSock)
			m.downConns[p.ID()] = c
			_, _ = c.conn.Recv()
			// go downWorker(p.ID(), m)
		}
		if e == mangos.PipeEventDetached {
			if c, ok := m.downConns[p.ID()]; ok {
				_ = c.conn.Close()
				delete(m.downConns, p.ID())
			}
		}
	})
	m.upstreamSock.SetPipeEventHook(func(e mangos.PipeEvent, p mangos.Pipe) {
		m.log.Debug("Upstream hook fired", map[string]interface{}{
			"event":   e,
			"id":      p.ID(),
			"address": p.Address(),
		})
		if e == mangos.PipeEventAttached {
			m.upConns[p.ID()], _ = newConn(m.ctx, m.upstreamSock)
			go upWorker(p.ID(), m)
		}
		if e == mangos.PipeEventDetached {
			if c, ok := m.upConns[p.ID()]; ok {
				_ = c.conn.Close()
				delete(m.upConns, p.ID())
			}
		}
	})
	return nil
}

type conn struct {
	targetID uint32
	conn     mangos.Context
	ctx      context.Context
	done     context.CancelFunc
}

func (c *conn) close() error {
	c.done()
	return c.conn.Close()
}

func newConn(ctx context.Context, sock mangos.Socket) (*conn, error) {
	var (
		err     error
		connCtx mangos.Context
	)
	if connCtx, err = sock.OpenContext(); err != nil {
		return nil, err
	}
	c := &conn{conn: connCtx}
	c.ctx, c.done = context.WithCancel(ctx)
	return c, nil
}

func upWorker(id uint32, m *mangosLoadBalancer) {
	upConn, ok := m.upConns[id]
	if !ok {
		return
	}
	var (
		err        error
		downConn   *conn = nil
		downTarget uint32
		msg        *mangos.Message = nil
	)
	defer func() {
		if downConn != nil {
			downConn.targetID = 0
		}
	}()
LOOP:
	if upConn.ctx.Err() != nil {
		return
	}
	for _downTarget, _downConn := range m.downConns {
		if _downConn.targetID == 0 {
			downTarget = _downTarget
			downConn = _downConn
			break
		}
	}
	if downConn == nil || downTarget == 0 {
		goto LOOP
	}
	upConn.targetID = downTarget
	downConn.targetID = id
	m.log.Info("Paired", map[string]interface{}{
		"upstreamID":   id,
		"downstreamID": downTarget,
	})

	for upConn.ctx.Err() == nil {
		if msg == nil {
			if msg, err = upConn.conn.RecvMsg(); err != nil {
				m.log.Error("Fail to read upstream message", map[string]interface{}{
					helpers.ErrField: errors.NewError(errors.UnknownError, err),
				})
				continue
			}
		}
		m.log.Debug("Forwarding from up to down", map[string]interface{}{
			"msg": string(msg.Body),
		})
		if err = downConn.conn.SendMsg(msg); err != nil {
			m.log.Error("Fail to write downstream message", map[string]interface{}{
				helpers.ErrField: errors.NewError(errors.UnknownError, err),
			})
			_ = upConn.conn.Send([]byte("broker: " + err.Error()))
			downConn = nil
			goto LOOP
		}
		if _, err = downConn.conn.Recv(); err != nil {
			m.log.Error("Fail to read downstream message", map[string]interface{}{
				helpers.ErrField: errors.NewError(errors.UnknownError, err),
			})
			_ = upConn.conn.Send([]byte("broker: " + err.Error()))
		}
		_ = upConn.conn.Send([]byte("OK"))

		msg.Free()
		m.log.Debug("Successful forwarded from up to down", map[string]interface{}{
			"msg": string(msg.Body),
		})
		msg = nil
	}
}
