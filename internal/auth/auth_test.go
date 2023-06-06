package auth

import (
	"strings"
	"testing"

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
	sid, key := s.engine.Generate()
	key2, err := s.engine.Validate(sid)
	s.Assert().NoError(err)
	s.Assert().Equal(key, key2)
}

func (s *AuthTestSuite) TestNegativeHexDecrypt() {
	sid := strings.Repeat("g", 16)
	_, err := s.engine.Validate(sid)
	s.Assert().Error(err)
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}
