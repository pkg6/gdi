package gdi

import (
	"testing"
)

func TestSetGet(t *testing.T) {
	c := NewContainer()
	c.Set("foo", "bar")
	v := c.Get("foo")
	if v != "bar" {
		t.Fatal("Expected bar, got", v)
	}
}

type Foo struct {
	container IContainer
}

func (f *Foo) Construct(container IContainer) {
	f.container = container
}

func TestSingleTon(t *testing.T) {
	c := NewContainer()
	foo := &Foo{}
	c.Set("foo", foo)
	v1 := c.Get("foo")
	c.Set("foo", foo)
	v2 := c.Get("foo")
	if v1 != v2 {
		t.Fatal("Expected single instance, got ", v1, v2)
	}
}
func TestNotExistKey(t *testing.T) {
	c := NewContainer()
	defer func() {
		err := recover()
		if err == nil {
			t.Fatal("Expected panic, but nil")
		}
	}()
	_ = c.Get("not_exist")
}

type DemoServiceProvider struct {
}

func (d DemoServiceProvider) Register(container IContainer) {
	container.Set("foo", "bar")
}

func TestRegisterAndGet(t *testing.T) {
	c := NewContainer()
	sp := &DemoServiceProvider{}
	c.Register(sp)
	v := c.Get("foo")
	if v != "bar" {
		t.Fatal("Expected bar, got", v)
	}
}
func TestUnset(t *testing.T) {
	c := NewContainer()
	c.Set("foo", "bar")
	c.Unset("foo")
	defer func() {
		err := recover()
		if err == nil {
			t.Fatal("Expected panic after unset, but nil")
		}
	}()
	_ = c.Get("foo")
}
