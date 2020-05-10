// +build !zmq

package cmd

import (
	"github.com/spf13/cobra"
)

var ZMQCommand = &cobra.Command{
	Use:   "zmq",
	Short: "0MQ commands",
	Run: func(*cobra.Command, []string) {
		panic(`Not supported: build with '-tags="zmq"'`)
	},
}
