package gdi

var (
	di = New()
)

func Register(provider IServiceProvider) IContainer {
	return di.Register(provider)
}
func Handler(handlers ...HandlerFunc) IContainer {
	return di.Handler(handlers...)
}
func Set(id string, value any) error {
	return di.Set(id, value)
}
func MustGet(id string) any {
	return di.MustGet(id)
}
func Get(id string) (any, error) {
	return di.Get(id)
}
func Exists(id string) bool {
	return di.Exists(id)
}
func Unset(id string) {
	di.Unset(id)
}
func Raw(id string) (any, error) {
	return di.Raw(id)
}
func Values() map[string]any {
	return di.Values()
}
func Decode(val any) error {
	return di.Decode(val)
}
