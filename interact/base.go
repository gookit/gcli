// Package interact collect some interactive methods for CLI
package interact

type Interactive struct {
	Name string
}

func New(name string) *Interactive {
	return &Interactive{Name: name}
}
