package main

import (
	"fmt"

	_ "github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/zeromq/goczmq"
)

func main() {
	fmt.Println("Working!")
}
