package utils

import (
	"math/rand"
)

var letter = []rune("abcdefghijklmnopqrstuvwxyz")

// RandomString generates a random string of n characters
func RandomString(n int) string {

	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

func RandomInt(min, max int) int {
	return rand.Intn(max-min) + min
}

func RandomName() string {
	return RandomString(6)
}

func RandomBalance() int64 {
	return int64(RandomInt(0, 1000))
}

func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "ILS"}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}
