package storage

import (
	"context"
	"crypto/rand"
	"github.com/stretchr/testify/suite"
	"math/big"
	"shortener/internal/cfg"
	"testing"
)

type MemoryTestSuite struct {
	suite.Suite
	config *cfg.ConfigT
	ctrl   DatabaseORM
	ctx    context.Context
}

func (s *MemoryTestSuite) SetupTest() {
	s.config = cfg.Initialize()
	s.config.Storage.StorageType = "file"
	s.config.Storage.SavePath = "../data"
	s.config.Storage.AutosaveInterval = 5
	Initialize(s.config)
	s.ctrl = Controller
	s.ctx = context.Background()
}

func (s *MemoryTestSuite) TestPing() {
	ping := s.ctrl.Ping(s.ctx)
	s.Assert().True(ping)
}

func (s *MemoryTestSuite) TestNegative() {
	uid := generateString(20)
	array, err := s.ctrl.GetURLByOwner(s.ctx, uid)
	s.Assert().NoError(err)
	s.Assert().Nil(array)

	uid2 := generateString(20)
	err = s.ctrl.Delete(s.ctx, []string{}, uid2)
	s.Assert().Error(err)

	shortURL := generateString(10)
	result, ok, err := s.ctrl.GetURL(s.ctx, shortURL)
	s.Assert().Error(err)
	s.Assert().False(ok)
	s.Assert().Equal(result, "")
}

func (s *MemoryTestSuite) TestAddGetURL() {
	uid := generateString(20)

	originalURL := generateString(20)
	shortURL := generateString(10)
	ok, _, err := s.ctrl.AddURL(s.ctx, originalURL, shortURL, uid)
	s.Assert().NoError(err)
	s.Assert().True(ok)

	originalURL2 := generateString(20)
	shortURL2 := generateString(10)
	ok, _, err = s.ctrl.AddURL(s.ctx, originalURL2, shortURL2, uid)
	s.Assert().NoError(err)
	s.Assert().True(ok)

	original, ok, err := s.ctrl.GetURL(s.ctx, shortURL)
	s.Assert().NoError(err)
	s.Assert().True(ok)
	s.Assert().Equal(originalURL, original)

	expectedArray := []URLOfOwner{
		{
			s.config.Server.BaseURL + "/" + shortURL,
			originalURL,
		},
		{
			s.config.Server.BaseURL + "/" + shortURL2,
			originalURL2,
		},
	}
	array, err := s.ctrl.GetURLByOwner(s.ctx, uid)
	s.Assert().NoError(err)
	s.Assert().Equal(expectedArray, array)
}

func (s *MemoryTestSuite) TestDeletion() {
	uid := generateString(20)

	originalURL := generateString(20)
	shortURL := generateString(10)
	ok, _, err := s.ctrl.AddURL(s.ctx, originalURL, shortURL, uid)
	s.Assert().NoError(err)
	s.Assert().True(ok)

	result, ok, err := s.ctrl.GetURL(s.ctx, shortURL)
	s.Assert().NoError(err)
	s.Assert().True(ok)
	s.Assert().Equal(result, originalURL)

	err = s.ctrl.Delete(s.ctx, []string{shortURL}, uid)
	s.Assert().NoError(err)

	result, ok, err = s.ctrl.GetURL(s.ctx, shortURL)
	s.Assert().NoError(err)
	s.Assert().True(ok)
	s.Assert().Equal(result, "")
}

func TestMemory(t *testing.T) {
	suite.Run(t, new(MemoryTestSuite))
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var length = big.NewInt(int64(len(letters)))

func generateString(n int) string {
	b := make([]rune, n)
	for i := range b {
		num, _ := rand.Int(rand.Reader, length)
		b[i] = letters[num.Int64()]
	}
	return string(b)
}
