package bootstrap

import (
	"context"
	_const "source-base/internal/const"
	"source-base/internal/ent/generate"
	"source-base/internal/ports"
	"source-base/pkg/database/ent"
)

func NewEntTxManager(client *generate.Client) ports.TxManager {
	return ent.NewTxManager(
		client.Tx,
		func(ctx context.Context, tx *generate.Tx) context.Context {
			return context.WithValue(ctx, _const.TxKey, tx)
		},
	)
}
