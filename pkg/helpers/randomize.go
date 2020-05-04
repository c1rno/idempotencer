package helpers

import (
	"fmt"
	"math/rand"
	"time"
)

var source = rand.NewSource(time.Now().Unix())

func RandInt() int64 {
	return source.Int63()
}

func UniqIdentity() string {
	return fmt.Sprintf("%d", RandInt())
}
