package urlgenerator

import (
	"crypto/rand"
	"math/big"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var length = big.NewInt(int64(len(letters)))

func RandSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		num, _ := rand.Int(rand.Reader, length)
		b[i] = letters[num.Int64()]
	}
	return string(b)
}
