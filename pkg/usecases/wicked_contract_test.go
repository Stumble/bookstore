package usecases

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
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

// The 250ms timeout of SimpleCachedQuery bounds the statement in the database:
// a statement blocked by a lock fails with wpgx.ErrDBTimeout long before the
// caller's deadline.
func (s *usecaseTestSuite) TestWickedDatabaseTimeout() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	locker, err := s.GetPool().RawPrimaryPool().Begin(ctx)
	s.Require().NoError(err)
	defer func() { _ = locker.Rollback(context.Background()) }()
	_, err = locker.Exec(ctx, "LOCK TABLE books IN ACCESS EXCLUSIVE MODE")
	s.Require().NoError(err)

	start := time.Now()
	_, err = books.New(s.GetPool().WConn(), nil).SimpleCachedQuery(ctx)
	elapsed := time.Since(start)
	s.Require().ErrorIs(err, wpgx.ErrDBTimeout)
	s.ErrorIs(err, context.DeadlineExceeded)
	s.GreaterOrEqual(elapsed, 200*time.Millisecond)
	s.Less(elapsed, 2*time.Second)
	s.NoError(ctx.Err(), "the query timeout must fire before the caller's deadline")
}

// The 250ms timeout starts once the query has a connection: waiting longer than
// that for the pool's only connection does not fail the query.
func (s *usecaseTestSuite) TestWickedTimeoutExcludesPoolWait() {
	config := s.GetConfig()
	config.MaxConns = 1
	config.ReadReplicas = nil
	pool, err := wpgx.NewPool(context.Background(), &config)
	s.Require().NoError(err)
	s.T().Cleanup(pool.Close)
	held, err := pool.RawPrimaryPool().Acquire(context.Background())
	s.Require().NoError(err)
	time.AfterFunc(400*time.Millisecond, held.Release)

	start := time.Now()
	_, err = books.New(pool.WConn(), nil).SimpleCachedQuery(context.Background())
	s.Require().NoError(err)
	s.GreaterOrEqual(time.Since(start), 350*time.Millisecond)
}

// unresponsiveRedis accepts connections and never answers.
func (s *usecaseTestSuite) unresponsiveRedis() string {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	s.Require().NoError(err)
	var mu sync.Mutex
	var conns []net.Conn
	s.T().Cleanup(func() {
		_ = listener.Close()
		mu.Lock()
		defer mu.Unlock()
		for _, c := range conns {
			_ = c.Close()
		}
	})
	go func() {
		for {
			c, err := listener.Accept()
			if err != nil {
				return
			}
			mu.Lock()
			conns = append(conns, c)
			mu.Unlock()
		}
	}()
	return listener.Addr().String()
}

// Cache lookups are bounded by the Redis client's own timeouts, not by the
// 250ms query timeout: with dcache's default configuration (no Redis single
// flight) an unresponsive Redis fails the read after the client's ReadTimeout,
// well before the caller's deadline, without falling through to the database.
func (s *usecaseTestSuite) TestWickedCacheTimeout() {
	client := redis.NewClient(&redis.Options{
		Addr: s.unresponsiveRedis(), ReadTimeout: 400 * time.Millisecond, MaxRetries: -1,
	})
	s.T().Cleanup(func() { _ = client.Close() })
	hook := &getDeadlineHook{}
	client.AddHook(hook)
	cache, err := dcache.NewDCache("timeout", client, nil, 100*time.Millisecond, false, false)
	s.Require().NoError(err)
	s.T().Cleanup(cache.Close)
	conn := &observedConn{WGConn: s.GetPool().WConn()}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	start := time.Now()
	_, err = books.New(conn, cache).GetBookByID(ctx, 123)
	elapsed := time.Since(start)
	s.Require().Error(err)
	s.GreaterOrEqual(elapsed, 350*time.Millisecond, "the 250ms query timeout does not cut the cache lookup")
	s.Less(elapsed, 3*time.Second, "the client's ReadTimeout bounds the lookup")
	s.NoError(ctx.Err(), "the caller's deadline is not what ended it")
	s.Zero(conn.reads.Load(), "cache timeout must not fall through to a database read")
	// go-redis ignores context deadlines for socket I/O by default, so check the
	// deadline handed to Redis: the caller's, not the query timeout.
	s.Require().True(hook.seen.Load(), "the GET reached Redis")
	s.Greater(time.Duration(hook.remaining.Load()), 5*time.Second)
}

// getDeadlineHook records the time left before the deadline of the context of
// the first GET.
type getDeadlineHook struct {
	seen      atomic.Bool
	remaining atomic.Int64
}

func (h *getDeadlineHook) DialHook(next redis.DialHook) redis.DialHook { return next }

func (h *getDeadlineHook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		if cmd.Name() == "get" && !h.seen.Swap(true) {
			remaining := time.Duration(1<<63 - 1)
			if deadline, ok := ctx.Deadline(); ok {
				remaining = time.Until(deadline)
			}
			h.remaining.Store(int64(remaining))
		}
		return next(ctx, cmd)
	}
}

func (h *getDeadlineHook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return next
}
