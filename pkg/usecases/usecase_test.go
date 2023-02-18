package usecases

import (
	"context"
	"fmt"
	// "math/big"
	"testing"
	"time"

	// "github.com/jackc/pgx/v5/pgtype"
	"github.com/coocood/freecache"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/suite"
	"github.com/stumble/dcache"
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
	usecase   *Usecase
	RedisConn redis.UniversalClient
	FreeCache *freecache.Cache
	DCache    *dcache.DCache
}

func newUsecaseTestSuite() *usecaseTestSuite {
	redisClient := redis.NewClient(&redis.Options{
		Addr:        "127.0.0.1:6379",
		ReadTimeout: 3 * time.Second,
		PoolSize:    50,
		Password:    "",
	})
	if redisClient.Ping(context.Background()).Err() != nil {
		panic(fmt.Errorf("redis connection failed to ping"))
	}
	memCache := freecache.NewCache(100 * 1024 * 1024)
	dCache, err := dcache.NewDCache(
		"test", redisClient, memCache, 100*time.Millisecond, true, true)
	if err != nil {
		panic(err)
	}
	return &usecaseTestSuite{
		WPgxTestSuite: wpgxtestsuite.NewWPgxTestSuiteFromEnv("testdb", []string{
			books.Schema,
			activities.Schema,
		}),
		RedisConn: redisClient,
		FreeCache: memCache,
		DCache:    dCache,
	}
}

func TestUsecaseTestSuite(t *testing.T) {
	suite.Run(t, newUsecaseTestSuite())
}

func (suite *usecaseTestSuite) SetupTest() {
	suite.WPgxTestSuite.SetupTest()
	suite.Require().NoError(suite.RedisConn.FlushAll(context.Background()).Err())
	suite.FreeCache.Clear()
	suite.usecase = NewUsecase(suite.GetPool(), suite.DCache)
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

func (suite *usecaseTestSuite) TestGetBySpec() {
	ctx := context.Background()
	bookserde := booksTableSerde{books: suite.usecase.books}
	suite.LoadState("TestUsecaseTestSuite/TestGetBySpec.books.input.json", bookserde)
	cond := "b%"
	rst, err := suite.usecase.books.GetBookBySpec(ctx, books.GetBookBySpecParams{
		Name: &cond,
	})
	suite.Nil(err)
	suite.GoldenVarJSON("with_names", rst)

	d := int32(999)
	rst, err = suite.usecase.books.GetBookBySpec(ctx, books.GetBookBySpecParams{
		Dummy: &d,
	})
	suite.Nil(err)
	suite.GoldenVarJSON("with_dummy", rst)

	name := "a1%"
	namePtr := &name
	rst, err = suite.usecase.books.GetBookByNameMaybe(ctx, &name)
	suite.Nil(err)
	suite.GoldenVarJSON("by_name", rst)

	err = suite.usecase.books.InsertWithInvalidate(ctx,
		books.InsertWithInvalidateParams{
			ID:          123,
			Name:        "a1-2",
			Description: "xxxx",
			Metadata:    []byte("[]"),
			Category:    books.BookCategoryComputerScience,
			DummyField:  &d,
			Price:       1.69,
		},
		&namePtr,
		&books.GetBookBySpecParams{
			Dummy: &d,
		})
	suite.Require().Nil(err)

	rst, err = suite.usecase.books.GetBookBySpec(ctx, books.GetBookBySpecParams{
		Dummy: &d,
	})
	suite.Nil(err)
	suite.Equal(1, len(rst))
	rst[0].CreatedAt = time.Unix(0, 0).UTC()
	rst[0].UpdatedAt = time.Unix(0, 0).UTC()
	suite.GoldenVarJSON("with_dummy_2", rst)

	rst, err = suite.usecase.books.GetBookByNameMaybe(ctx, &name)
	suite.Nil(err)
	suite.Equal(2, len(rst))
	rst[1].CreatedAt = time.Unix(0, 0).UTC()
	rst[1].UpdatedAt = time.Unix(0, 0).UTC()
	suite.GoldenVarJSON("by_name_2", rst)

}
