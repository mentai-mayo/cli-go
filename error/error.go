package error

import (
	"fmt"
)

type Error struct {
	name    string
	message string
}

func New(name string, message string) Error {
	return Error{name, message}
}

func NewFormat(name string, format string, a ...any) Error {
	message := fmt.Sprintf(format, a...)
	return Error{name, message}
}

func (e Error) Error() string {
	return fmt.Sprintf("\x1B[31m%s Error\x1B[0m: %s", e.name, e.message)
}

func (e Error) Name() string {
	return e.name
}
