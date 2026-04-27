package app

import (
	"context"
	"embed"
	"went/app/bootstrap"

	"github.com/zeroSal/went-clio/clio"
	"go.uber.org/fx"
)

type Kernel struct {
	EmbedFS embed.FS
	Specs   *Specs
	Clio    *clio.Clio
}

func NewKernel(
	embedFS embed.FS,
	specs *Specs,
	clio *clio.Clio,
) *Kernel {
	return &Kernel{
		embedFS,
		specs,
		clio,
	}
}

func (a *Kernel) Run(invoke any, opts ...fx.Option) error {
	di := []fx.Option{
		Container,
		bootstrap.Init,
		fx.Supply(a.Clio),
		fx.Supply(a.EmbedFS),
		fx.Supply(a.Specs),
		fx.Invoke(invoke),
		fx.NopLogger,
	}

	app := fx.New(append(di, opts...)...)

	if err := app.Start(context.Background()); err != nil {
		return err
	}

	if err := app.Stop(context.Background()); err != nil {
		return err
	}

	return nil
}
