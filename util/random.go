package util

import (
	"math/rand"
	"strings"
	"time"
)

// generate a random uint64

const alp = "abcdefghijklmnopqrstuvwxyz"

func init()  {
	rand.Seed(time.Now().UnixNano())
}

// generate random integer btw min and max

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max - min + 1)	
}

// generates random string of length n
func RandomStr(n int) string {
	var sb strings.Builder
	
	k := len(alp)

	for i := 0; i < n; i++ {
		c := alp[rand.Intn(k)]	
		sb.WriteByte(c)
	}

	return sb.String()
}

// generates a random owner name
func RandomOwner()string  {
	return RandomStr(6)
}

// generates a random amount of money
func RandomMoney() int64 {
	return RandomInt(20, 1000)
}

// returns a random currency from the slice decalred
func RandomCurrency() string  {
	currencies := []string{EUR, NGN, USD}

	n := len(currencies)

	return currencies[rand.Intn(n)]
}