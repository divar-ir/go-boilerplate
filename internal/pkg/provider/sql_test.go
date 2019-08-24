package provider_test

import (
	"context"
	"io"
	"testing"

	"github.com/cafebazaar/go-boilerplate/internal/pkg/provider"
	"github.com/cafebazaar/go-boilerplate/pkg/postview"
	"github.com/cafebazaar/go-boilerplate/pkg/sql"
	"github.com/stretchr/testify/suite"
	"golang.org/x/xerrors"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type SQLProviderTestSuite struct {
	suite.Suite

	provider provider.PostProvider
}

func TestSQLProviderTestSuite(t *testing.T) {
	suite.Run(t, new(SQLProviderTestSuite))
}

func (s *SQLProviderTestSuite) TestGetPostShouldReturnNotFoundInitially() {
	_, err := s.provider.GetPost(context.Background(), "myToken")
	s.True(xerrors.Is(err, provider.ErrNotFound))
}

func (s *SQLProviderTestSuite) TestShouldReturnPostAfterAdd() {
	err := s.provider.AddPost(context.Background(), &postview.Post{
		Token: "abcd",
		Title: "a title",
	})
	s.Nil(err)
	if err != nil {
		return
	}

	post, err := s.provider.GetPost(context.Background(), "abcd")
	s.Nil(err)
	if err != nil {
		return
	}

	s.Equal("abcd", post.Token)
	s.Equal("a title", post.Title)
}

func (s *SQLProviderTestSuite) SetupTest() {
	db, err := sql.GetDatabase(sql.SqliteConfig{
		InMemory: true,
	})
	if err != nil {
		s.FailNow(err.Error(), "unable to instantiate SQLite instance")
		return
	}

	s.provider = provider.NewSQL(db)

	err = s.provider.(sql.Migrater).Migrate()
	if err != nil {
		s.FailNow(err.Error(), "unable to migrate SQLite database")
	}
}

func (s *SQLProviderTestSuite) TearDownTest() {
	err := s.provider.(io.Closer).Close()
	if err != nil {
		s.FailNow(err.Error(), "unable to close SQLite database")
	}
}
