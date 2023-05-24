package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"shortener/internal/cfg"
)

var AuthEngine SIDValidator

func init() {
	key := sha256.Sum256([]byte(cfg.Shortener.Secret))
	aesBlock, err := aes.NewCipher(key[:])
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	AuthEngine = &engineT{
		crypt:  aesBlock,
		secret: key[:],
	}
}

type engineT struct {
	crypt  cipher.Block
	secret []byte
}

type SIDValidator interface {
	Validate(string) (string, error)
	Generate() (string, string)
}

func (e *engineT) Validate(s string) (string, error) {
	src, _ := hex.DecodeString(s)
	dst := make([]byte, aes.BlockSize)
	e.crypt.Decrypt(dst, src)
	res := hex.EncodeToString(dst)
	return res, nil
}

func (e *engineT) Generate() (string, string) {

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", ""
	}
	dst := make([]byte, aes.BlockSize)
	e.crypt.Encrypt(dst, b)

	return hex.EncodeToString(dst), hex.EncodeToString(b)
}
