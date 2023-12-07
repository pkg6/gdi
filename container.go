package gdi

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"sync"
)

const (
	ErrNotFind = iota + 1
	ErrExist
)

type HandlerFunc func(container IContainer)

type IContainer interface {
	Register(provider IServiceProvider) IContainer
	Handler(handlers ...HandlerFunc) IContainer
	Set(id string, value any) error
	MustGet(id string) any
	Get(id string) (any, error)
	Exists(id string) bool
	Unset(id string)
	Raw(id string) (any, error)
	Values() map[string]any
	Decode(val any) error
}

type ContainerErr struct {
	Index   int
	ID      string
	Message string
}

func NewErr(index int, ID string) *ContainerErr {
	e := &ContainerErr{ID: ID, Index: index}
	switch index {
	case ErrNotFind:
		e.Message = fmt.Sprintf("identifier %s is not defined", ID)
	case ErrExist:
		e.Message = fmt.Sprintf("cannot override frozen service %s", ID)
	}
	return e
}

func (e *ContainerErr) Error() string {
	return e.Message
}

type Container struct {
	lock   sync.RWMutex
	values map[string]any
	raw    map[string]any
	frozen map[string]bool
}

func New() IContainer {
	return &Container{
		lock:   sync.RWMutex{},
		frozen: map[string]bool{},
		values: map[string]any{},
		raw:    map[string]any{},
	}
}
func (c *Container) Register(provider IServiceProvider) IContainer {
	c.lock.RLock()
	defer c.lock.RUnlock()
	provider.Register(c)
	return c
}

func (c *Container) Handler(handlers ...HandlerFunc) IContainer {
	for _, handler := range handlers {
		handler(c)
	}
	return c
}

func (c *Container) Set(id string, value any) error {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if _, frozenExist := c.frozen[id]; frozenExist {
		return NewErr(ErrExist, id)
	}
	c.values[id] = value
	return nil
}
func (c *Container) MustGet(id string) any {
	val, err := c.Get(id)
	if err != nil {
		panic(err)
	}
	return val
}
func (c *Container) Get(id string) (any, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	raw, keyExist := c.values[id]
	if !keyExist {
		return nil, NewErr(ErrNotFind, id)
	}
	if object, ok := raw.(IObject); ok {
		object.Construct(c)
		return object, nil
	}
	if handler, ok := raw.(func(IContainer) any); ok {
		return handler(c), nil
	}
	c.raw[id] = raw
	c.frozen[id] = true
	return c.raw[id], nil
}

func (c *Container) Exists(id string) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if _, keyExist := c.values[id]; keyExist {
		return true
	}
	return false
}

func (c *Container) Unset(id string) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	delete(c.frozen, id)
	delete(c.values, id)
	delete(c.raw, id)
}

func (c *Container) Raw(id string) (any, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	rawVal, keyExist := c.values[id]
	if !keyExist {
		return nil, NewErr(ErrNotFind, id)
	}
	return rawVal, nil
}

func (c *Container) Values() map[string]any {
	return c.values
}

func (c *Container) Decode(val any) error {
	return mapstructure.Decode(c.Values(), val)
}
