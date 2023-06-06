package storage

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"shortener/internal/cfg"
)

type DBTestSuite struct {
	suite.Suite
}

func (s *DBTestSuite) SetupTest() {
	cfg.New(cfg.FromDefaults())
}

func (s *DBTestSuite) TestGetPostRequest() {
}

func TestDatabase(t *testing.T) {
	suite.Run(t, new(DBTestSuite))
}
