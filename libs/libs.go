package libs

import (
	"github.com/PhantomX7/dhamma/libs/casbin"
	"github.com/PhantomX7/dhamma/libs/transaction_manager"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		transaction_manager.New,
		casbin.New,
	),
)
