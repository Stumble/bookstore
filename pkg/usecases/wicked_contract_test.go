package usecases

import (
	"context"
	"encoding/json"
	"errors"
	"sync/atomic"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5"
	"github.com/stumble/bookstore/pkg/repos/activities"
	"github.com/stumble/bookstore/pkg/repos/books"
	"github.com/stumble/dcache"
	"github.com/stumble/wpgx"
)

type observedConn struct {
	wpgx.WGConn
	reads   atomic.Int64
	intents atomic.Int64
}

func (s *usecaseTestSuite) TestWickedCopyFrom() {
	ctx := context.Background()
	parameter := "copied"
	q := s.usecase.activities.WithCache(nil)
	count, err := q.BulkInsert(ctx, []activities.BulkInsertParams{
		{Action: "bulk", Parameter: &parameter, CreatedAt: time.Unix(100, 0).UTC()},
		{Action: "bulk", CreatedAt: time.Unix(200, 0).UTC()},
	})
	s.Require().NoError(err)
	s.Equal(int64(2), count)
	rows, err := q.GetActivitiesByAction(ctx, "bulk")
	s.Require().NoError(err)
	s.Require().Len(rows, 2)
	parameters := map[int64]*string{}
	for _, row := range rows {
		parameters[row.CreatedAt.Unix()] = row.Parameter
	}
	s.Require().NotNil(parameters[100])
	s.Equal(parameter, *parameters[100])
	s.Nil(parameters[200])
}

func (c *observedConn) WQuery(ctx context.Context, name, sql string, args ...interface{}) (pgx.Rows, error) {
	c.reads.Add(1)
	return c.WGConn.WQuery(ctx, name, sql, args...)
}
func (c *observedConn) WQueryRow(ctx context.Context, name, sql string, args ...interface{}) pgx.Row {
	c.reads.Add(1)
	return c.WGConn.WQueryRow(ctx, name, sql, args...)
}
func (c *observedConn) CountIntent(name string) { c.intents.Add(1); c.WGConn.CountIntent(name) }

func (s *usecaseTestSuite) seedContractBook() int32 {
	id, err := s.usecase.books.WithCache(nil).InsertAndReturnID(context.Background(), books.InsertAndReturnIDParams{
		Name: "contract", Description: "original", Metadata: json.RawMessage(`{"edition":1}`), Category: books.BookCategoryComic, Price: 10,
	})
	s.Require().NoError(err)
	s.Require().NotNil(id)
	return *id
}

func (s *usecaseTestSuite) TestWickedTransactionInvalidation() {
	ctx := context.Background()
	id := s.seedContractBook()
	conn := &observedConn{WGConn: s.GetPool().WConn()}
	q := books.New(conn, s.DCache)
	before, err := q.GetBookByID(ctx, id)
	s.Require().NoError(err)
	s.Require().NotNil(before)
	_, err = q.GetBookByID(ctx, id)
	s.Require().NoError(err)
	s.Equal(int64(1), conn.reads.Load(), "second read must hit the cache")
	s.Equal(int64(2), conn.intents.Load(), "cache hits still count query intent")
	abort := errors.New("intentional rollback")
	mutate := func(rollback bool) error {
		_, err := s.GetPool().Transact(ctx, pgx.TxOptions{}, func(ctx context.Context, tx *wpgx.WTx) (any, error) {
			err := q.WithTx(tx).UpdateBookByID(ctx, books.UpdateBookByIDParams{ID: id, Description: "changed", Meta: before.Metadata, Price: before.Price}, &id)
			if err != nil {
				return nil, err
			}
			cached, err := q.GetBookByID(ctx, id)
			s.Require().NoError(err)
			s.Equal("original", cached.Description)
			s.Equal(int64(1), conn.reads.Load(), "invalidation must wait for commit")
			if rollback {
				return nil, abort
			}
			return nil, nil
		})
		return err
	}
	s.Require().ErrorIs(mutate(true), abort)
	after, err := q.GetBookByID(ctx, id)
	s.Require().NoError(err)
	s.Equal("original", after.Description)
	s.Equal(int64(1), conn.reads.Load(), "rollback must leave the cached value intact")
	uncached, err := q.WithCache(nil).GetBookByID(ctx, id)
	s.Require().NoError(err)
	s.Equal("original", uncached.Description)
	// Reset only the observer; the cache remains populated with the original row.
	conn.reads.Store(1)
	s.Require().NoError(mutate(false))
	after, err = q.GetBookByID(ctx, id)
	s.Require().NoError(err)
	s.Equal("changed", after.Description)
	s.Equal(int64(2), conn.reads.Load(), "committed mutation must invalidate and cause one new read")
}

