package cmd

import (
	"went/registry"

	"github.com/zeroSal/went-clio/clio"
	"github.com/zeroSal/went-command/command"
)

var _ command.Interface = (*Init)(nil)
type Init struct {
	command.Base
}

func init() {
    registry.Command.Register(&Init{})
}

func (c *Init) GetHeader() command.Header {
	return command.Header{
		Use:   "init",
		Short: "Initialize a project",
		Long:  "Creates the base scaffolding for a Went project.",
	}
}

func (c *Init) Invoke() any {
	return c.run
}

func (c *Init) run(
	clio *clio.Clio,
) error {
	clio.Banner()

	clio.Success("The application works!")

	return nil
}
