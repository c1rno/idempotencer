// +build zmq

package cmd

import (
	"github.com/c1rno/idempotencer/cmd/zeromq"
	"github.com/spf13/cobra"
)

func init() {
	ZMQCommand.AddCommand(
		zeromq.BrokerCommand,
		zeromq.DownstreamCommand,
		zeromq.UpstreamCommand,
	)
}

var ZMQCommand = &cobra.Command{
	Use:   "zmq",
	Short: "0MQ commands",
}
