package util

import (
	"math/rand"
	"strings"
	"time"
)

const (
	alphabet = "abcdefghijklmnopqrstuvwxyz"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

//RandomInt generates a random integer between min and max
func RandomInt(min int64, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

//RandomString generates a random string with length n
func RandomString(n int) string {
	var sb strings.Builder
	for i := 0; i < n; i++ {
		char := alphabet[rand.Intn(len(alphabet))]
		sb.WriteByte(char)
	}
	return sb.String()
}

//RandomOwner generates a random owner name
func RandomOwner() string {
	return RandomString(6)
}

//RandomMoney generates a random amount of money
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

//RandomCurrency generates a random currency code
func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "CAD", "BRL"}
	return currencies[rand.Intn(len(currencies))]
}
