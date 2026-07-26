package repository

import (
	"context"
	_const "github.com/himbo22/source-base/internal/const"
	"github.com/himbo22/source-base/internal/ent/generate"
)

func GetClient(ctx context.Context, client *generate.Client) *generate.Client {
	tx, ok := ctx.Value(_const.TxKey).(*generate.Tx)
	if !ok || tx == nil {
		return client
	}
	return tx.Client()
}
