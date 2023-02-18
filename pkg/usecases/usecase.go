package usecases

import (
	"context"
	"sync"

	"github.com/stumble/bookstore/pkg/repos/activities"
	"github.com/stumble/bookstore/pkg/repos/books"
	"github.com/stumble/dcache"
	"github.com/stumble/wpgx"
)

type Usecase struct {
	books      *books.Queries
	activities *activities.Queries

	pool *wpgx.Pool
}

func NewUsecase(pool *wpgx.Pool, cache *dcache.DCache) *Usecase {
	return &Usecase{
		books:      books.New(pool.WConn(), cache),
		activities: activities.New(pool.WConn(), cache),
		pool:       pool,
	}
}

func (u *Usecase) Search(ctx context.Context, bookName string) (n int, err error) {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		rst, searchErr := u.books.SearchBooks(ctx, bookName+"%")
		if searchErr != nil {
			n = 0
			err = searchErr
		} else {
			n = len(rst)
		}
	}()
	go func() {
		defer wg.Done()
		_ = u.activities.Insert(ctx, activities.InsertParams{
			Action:    "search",
			Parameter: &bookName,
		})
	}()
	wg.Wait()
	return
}
