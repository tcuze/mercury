package stream

type Stream interface {
	Init(opts ...Option) error
	Register() error
	Deregister() error
	Options() Options
	Start() error
	Stop() error
}

type Option func(*Options)

var (
	DEFAULT_SERVER = NewStream()
)

func Register() error {
	return DEFAULT_SERVER.Register()
}

func Deregister() error {
	return DEFAULT_SERVER.Deregister()
}

func Start() error {
	return DEFAULT_SERVER.Start()
}

func Stop() error {
	return DEFAULT_SERVER.Stop()
}
