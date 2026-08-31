package bootstrap

import (
	"context"

	_const "github.com/himbo22/source-base/internal/const"
	"github.com/himbo22/source-base/internal/ent/generate"
	"github.com/himbo22/source-base/internal/ports"
	"github.com/himbo22/source-base/pkg/database/ent"
)

func NewEntTxManager(client *generate.Client) ports.TxManager {
	return ent.NewTxManager(
		client.Tx,
		func(ctx context.Context, tx *generate.Tx) context.Context {
			return context.WithValue(ctx, _const.TxKey, tx)
		},
	)
}
