//go:build unit
// +build unit

package storage

import (
	"context"
	"errors"
	"net/url"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"

	"shortener/internal/cfg"

	pgxV4 "github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v5"

	"github.com/golang/mock/gomock"

	"github.com/chrisyxlee/pgxpoolmock"
	"github.com/stretchr/testify/suite"
)

type newRows struct {
	pgxV4.Rows
}

func (r newRows) FieldDescriptions() []pgconn.FieldDescription {
	return nil
}

func (r newRows) CommandTag() pgconn.CommandTag {
	return pgconn.CommandTag{}
}

func (r newRows) Conn() *pgx.Conn {
	return nil
}

type poolMockWrapper struct {
	pool *pgxpoolmock.MockPgxIface
}

func (p poolMockWrapper) Close() {
	p.pool.Close()
}

func (p poolMockWrapper) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	exec, err := p.pool.Exec(ctx, sql, arguments...)
	commTag := pgconn.NewCommandTag(exec.String())
	return commTag, err //nolint
}

func (p poolMockWrapper) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	rows, err := p.pool.Query(ctx, sql, args...)
	rows2 := newRows{rows}
	return rows2, err //nolint
}

func (p poolMockWrapper) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return p.pool.QueryRow(ctx, sql, args...)
}

func (p poolMockWrapper) Ping(_ context.Context) error {
	return nil
}

type DBTestSuite struct {
	suite.Suite
	ctx    context.Context
	config *cfg.ConfigT
	db     DatabaseORM
	expect *pgxpoolmock.MockPgxIfaceMockRecorder
}

func (s *DBTestSuite) SetupTest() {
	s.config = cfg.New(cfg.FromDefaults())
	pgxPool := &poolMockWrapper{pool: pgxpoolmock.NewMockPgxIface(gomock.NewController(s.T()))}
	s.ctx = context.Background()
	s.expect = pgxPool.pool.EXPECT()
	s.db = &databaseT{
		database: pgxPool,
		config:   s.config,
	}
}

func (s *DBTestSuite) TestPing() {
	ping := s.db.Ping(s.ctx)
	s.Assert().True(ping)
}

func (s *DBTestSuite) TestInit() {
	s.expect.Exec(gomock.Any(), gomock.Any())
	s.expect.Exec(gomock.Any(), gomock.Any())
	s.db.Initialize()

	s.expect.Exec(gomock.Any(), gomock.Any())
	s.expect.Exec(gomock.Any(), gomock.Any()).Return(nil, pgx.ErrTxClosed)
	s.db.Initialize()

	s.expect.Exec(gomock.Any(), gomock.Any()).Return(nil, pgx.ErrTxClosed)
	s.db.Initialize()
}

func (s *DBTestSuite) TestCreation() {
	newC := cfg.New(cfg.FromDefaults(), cfg.WithStorage(cfg.StorageT{DatabaseDSN: "dbname=test"}))
	New(newC)
	// To-Do finish test case
}

func (s *DBTestSuite) TestNegative() {
	uid := generateString(20)
	s.expect.Query(gomock.Any(), gomock.Any(), uid).Return(
		pgxpoolmock.NewRows([]string{"short_url", "original_url"}).ToPgxRows(),
		pgx.ErrNoRows)
	array, err := s.db.GetURLByOwner(s.ctx, uid)
	s.Assert().Error(err)
	s.Assert().Nil(array)

	uid2 := generateString(20)
	s.expect.Exec(gomock.Any(), gomock.Any(), uid2).Return(nil, nil)
	err = s.db.Delete(s.ctx, []string{}, uid2)
	s.Assert().NoError(err)

	testErr := errors.New("testerr")
	s.expect.Exec(gomock.Any(), gomock.Any(), uid2).Return(nil, testErr)
	err = s.db.Delete(s.ctx, []string{}, uid2)
	s.Assert().ErrorIs(err, testErr)

	shortURL := generateString(10)
	s.expect.QueryRow(gomock.Any(), gomock.Any(), shortURL).Return(pgxpoolmock.NewRow("", false).WithError(pgx.ErrNoRows))
	result, ok, err := s.db.GetURL(s.ctx, shortURL)
	s.Assert().Error(err)
	s.Assert().False(ok)
	s.Assert().Equal(result, "")
}

