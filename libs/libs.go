package libs

import (
	"github.com/PhantomX7/dhamma/libs/database_transaction"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		database_transaction.New,
	),
)
