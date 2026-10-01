package app

import (
	"context"
	"fmt"
	"net/http"
)

type Module interface {
	Name() string
	Register(*Container) error
}

type Container struct {
	mux     *http.ServeMux
	modules []Module
	startup []func(context.Context) error
}

func NewContainer() *Container {
	return &Container{mux: http.NewServeMux()}
}

func (c *Container) Mux() *http.ServeMux {
	return c.mux
}

func (c *Container) RegisterRoutes(fn func(*http.ServeMux)) {
	if fn != nil {
		fn(c.mux)
	}
}

func (c *Container) RegisterStartup(fn func(context.Context) error) {
	if fn != nil {
		c.startup = append(c.startup, fn)
	}
}

func (c *Container) RegisterModules(modules ...Module) error {
	for _, module := range modules {
		if module == nil {
			continue
		}
		if err := module.Register(c); err != nil {
			return fmt.Errorf("register module %s: %w", module.Name(), err)
		}
		c.modules = append(c.modules, module)
	}
	return nil
}

func (c *Container) Start(ctx context.Context) error {
	for _, fn := range c.startup {
		if err := fn(ctx); err != nil {
			return err
		}
	}
	return nil
}
