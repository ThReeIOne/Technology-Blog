package lib

import (
	"math/rand"
	"time"
)

var random *rand.Rand

func init() {
	random = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func randBySource(source []rune, length int) string {
	target := make([]rune, length)
	for i := range target {
		target[i] = source[random.Intn(len(source))]
	}
	return string(target)
}

func RandNumbers(length int) string {
	source := []rune("0123456789")
	return randBySource(source, length)
}

func RandNumbersLetters(length int) string {
	source := []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	return randBySource(source, length)
}

func RandIndex(length int) int {
	return random.Intn(length)
}
