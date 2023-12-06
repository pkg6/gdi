package gdi

import (
	"testing"
)

func TestSetGet(t *testing.T) {
	c := New()
	c.Set("foo", "bar")
	v, _ := c.Get("foo")
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
	c := New()
	foo := &Foo{}
	c.Set("foo", foo)
	v1, _ := c.Get("foo")
	c.Set("foo", foo)
	v2, _ := c.Get("foo")
	if v1 != v2 {
		t.Fatal("Expected single instance, got ", v1, v2)
	}
}
func TestNotExistKey(t *testing.T) {
	c := New()
	_, err := c.Get("not_exist")
	if e, ok := err.(*ContainerErr); ok {
		if e.Index != ErrNotFind {
			t.Fatal("not_exist not defined ")
		}
	}
}

type DemoServiceProvider struct {
}

func (d DemoServiceProvider) Register(container IContainer) {
	container.Set("foo", "bar")
}

func TestRegisterAndGet(t *testing.T) {
	c := New()
	sp := &DemoServiceProvider{}
	c.Register(sp)
	v, _ := c.Get("foo")
	if v != "bar" {
		t.Fatal("Expected bar, got", v)
	}
}
func TestHandler(t *testing.T) {
	c := New()
	c.Handler(func(container IContainer) {
		container.Set("handler", "handler")
	})
	v, _ := c.Get("handler")
	if v != "handler" {
		t.Fatal("Expected bar, got", v)
	}
}
func TestUnset(t *testing.T) {
	c := New()
	c.Set("foo", "bar")
	c.Unset("foo")

	_, err := c.Get("foo")
	if e, ok := err.(*ContainerErr); ok {
		if e.Index != ErrNotFind {
			t.Fatal("not_exist not defined ")
		}
	}
}

type app struct {
	Demo *demo `json:"demo"`
}

type demo struct {
	P string
}

func TestUnmarshal(t *testing.T) {
	c := New()
	var apps app
	c.Set("demo", &demo{P: "p"})
	c.Decode(&apps)
	if apps.Demo.P != "p" {
		t.Fatal("Expected Unmarshal")
	}
}
