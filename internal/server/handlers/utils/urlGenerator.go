package utils

import (
	"crypto/rand"
	"math/big"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var runeLength = big.NewInt(int64(len(letters)))

// RandSeq
// Generates random sequence of letters (both upper- and lower-case) of desired length
func RandSeq(length int) string {
	b := make([]rune, length)
	for i := range b {
		num, _ := rand.Int(rand.Reader, runeLength)
		b[i] = letters[num.Int64()]
	}
	return string(b)
}
