package util

import (
	"math/rand"
	"strings"
)

const alphabets = "abcdefghijklmnopqrstuvwxyz"

// RandInt generates a random integer between min and max
func RandInt(min , max int64) int64 {
	return min + rand.Int63n(max - min + 1)
}


// RandString generates a random string 
func RandString(n int) string {

	c := len(alphabets)
	var sb strings.Builder
	
	for i := 0 ; i < n ; i ++ {
		owner := alphabets[rand.Intn(c)]
		sb.WriteByte(owner)
	}
	
	return sb.String()
}

// RandOwner generates a random owner name
func RandOwner() string {
	return RandString(6)
}

// RandMoney generates a random money amount
func RandMoney() int64 {
	return RandInt(0,333280)
}

// RandAccountId generates a random Account ID
func RandAccountId() int64 {
	return RandInt(0, 9999999999)
}

// RandAmount generates a random Amount
func RandAmount() int64 {
	return RandInt(0, 99999999)
}


//  RandCurrency generates a random currency
func RandCurrency() string {
	currency := []string{"EUR","USD","YEN"}
	k := len(currency)
	return currency[rand.Intn(k)]
}