package bus

type Event interface {
	Name() string
	Type() string
	Data() map[string]interface{}
}
