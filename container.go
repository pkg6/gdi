package gdi

import (
	"fmt"
	"sync"
)

type IContainer interface {
	Register(provider IServiceProvider) IContainer
	Set(id string, value any)
	Get(id string) any
	Exists(id string) bool
	Unset(id string)
	Raw(id string) (any, error)
}

type Container struct {
	lock   sync.RWMutex
	values map[string]any
	raw    map[string]any
	frozen map[string]bool
}

func NewContainer() IContainer {
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

func (c *Container) Set(id string, value any) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if _, frozenExist := c.frozen[id]; frozenExist {
		panic(fmt.Errorf("cannot override frozen service %s", id))
	}
	c.values[id] = value
}

func (c *Container) Get(id string) any {
	c.lock.RLock()
	defer c.lock.RUnlock()
	raw, keyExist := c.values[id]
	if !keyExist {
		panic(fmt.Errorf("identifier %s is not defined", id))
	}
	if object, ok := raw.(IObject); ok {
		object.Construct(c)
		return object
	}
	if handler, ok := raw.(func(IContainer) any); ok {
		return handler(c)
	}
	c.raw[id] = raw
	c.frozen[id] = true
	return c.raw[id]
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
		panic(fmt.Errorf("identifier %s is not defined", id))
	}
	return rawVal, nil
}