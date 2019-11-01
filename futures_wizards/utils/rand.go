package utils

import (
	"math/rand"
	"time"
)

func RandString(l int) string {
	const (
		src = "ABCDEFGHIJKLMNOPQLSTUVWXYZabcdefghijklmnopqlstuvwxyz0123456789"
	)

	rand.Seed(time.Now().UnixNano())

	b := make([]byte, l)
	for i := 0; i < l; i++ {
		b[i] = src[rand.Intn(len(src))]
	}

	return string(b)
}
