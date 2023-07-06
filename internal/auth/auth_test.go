//go:build unit
// +build unit

package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/suite"

	"shortener/internal/cfg"
)

type AuthTestSuite struct {
	suite.Suite
	engine *EngineT
	config *cfg.ConfigT
}

func (s *AuthTestSuite) SetupSuite() {
	s.config = cfg.New(cfg.FromDefaults())
	s.engine = New(s.config)
}

func (s *AuthTestSuite) TestCreateAndValidate() {
	sid, key, err := s.engine.Generate()
	s.Assert().NoError(err)
	key2, err := s.engine.Validate(sid)
	s.Assert().NoError(err)
	s.Assert().Equal(key, key2)
}

func (s *AuthTestSuite) TestNegativeHexDecrypt() {
	sid := strings.Repeat("g", 16)
	_, err := s.engine.Validate(sid)
	s.Assert().Error(err)
}

func (s *AuthTestSuite) TestErrorWhileReading() {
	faultyEngine := New(s.config)

	errText := "reader returned error"
	faultyEngine.randomReader = iotest.ErrReader(errors.New(errText))
	_, _, err := faultyEngine.Generate()
	s.Assert().ErrorContains(err, errText)

	faultyEngine.randomReader = iotest.ErrReader(errors.New(errText))

	_, _, err = faultyEngine.GetCertificate()
	s.Assert().ErrorContains(err, "reading random bytes caused a panic")
}

func (s *AuthTestSuite) TestCertGenerator() {
	_, certKey, err := s.engine.GetCertificate()
	s.Assert().NoError(err)
	message := "the test message for checking certificate + key pair"

	keyBlock, _ := pem.Decode([]byte(certKey))
	privateKey, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	s.Assert().NoError(err)

	encMessage, err := rsa.EncryptPKCS1v15(s.engine.randomReader, &privateKey.PublicKey, []byte(message))
	s.Assert().NoError(err)

	decMessage, err := rsa.DecryptPKCS1v15(s.engine.randomReader, privateKey, encMessage)
	s.Assert().NoError(err)

	s.Assert().Equal(message, string(decMessage))
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}
