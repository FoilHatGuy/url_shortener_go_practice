package storage

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"shortener/internal/cfg"
)

type DBTestSuite struct {
	suite.Suite
	config *cfg.ConfigT
	db     DatabaseORM
}

func (s *DBTestSuite) SetupTest() {
	s.config = cfg.New(cfg.FromDefaults())
	s.db = databaseInitialize(s.config)
}

func (s *DBTestSuite) TestGetPostRequest() {
}

func TestDatabase(t *testing.T) {
	suite.Run(t, new(DBTestSuite))
}
