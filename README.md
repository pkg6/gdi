A simple yet powerful dependency injection container for Go.

## Installation

```
go get -u github.com/pkg6/gdi
```

## Usage

Create a container

```
container := container.New()
```

##  Set a value

```
container.Set("foo", "bar")
```

##  Get a value

```
value := container.Get("foo") // value is "bar"
```

## Check if a value exists

```
exists := container.Exists("foo") // exists is true
```

##  Unset a value

```
container.Unset("foo")
```

## Register a service provider

```
type DemoServiceProvider struct {
}
func (d DemoServiceProvider) Register(container IContainer) {
	container.Set("foo", "bar")
}
sp := &DemoServiceProvider{}
container.Register(sp)
```

## Register a Object

```
type Foo struct {
	container IContainer
}

func (f *Foo) Construct(container IContainer) {
	f.container = container
}

foo := &Foo{}
c.Set("foo", foo)
```

The service provider can then set values in the container

## Register a HandlerFunc

```
c.Handler(func(container IContainer) {
  container.Set("handler", "handler")
})
```

The service provider can then set values in the container

## Panic handling

Methods like Get and Raw will panic if the identifier does not exist. You can recover from panics to handle the error in a graceful way:

```
defer func() {
  if err := recover(); err != nil {
    // handle error
  }
}()
value = container.Get("not-exist")  // Will panic
```

## Contributors

[@pkg6](https://github.com/pkg6/gdi)

## License

Container is licensed under the MIT license. See the [LICENSE](https://github.com/pkg6/gdi/blob/main/LICENSE) file for details.
