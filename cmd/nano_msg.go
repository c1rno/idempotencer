package cmd

import (
	"github.com/c1rno/idempotencer/cmd/nanomsg"
	"github.com/spf13/cobra"
)

func init() {
	NanomsgCommand.AddCommand(
		nanomsg.UpstreamCmd,
		nanomsg.DownstreamCmd,
		nanomsg.BrokerCmd,
	)
}

var NanomsgCommand = &cobra.Command{
	Use:   "nanomsg",
	Short: "Nanomsg commands",
}
