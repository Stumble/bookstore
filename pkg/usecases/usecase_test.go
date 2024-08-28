package usecases

import (
	"context"
	"fmt"
	"os"

	// "math/big"
	"testing"
	"time"

	// "github.com/jackc/pgx/v5/pgtype"
	"github.com/coocood/freecache"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/suite"
	"github.com/stumble/dcache"
	"github.com/stumble/wpgx"
	wpgxtestsuite "github.com/stumble/wpgx/testsuite"

	"github.com/stumble/bookstore/pkg/repos/activities"
	"github.com/stumble/bookstore/pkg/repos/books"
)

type booksTableSerde struct {
	books *books.Queries
}

func (b booksTableSerde) Load(data []byte) error {
	err := b.books.Load(context.Background(), data)
	if err != nil {
		return err
	}
	// reset books's id column to the next value.
	return b.books.RefreshIDSerial(context.Background())
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
	// setup "fake" replica database in envvar by using default configuration
	// that is almost the same as the primary database.
	// In real world, you will need to set up a real replica and set their parameters
	// in the envvar. For example:
	// os.Setenv("POSTGRES_REPLICAPREFIXES", "R1,R2")
	// os.Setenv("R1_NAME", "R1")
	// os.Setenv("R1_USERNAME", "user1")
	// os.Setenv("R1_PASSWORD", "password1")
	// os.Setenv("R1_HOST", "192.168.0.1")
	// ....
	// os.Setenv("R2_NAME", "R2")
	// ...
	os.Setenv("POSTGRES_REPLICAPREFIXES", "R1")
	os.Setenv("R1_NAME", "R1")
	os.Setenv("R1_DBNAME", "testdb")
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

func (suite *usecaseTestSuite) TestReplicaIsUp() {
	suite.Require().Equal(1, len(suite.Pool.ReplicaPools()))
	name := wpgx.ReplicaName("R1")

	// access 1 ok
	r1 := suite.Pool.ReplicaPools()[name]
	suite.Require().NotNil(r1)
	suite.Require().NoError(r1.Ping(context.Background()))

	// access 2 ok
	r1 = suite.Pool.MustReplicaPool(name)
	suite.Require().NoError(r1.Ping(context.Background()))

	// access 3 ok
	r1, ok := suite.Pool.ReplicaPool(name)
	suite.Require().True(ok)
	suite.Require().NoError(r1.Ping(context.Background()))
}

func (suite *usecaseTestSuite) TestReadFromReplica() {
	ctx := context.Background()
	bookserde := booksTableSerde{books: suite.usecase.books}
	suite.LoadState("TestUsecaseTestSuite/TestGetBySpec.books.input.json", bookserde)

	replicaName := wpgx.ReplicaName("R1")
	r1, err := suite.Pool.WQuerier(&replicaName)
	suite.Require().Nil(err)
	cond := "b%"
	rst, err := suite.usecase.books.UseReplica(r1).GetBookBySpec(ctx, books.GetBookBySpecParams{
		Name: &cond,
	})
	suite.Nil(err)
	for i := range rst {
		rst[i].CreatedAt = time.Unix(0, 0).UTC()
		rst[i].UpdatedAt = time.Unix(0, 0).UTC()
	}
	suite.GoldenVarJSON("with_names", rst)
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
	for i := range rst {
		rst[i].CreatedAt = time.Unix(0, 0).UTC()
		rst[i].UpdatedAt = time.Unix(0, 0).UTC()
	}
	suite.GoldenVarJSON("with_names", rst)

	d := int32(999)
	rst, err = suite.usecase.books.GetBookBySpec(ctx, books.GetBookBySpecParams{
		Dummy: &d,
	})
	suite.Nil(err)
	suite.Equal(0, len(rst))
	suite.GoldenVarJSON("with_dummy", rst)

	name := "a1%"
	namePtr := &name
	rst, err = suite.usecase.books.GetBookByNameMaybe(ctx, &name)
	suite.Nil(err)
	for i := range rst {
		rst[i].CreatedAt = time.Unix(0, 0).UTC()
		rst[i].UpdatedAt = time.Unix(0, 0).UTC()
	}
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
	for i := range rst {
		rst[i].CreatedAt = time.Unix(0, 0).UTC()
		rst[i].UpdatedAt = time.Unix(0, 0).UTC()
	}
	suite.GoldenVarJSON("with_dummy_2", rst)

	rst, err = suite.usecase.books.GetBookByNameMaybe(ctx, &name)
	suite.Nil(err)
	suite.Equal(2, len(rst))
	for i := range rst {
		rst[i].CreatedAt = time.Unix(0, 0).UTC()
		rst[i].UpdatedAt = time.Unix(0, 0).UTC()
	}
	suite.GoldenVarJSON("by_name_2", rst)

}

func (suite *usecaseTestSuite) TestNoPanicEvenDBErrorWhenNoCache() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel the context directly
	// This should cause the db query to fail due to context cancellation.
	b := suite.usecase.books.WithCache(nil)
	_, err := b.SimpleCachedQuery(ctx)
	suite.NotNil(err)
}

func (suite *usecaseTestSuite) TestListNewComicBook() {
	// must init after the last SetupTest()
	bookserde := booksTableSerde{books: suite.usecase.books}

	// load state
	suite.LoadState("TestUsecaseTestSuite/TestListNewComicBook.books.input.json", bookserde)

	allBooks, err := suite.usecase.books.GetAllBooks(context.Background())
	suite.Require().Nil(err)
	suite.Require().Equal(3, len(allBooks))

	// run search
	rst, err := suite.usecase.ListNewComicBookTx(context.Background(), "iron man 2", 10.88)

	// check return value
	suite.Nil(err)
	suite.Equal(4, rst)

	// test if invalidate cache works during a transaction.
	updatedAllBooks, err := suite.usecase.books.GetAllBooks(context.Background())
	suite.Require().Nil(err)
	suite.Require().Equal(4, len(updatedAllBooks))

	// verify db state
	suite.Golden("books_table", bookserde)
	suite.Golden("activitives_table", activitiesTableSerde{
		activities: suite.usecase.activities})
}
