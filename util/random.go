package util

import (
	"fmt"
	"math/rand"
)

var alphabets = "abcdefghijklmnopqrstuvwxyz"
var currencies = []string{USD, EUR, RUPEE}

func RandomName() string {
	name := []byte{}
	for i := 0; i < 6; i++ {
		name = append(name, alphabets[rand.Intn(len(alphabets))])
	}
	return string(name)
}

func RandomCurrency() string {
	return currencies[rand.Intn(len(currencies))]
}

func RandomAmount() int64 {
	return RandomNumber(1000)
}

func RandomNumber(max int) int64 {
	number := rand.Intn(max) + 1
	return int64(number)
}

func RandomEmail() string {
	return fmt.Sprintf("%v@gmail.com", RandomName())
}