func (s *DBTestSuite) TestAddGetURL() {
	uid := generateString(20)

	originalURL := generateString(20)
	shortURL := generateString(10)

	s.expect.Exec(gomock.Any(), gomock.Any(), shortURL, originalURL)
	s.expect.QueryRow(gomock.Any(), gomock.Any(), originalURL).Return(pgxpoolmock.NewRow(shortURL, originalURL))
	s.expect.Exec(gomock.Any(), gomock.Any(), uid, shortURL)
	ok, _, err := s.db.AddURL(s.ctx, originalURL, shortURL, uid)
	s.Assert().NoError(err)
	s.Assert().True(ok)

	originalURL2 := generateString(20)
	shortURL2 := generateString(10)
	s.expect.Exec(gomock.Any(), gomock.Any(), shortURL2, originalURL2)
	s.expect.QueryRow(gomock.Any(), gomock.Any(), originalURL2).Return(pgxpoolmock.NewRow(shortURL2, originalURL2))
	s.expect.Exec(gomock.Any(), gomock.Any(), uid, shortURL2)
	ok, _, err = s.db.AddURL(s.ctx, originalURL2, shortURL2, uid)
	s.Assert().NoError(err)
	s.Assert().True(ok)

	s.expect.QueryRow(gomock.Any(), gomock.Any(), shortURL).Return(pgxpoolmock.NewRow(originalURL, false))
	original, ok, err := s.db.GetURL(s.ctx, shortURL)
	s.Assert().NoError(err)
	s.Assert().True(ok)
	s.Assert().Equal(originalURL, original)

	u1, err := url.JoinPath(s.config.Server.BaseURL, shortURL)
	s.Assert().NoError(err)
	u2, err := url.JoinPath(s.config.Server.BaseURL, shortURL2)
	s.Assert().NoError(err)
	expectedArray := []URLOfOwner{
		{
			u1,
			originalURL,
		},
		{
			u2,
			originalURL2,
		},
	}
	s.expect.Query(gomock.Any(), gomock.Any(), uid).Return(
		pgxpoolmock.NewRows([]string{"short_url", "original_url"}).
			AddRow(shortURL, originalURL).
			AddRow(shortURL2, originalURL2).ToPgxRows(),
		nil)
	array, err := s.db.GetURLByOwner(s.ctx, uid)
	s.Assert().NoError(err)
	s.Assert().Equal(expectedArray, array)
}

func (s *DBTestSuite) TestMultipleSameOriginals() {
	uid := generateString(20)

	originalURL := generateString(20)
	shortURL := generateString(10)
	s.expect.Exec(gomock.Any(), gomock.Any(), shortURL, originalURL)
	s.expect.QueryRow(gomock.Any(), gomock.Any(), originalURL).Return(pgxpoolmock.NewRow(shortURL, originalURL))
	s.expect.Exec(gomock.Any(), gomock.Any(), uid, shortURL)
	ok, _, err := s.db.AddURL(s.ctx, originalURL, shortURL, uid)
	s.Assert().NoError(err)
	s.Assert().True(ok)

	shortURL2 := generateString(10)
	s.expect.Exec(gomock.Any(), gomock.Any(), shortURL2, originalURL)
	s.expect.QueryRow(gomock.Any(), gomock.Any(), originalURL).Return(pgxpoolmock.NewRow(shortURL, originalURL))
	// s.expect.Exec(gomock.Any(), gomock.Any(), uid, shortURL2)
	ok, existing, err := s.db.AddURL(s.ctx, originalURL, shortURL2, uid)
	s.Assert().NoError(err)
	s.Assert().Equal(shortURL, existing)
	s.Assert().False(ok)

	s.expect.QueryRow(gomock.Any(), gomock.Any(), shortURL).Return(pgxpoolmock.NewRow(originalURL, false))
	result, ok, err := s.db.GetURL(s.ctx, shortURL)
	s.Assert().NoError(err)
	s.Assert().True(ok)
	s.Assert().Equal(result, originalURL)

	// in case it was deleted
	s.expect.QueryRow(gomock.Any(), gomock.Any(), shortURL).Return(pgxpoolmock.NewRow(originalURL, true))
	result, ok, err = s.db.GetURL(s.ctx, shortURL)
	s.Assert().NoError(err)
	s.Assert().True(ok)
	s.Assert().Equal(result, "")
}

func (s *DBTestSuite) TestAddConditions() {
	uid := generateString(20)

	originalURL := generateString(20)
	shortURL := generateString(10)
	s.expect.Exec(gomock.Any(), gomock.Any(), shortURL, originalURL).Return(nil, pgx.ErrTxClosed)
	ok, _, err := s.db.AddURL(s.ctx, originalURL, shortURL, uid)
	s.Assert().Error(err)
	s.Assert().False(ok)

	s.expect.Exec(gomock.Any(), gomock.Any(), shortURL, originalURL)
	s.expect.QueryRow(gomock.Any(), gomock.Any(), originalURL).
		Return(pgxpoolmock.NewRow(shortURL, originalURL).WithError(pgx.ErrTxClosed))
	ok, _, err = s.db.AddURL(s.ctx, originalURL, shortURL, uid)
	s.Assert().Error(err)
	s.Assert().False(ok)

	s.expect.Exec(gomock.Any(), gomock.Any(), shortURL, originalURL)
	s.expect.QueryRow(gomock.Any(), gomock.Any(), originalURL).Return(pgxpoolmock.NewRow(shortURL, originalURL))
	s.expect.Exec(gomock.Any(), gomock.Any(), uid, shortURL).Return(nil, pgx.ErrTxClosed)
	ok, _, err = s.db.AddURL(s.ctx, originalURL, shortURL, uid)
	s.Assert().Error(err)
	s.Assert().False(ok)
}

func TestDatabase(t *testing.T) {
	suite.Run(t, new(DBTestSuite))
}
