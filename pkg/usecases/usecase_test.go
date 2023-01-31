package usecases

import (
	"context"
	// "math/big"
	"testing"
	"time"

	// "github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/suite"
	wpgxtestsuite "github.com/stumble/wpgx/testsuite"

	"github.com/stumble/bookstore/pkg/repos/activities"
	"github.com/stumble/bookstore/pkg/repos/books"
)

type booksTableSerde struct {
	books *books.Queries
}

func (b booksTableSerde) Load(data []byte) error {
	return b.books.Load(context.Background(), data)
}

func (b booksTableSerde) Dump() ([]byte, error) {
	return b.books.Dump(context.Background(), func(m *books.Book) {
		m.CreatedAt = time.Unix(0, 0).UTC()
		m.UpdatedAt = time.Unix(0, 0).UTC()
	})
}

type activitiesTableSerde struct {
	activities *activities.Queries
}

func (a activitiesTableSerde) Load(data []byte) error {
	return a.activities.Load(context.Background(), data)
}

func (a activitiesTableSerde) Dump() ([]byte, error) {
	return a.activities.Dump(context.Background(), func(m *activities.Activity) {
		m.CreatedAt = time.Unix(0, 0).UTC()
	})
}

type usecaseTestSuite struct {
	*wpgxtestsuite.WPgxTestSuite
	usecase *Usecase
}

func newUsecaseTestSuite() *usecaseTestSuite {
	return &usecaseTestSuite{
		WPgxTestSuite: wpgxtestsuite.NewWPgxTestSuiteFromEnv("testdb", []string{
			books.Schema,
			activities.Schema,
		}),
	}
}

func TestUsecaseTestSuite(t *testing.T) {
	suite.Run(t, newUsecaseTestSuite())
}

func (suite *usecaseTestSuite) SetupTest() {
	suite.WPgxTestSuite.SetupTest()
	suite.usecase = NewUsecase(suite.GetPool())
}

func (suite *usecaseTestSuite) TestBulkInsertBooks() {
	err := suite.usecase.books.Insert(context.Background(),
		books.InsertParams{
			Name:        "DDIA",
			Description: "backend must read!",
			Metadata:    nil,
			Category:    books.BookCategoryComputerScience,
			// Price:       pgtype.Numeric{Valid: true, Int: big.NewInt(3888), Exp: -2},
			Price: 38.88,
		},
	)
	suite.Require().Nil(err)
	err = suite.usecase.books.Insert(context.Background(),
		books.InsertParams{
			Name:        "HAHAHAH",
			Description: "very funny",
			Metadata:    nil,
			Category:    books.BookCategoryComic,
			// Price:       pgtype.Numeric{Valid: true, Int: big.NewInt(55), Exp: -1},
			Price: 5.5,
		},
	)
	suite.Require().Nil(err)
	suite.Golden("books", booksTableSerde{books: suite.usecase.books})
}

func (suite *usecaseTestSuite) TestSearch() {
	for _, tc := range []struct {
		tcName      string
		s           string
		n           int
		expectedErr error
	}{
		{"start_with_a", "a", 1, nil},
		{"start_with_b", "b", 2, nil},
	} {
		suite.Run(tc.tcName, func() {
			suite.SetupTest()

			// must init after the last SetupTest()
			bookserde := booksTableSerde{books: suite.usecase.books}
			// load state
			suite.LoadState("TestUsecaseTestSuite/TestSearch.books.input.json", bookserde)

			// run search
			rst, err := suite.usecase.Search(context.Background(), tc.s)

			// check return value
			suite.Equal(tc.expectedErr, err)
			suite.Equal(tc.n, rst)

			// verify db state
			suite.Golden("books_table", bookserde)
			suite.Golden("activitives_table", activitiesTableSerde{
				activities: suite.usecase.activities})
		})
	}
}
