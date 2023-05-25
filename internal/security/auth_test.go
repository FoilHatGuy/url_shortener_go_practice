package security

import (
	"github.com/stretchr/testify/suite"
	"shortener/internal/cfg"
	"strings"
	"testing"
)

type AuthTestSuite struct {
	suite.Suite
	engine engineT
	config *cfg.ConfigT
}

func (s *AuthTestSuite) SetupSuite() {
	config := cfg.Initialize()
	s.config = config
	Init(config)
}

func (s *AuthTestSuite) TestCreateAndValidate() {
	sid, key := AuthEngine.Generate()
	key2, err := AuthEngine.Validate(sid)
	s.Assert().NoError(err)
	s.Assert().Equal(key, key2)
}

func (s *AuthTestSuite) TestNegativeHexDecrypt() {
	sid := strings.Repeat("g", 16)
	_, err := AuthEngine.Validate(sid)
	s.Assert().Error(err)
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}
