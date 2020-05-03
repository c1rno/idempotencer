package debug

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/c1rno/idempotencer/pkg/config"
	"github.com/c1rno/idempotencer/pkg/errors"
	"github.com/c1rno/idempotencer/pkg/logging"
	"github.com/c1rno/idempotencer/pkg/metrics"
	"github.com/davecgh/go-spew/spew"
	zmq "github.com/pebbe/zmq4"
	"github.com/spf13/cobra"
)

const (
	NBR_CLIENTS  = 10
	NBR_WORKERS  = 3
	WORKER_READY = "\001" //  Signals worker is ready
)

// This http://zguide.zeromq.org/page:all#The-CZMQ-High-Level-API
// https://github.com/pebbe/zmq4/blob/44644726cbe40e63c27ee9c4d822449e7754971e/examples/lbbroker3.go#L93
var Command = &cobra.Command{
	Use:   `debug`,
	Short: `Don't use it`,
	Run: func(cmd *cobra.Command, args []string) {
		spew.Config.DisablePointerAddresses = true
		spew.Config.DisableCapacities = true
		spew.Config.SortKeys = true
		spew.Config.SpewKeys = true

		logger := logging.NewLogger(0)
		conf, err := config.NewConfig(logger)
		logger.Debug(fmt.Sprintf("Env config: %s", spew.Sdump(conf)), map[string]interface{}{
			"err": err,
		})
		logger = logging.NewLogger(conf.LogLevel)
		shutdown := metrics.RunMetricsSrv(conf.Metrics, func() errors.Error { return nil }, logger)
		defer shutdown(context.Background())

		lbbroker := &lbbroker_t{}
		lbbroker.frontend, _ = zmq.NewSocket(zmq.ROUTER)
		lbbroker.backend, _ = zmq.NewSocket(zmq.ROUTER)
		defer lbbroker.frontend.Close()
		defer lbbroker.backend.Close()
		lbbroker.frontend.Bind("ipc://frontend.ipc")
		lbbroker.backend.Bind("ipc://backend.ipc")

		for client_nbr := 0; client_nbr < NBR_CLIENTS; client_nbr++ {
			go client_task()
		}
		for worker_nbr := 0; worker_nbr < NBR_WORKERS; worker_nbr++ {
			go worker_task()
		}

		//  Queue of available workers
		lbbroker.workers = make([]string, 0, 10)

		//  Prepare reactor and fire it up
		lbbroker.reactor = zmq.NewReactor()
		lbbroker.reactor.AddSocket(lbbroker.backend, zmq.POLLIN,
			func(e zmq.State) error { return handle_backend(lbbroker) })
		lbbroker.reactor.Run(-1)
	},
}

//  Basic request-reply client using REQ socket
//
func client_task() {
	client, _ := zmq.NewSocket(zmq.REQ)
	defer client.Close()
	client.Connect("ipc://frontend.ipc")

	//  Send request, get reply
	for {
		client.SendMessage("HELLO")
		reply, _ := client.RecvMessage(0)
		if len(reply) == 0 {
			break
		}
		fmt.Println("Client:", strings.Join(reply, "\n\t"))
		time.Sleep(time.Second)
	}
}

//  Worker using REQ socket to do load-balancing
//
func worker_task() {
	worker, _ := zmq.NewSocket(zmq.REQ)
	defer worker.Close()
	worker.Connect("ipc://backend.ipc")

	//  Tell broker we're ready for work
	worker.SendMessage(WORKER_READY)

	//  Process messages as they arrive
	for {
		msg, e := worker.RecvMessage(0)
		if e != nil {
			break //  Interrupted
		}
		msg[len(msg)-1] = "OK"
		worker.SendMessage(msg)
	}
}

//  Our load-balancer structure, passed to reactor handlers
type lbbroker_t struct {
	frontend *zmq.Socket //  Listen to clients
	backend  *zmq.Socket //  Listen to workers
	workers  []string    //  List of ready workers
	reactor  *zmq.Reactor
}

//  In the reactor design, each time a message arrives on a socket, the
//  reactor passes it to a handler function. We have two handlers; one
//  for the frontend, one for the backend:

//  Handle input from client, on frontend
func handle_frontend(lbbroker *lbbroker_t) error {

	//  Get client request, route to first available worker
	msg, err := lbbroker.frontend.RecvMessage(0)
	if err != nil {
		return err
	}
	lbbroker.backend.SendMessage(lbbroker.workers[0], "", msg)
	lbbroker.workers = lbbroker.workers[1:]

	//  Cancel reader on frontend if we went from 1 to 0 workers
	if len(lbbroker.workers) == 0 {
		lbbroker.reactor.RemoveSocket(lbbroker.frontend)
	}
	return nil
}

//  Handle input from worker, on backend
func handle_backend(lbbroker *lbbroker_t) error {
	//  Use worker identity for load-balancing
	msg, err := lbbroker.backend.RecvMessage(0)
	if err != nil {
		return err
	}
	identity, msg := unwrap(msg)
	lbbroker.workers = append(lbbroker.workers, identity)

	//  Enable reader on frontend if we went from 0 to 1 workers
	if len(lbbroker.workers) == 1 {
		lbbroker.reactor.AddSocket(lbbroker.frontend, zmq.POLLIN,
			func(e zmq.State) error { return handle_frontend(lbbroker) })
	}

	//  Forward message to client if it's not a READY
	if msg[0] != WORKER_READY {
		lbbroker.frontend.SendMessage(msg)
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
