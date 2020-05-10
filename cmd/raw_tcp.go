package cmd

import (
	"github.com/c1rno/idempotencer/cmd/rawtcp"
	"github.com/spf13/cobra"
)

func init() {
	RawTCPCommand.AddCommand(
		rawtcp.DownstreamTestCmd,
		rawtcp.UpstreamTestCmd,
	)
}

var RawTCPCommand = &cobra.Command{
	Use:   "tcp",
	Short: "Raw tcp commands",
}
