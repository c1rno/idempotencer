package main

import (
	"github.com/c1rno/idempotencer/cmd"

	// _ "github.com/jackc/pgx/v4"
	"github.com/spf13/cobra"
)

var root = &cobra.Command{}

func main() {
	root.AddCommand(
		cmd.MigrateCommand,
		cmd.NanomsgCommand,
		cmd.RawTCPCommand,
		cmd.ZMQCommand,
	)
	if err := root.Execute(); err != nil {
		panic(err)
	}
}
