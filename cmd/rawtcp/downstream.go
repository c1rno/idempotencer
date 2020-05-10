package rawtcp

import (
	"bufio"
	"fmt"
	"net"

	"github.com/spf13/cobra"
)

const (
	testTCPTarget = "localhost"
)

var DownstreamTestCmd = &cobra.Command{
	Use:   `downstream`,
	Short: `Raw downstream tcp socket to compare with 0MQ`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Connect to %s:%s\n", testTCPTarget, testTCPPort)
		conn, err := net.Dial("tcp", testTCPTarget+":"+testTCPPort)
		if err != nil {
			fmt.Printf("Error connecting: %v\n", err)
			panic(err)
		}
		defer conn.Close()

		i := 0
		for {
			i += 1
			msg := fmt.Sprintf("msg-%d", i)
			fmt.Printf("Text to send: %s\n", msg)
			fmt.Fprintf(conn, "%s%s", msg, "\n")
			recv, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				fmt.Printf("Receive err: %v", err)
			} else {
				fmt.Printf("Message from server: %s", recv)
			}
		}
	},
}
