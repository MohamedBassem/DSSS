package server

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func generateRandomString(l int) string {

	seed := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	ret := ""
	for ; l > 0; l-- {
		ret += string(seed[rand.Intn(len(seed))])
	}

	return ret
}
