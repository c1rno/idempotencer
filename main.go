package main

import (
	"github.com/c1rno/idempotencer/cmd/broker"
	"github.com/c1rno/idempotencer/cmd/downstream"
	"github.com/c1rno/idempotencer/cmd/migrate"
	"github.com/c1rno/idempotencer/cmd/upstream"
	"github.com/spf13/cobra"
)

var root = &cobra.Command{}

func main() {
	root.AddCommand(
		broker.Command,
		downstream.Command,
		migrate.Command,
		upstream.Command,
	)
	if err := root.Execute(); err != nil {
		panic(err)
	}
}
