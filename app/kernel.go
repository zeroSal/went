package app

import (
	"context"
	"embed"
	"went/app/bootstrap"

	"github.com/zeroSal/went-clio/clio"
	"go.uber.org/fx"
)

type Kernel struct {
	Context context.Context
	EmbedFS embed.FS
	Specs   *Specs
	Clio    *clio.Clio
}

func NewKernel(
	context context.Context,
	embedFS embed.FS,
	specs *Specs,
	clio *clio.Clio,
) *Kernel {
	return &Kernel{
		context,
		embedFS,
		specs,
		clio,
	}
}

func (a *Kernel) Run(invoke any, opts ...fx.Option) error {
	contextProvider := func () context.Context {
		return a.Context
	}

	di := []fx.Option{
		Container,
		bootstrap.Init,
		fx.Supply(a.Clio),
		fx.Provide(contextProvider),
		fx.Invoke(invoke),
		fx.NopLogger,
	}

	app := fx.New(append(di, opts...)...)

	if err := app.Start(a.Context); err != nil {
		return err
	}

	if err := app.Stop(a.Context); err != nil {
		return err
	}

	return nil
}
