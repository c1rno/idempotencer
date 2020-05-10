package rawtcp

import (
	"bufio"
	"fmt"
	"net"

	"github.com/spf13/cobra"
)

const (
	testTCPPort = "8888"
	testTCPHost = "0.0.0.0"
	buffSize    = 1024
)

var UpstreamTestCmd = &cobra.Command{
	Use:   `upstream`,
	Short: `Raw upstream tcp socket to compare with 0MQ`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Bind at %s:%s\n", testTCPHost, testTCPPort)
		l, err := net.Listen("tcp", testTCPHost+":"+testTCPPort)
		if err != nil {
			fmt.Printf("Error listening: %v\n", err)
			panic(err)
		}
		defer l.Close()

		i := 0
		for {
			conn, err := l.Accept()
			if err != nil {
				fmt.Printf("%d Error accepting: %v\n", i, err)
				panic(err)
			}
			go func(i int, conn net.Conn) {
				for {
					i += 1
					msg, err := bufio.NewReader(conn).ReadString('\n')
					if err != nil {
						fmt.Printf("%d, Error reading: %v\n", i, err)
						return
					}
					fmt.Printf("%d Received: %s\n", i, msg)
					if _, err = conn.Write([]byte(fmt.Sprintf("%d Message received\n", i))); err != nil {
						fmt.Printf("%d, Error writing: %v\n", i, err)
					}
				}
				if err = conn.Close(); err != nil {
					fmt.Printf("%d, Error conn closing: %v\n", i, err)
				}
			}(i, conn)
		}
	},
}
