package tx

import "context"

type Manager interface {
	DoInTx(ctx context.Context, fn func(context.Context) error) error
}
