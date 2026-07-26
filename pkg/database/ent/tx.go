package ent

import (
	"context"
	"fmt"

	"github.com/himbo22/source-base/pkg/common/tx"
)

type Tx interface {
	Commit() error
	Rollback() error
}

func WithTx[T Tx](ctx context.Context, begin func(context.Context) (T, error), fn func(T) error) error {
	t, err := begin(ctx)
	if err != nil {
		return fmt.Errorf("begin Tx: %w", err)
	}
	defer func() {
		if v := recover(); v != nil {
			_ = t.Rollback()
			panic(v)
		}
	}()
	if err := fn(t); err != nil {
		if rerr := t.Rollback(); rerr != nil {
			return fmt.Errorf("rollback: %w", rerr)
		}
		return err
	}
	if err := t.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}
	return nil
}

type txManager[T Tx] struct {
	begin  func(context.Context) (T, error)
	inject func(context.Context, T) context.Context
}

func NewTxManager[T Tx](
	begin func(context.Context) (T, error),
	inject func(context.Context, T) context.Context,
) tx.Manager {
	return txManager[T]{
		begin:  begin,
		inject: inject,
	}
}

func (m txManager[T]) DoInTx(ctx context.Context, fn func(context.Context) error) error {
	return WithTx(ctx, m.begin, func(t T) error {
		return fn(m.inject(ctx, t))
	})
}
