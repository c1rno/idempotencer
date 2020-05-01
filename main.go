package main

import (
	"github.com/c1rno/idempotencer/cmd/debug"
	_ "github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/jackc/pgx/v4"
	"github.com/spf13/cobra"
	_ "github.com/zeromq/goczmq"
)

var root = &cobra.Command{}

func main() {
	root.AddCommand(
		debug.Command,
	)
	if err := root.Execute(); err != nil {
		panic(err)
	}
}
