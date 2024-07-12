package usecases

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/jackc/pgx/v5"

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
		act := "search"
		_ = u.activities.Insert(ctx, activities.InsertParams{
			Action:    act,
			Parameter: &bookName,
		}, &act)
	}()
	wg.Wait()
	return
}

func (u *Usecase) ListNewComicBookTx(ctx context.Context, bookName string, price float32) (id int, err error) {
	rst, err := u.pool.Transact(ctx, pgx.TxOptions{}, func(ctx context.Context, tx *wpgx.WTx) (any, error) {
		booksTx := u.books.WithTx(tx)
		activitiesTx := u.activities.WithTx(tx)
		id, err := booksTx.InsertAndReturnID(ctx, books.InsertAndReturnIDParams{
			Name:        bookName,
			Description: "book desc",
			Metadata:    []byte("{}"),
			Category:    books.BookCategoryComic,
			Price:       0.1,
		})
		if err != nil {
			return nil, err
		}
		if id == nil {
			return nil, fmt.Errorf("nil id?")
		}
		param := strconv.Itoa(int(*id))
		action := "list a new comic book"
		err = activitiesTx.Insert(ctx, activities.InsertParams{
			Action:    action,
			Parameter: &param,
		}, &action)
		return int(*id), err
	})
	if err != nil {
		return 0, err
	}
	return rst.(int), nil
}
