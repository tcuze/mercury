// Package registry is an interface for service discovery
package registry

import (
	"errors"
)

// 服务发现interface接口
type Registry interface {
	Init(...Option) error
	Options() Options
	Register(*Service, ...RegisterOption) error
	Deregister(*Service) error
	GetService(string) ([]*Service, error)
	ListServices() ([]*Service, error)
	Watch(...WatchOption) (Watcher, error)
	String() string
}

type Option func(*Options)

type RegisterOption func(*RegisterOptions)

type WatchOption func(*WatchOptions)

var (
	DEFAULT_REGISTRY = NewRegistry()
	// Not found error when GetService is called
	ERR_NOT_FOUND = errors.New("not found")
	// Watcher stopped error when watcher is stopped
	ERR_WATCHER_STOPPED = errors.New("watcher stopped")
)

// Register a service node. Additionally supply options such as TTL.
func Register(s *Service, opts ...RegisterOption) error {
	return DEFAULT_REGISTRY.Register(s, opts...)
}

// Deregister a service node
func Deregister(s *Service) error {
	return DEFAULT_REGISTRY.Deregister(s)
}

func GetService(name string) ([]*Service, error) {
	return DEFAULT_REGISTRY.GetService(name)
}

// List the services. Only returns service names
func ListServices() ([]*Service, error) {
	return DEFAULT_REGISTRY.ListServices()
}

// Watch returns a watcher which allows you to track updates to the
func Watch(opts ...WatchOption) (Watcher, error) {
	return DEFAULT_REGISTRY.Watch(opts...)
}

func String() string {
	return DEFAULT_REGISTRY.String()
}