func (s *usecaseTestSuite) TestWickedMultipleNoArgumentInvalidations() {
	ctx := context.Background()
	s.seedContractBook()
	conn := &observedConn{WGConn: s.GetPool().WConn()}
	q := books.New(conn, s.DCache)
	a, err := q.GetAllBooks(ctx)
	s.Require().NoError(err)
	s.Len(a, 1)
	b, err := q.GetAllBooks2(ctx)
	s.Require().NoError(err)
	s.Len(b, 1)
	_, err = q.InsertAndReturnID(ctx, books.InsertAndReturnIDParams{Name: "second", Description: "new", Metadata: json.RawMessage(`[]`), Category: books.BookCategoryComic, Price: 2})
	s.Require().NoError(err)
	a, err = q.GetAllBooks(ctx)
	s.Require().NoError(err)
	s.Len(a, 2)
	b, err = q.GetAllBooks2(ctx)
	s.Require().NoError(err)
	s.Len(b, 2)
	s.Equal(int64(5), conn.reads.Load(), "two initial reads, INSERT RETURNING, and two invalidated reads")
}

func (s *usecaseTestSuite) TestWickedNilCacheAndMissingRows() {
	ctx := context.Background()
	id := s.seedContractBook()
	q := s.usecase.books.WithCache(nil)
	s.Require().NoError(q.UpdateBookByID(ctx, books.UpdateBookByIDParams{ID: id, Description: "uncached", Meta: json.RawMessage(`{}`), Price: 3}, &id))
	row, err := q.GetBookByID(ctx, id)
	s.Require().NoError(err)
	s.Equal("uncached", row.Description)
	missing, err := q.GetBookByID(ctx, 999999)
	s.Require().NoError(err)
	s.Nil(missing)
	conn := &observedConn{WGConn: s.GetPool().WConn()}
	cached := books.New(conn, s.DCache)
	for i := 0; i < 2; i++ {
		missing, err = cached.GetBookByID(ctx, 999999)
		s.Require().NoError(err)
		s.Nil(missing)
	}
	s.Equal(int64(1), conn.reads.Load(), "a missing row is cacheable")
}

func (s *usecaseTestSuite) TestWickedPointerKeysAndJSON() {
	ctx := context.Background()
	s.seedContractBook()
	conn := &observedConn{WGConn: s.GetPool().WConn()}
	q := books.New(conn, s.DCache)
	first, second := "contract", "contract"
	rows, err := q.GetBookByNameMaybe(ctx, &first)
	s.Require().NoError(err)
	s.Require().Len(rows, 1)
	s.JSONEq(`{"edition":1}`, string(rows[0].Metadata))
	rows, err = q.GetBookByNameMaybe(ctx, &second)
	s.Require().NoError(err)
	s.Require().Len(rows, 1)
	s.Equal(int64(1), conn.reads.Load(), "pointer arguments use values, not addresses, as cache keys")
	name := "contract%"
	dummy := int32(7)
	s.Equal("books:GetBookBySpec:contract%,<nil>,7", (books.GetBookBySpecParams{Name: &name, Dummy: &dummy}).CacheKey())
}

type deadlineConn struct {
	wpgx.WGConn
	remaining time.Duration
}

func (c *deadlineConn) WQuery(ctx context.Context, _ string, _ string, _ ...interface{}) (pgx.Rows, error) {
	deadline, _ := ctx.Deadline()
	c.remaining = time.Until(deadline)
	<-ctx.Done()
	return nil, ctx.Err()
}

func (s *usecaseTestSuite) TestWickedDatabaseTimeout() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	conn := &deadlineConn{WGConn: s.GetPool().WConn()}
	_, err := books.New(conn, nil).SimpleCachedQuery(ctx)
	s.Require().ErrorIs(err, context.DeadlineExceeded)
	s.Greater(conn.remaining, time.Duration(0))
	s.LessOrEqual(conn.remaining, 250*time.Millisecond)
	s.NoError(ctx.Err(), "the per-query deadline must fire before the parent deadline")
}

type deadlineHook struct{ remaining atomic.Int64 }

func (h *deadlineHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	if cmd.Name() == "get" {
		deadline, _ := ctx.Deadline()
		h.remaining.Store(int64(time.Until(deadline)))
		<-ctx.Done()
		return ctx, ctx.Err()
	}
	return ctx, nil
}
func (h *deadlineHook) AfterProcess(context.Context, redis.Cmder) error { return nil }
func (h *deadlineHook) BeforeProcessPipeline(ctx context.Context, _ []redis.Cmder) (context.Context, error) {
	return ctx, nil
}
func (h *deadlineHook) AfterProcessPipeline(context.Context, []redis.Cmder) error { return nil }

func (s *usecaseTestSuite) TestWickedCacheTimeout() {
	client := redis.NewClient(&redis.Options{Addr: s.RedisConn.(*redis.Client).Options().Addr})
	s.T().Cleanup(func() { _ = client.Close() })
	hook := &deadlineHook{}
	client.AddHook(hook)
	cache, err := dcache.NewDCache("timeout", client, nil, 100*time.Millisecond, false, false)
	s.Require().NoError(err)
	s.T().Cleanup(cache.Close)
	conn := &observedConn{WGConn: s.GetPool().WConn()}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err = books.New(conn, cache).GetBookByID(ctx, 123)
	s.Require().Error(err)
	s.Greater(hook.remaining.Load(), int64(0))
	s.LessOrEqual(hook.remaining.Load(), int64(250*time.Millisecond))
	s.Zero(conn.reads.Load(), "cache timeout must not fall through to a database read")
	s.NoError(ctx.Err())
}
