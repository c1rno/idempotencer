package helpers

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const identityLen = 8

var letters = []rune("abcdefghijklmnopqrstuvwxyz")

func UniqIdentity() string {
	b := make([]rune, identityLen)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
