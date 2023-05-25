package storage

import (
	"github.com/stretchr/testify/suite"
	"shortener/internal/cfg"
	"testing"
)

type DBTestSuite struct {
	suite.Suite
}

func (s *DBTestSuite) SetupTest() {
	cfg.Initialize()

}

func (s *DBTestSuite) TestGetPostRequest() {

}

func TestDatabase(t *testing.T) {
	suite.Run(t, new(DBTestSuite))
}
