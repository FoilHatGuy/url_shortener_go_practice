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

var AuthEngine sessionValidator

func Init(config *cfg.ConfigT) {
	key := sha256.Sum256([]byte(config.Shortener.Secret))
	aesBlock, err := aes.NewCipher(key[:])
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	AuthEngine = &engineT{
		crypt:  aesBlock,
		secret: key[:],
		config: config,
	}
}

type engineT struct {
	crypt  cipher.Block
	secret []byte
	config *cfg.ConfigT
}

type sessionValidator interface {
	Validate(cookie string) (key string, err error)
	Generate() (cookie string, key string)
}

func (e *engineT) Validate(cookie string) (key string, err error) {
	src, err := hex.DecodeString(cookie)
	if err != nil {
		return "", err
	}
	dst := make([]byte, aes.BlockSize)
	e.crypt.Decrypt(dst, src)
	res := hex.EncodeToString(dst)
	return res, nil
}

func (e *engineT) Generate() (cookie string, key string) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", ""
	}
	dst := make([]byte, aes.BlockSize)
	e.crypt.Encrypt(dst, b)

	return hex.EncodeToString(dst), hex.EncodeToString(b)
}
