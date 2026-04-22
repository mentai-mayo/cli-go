package expects

import (
	"fmt"
	"reflect"
	"unicode"

	clierr "github.com/mentai-mayo/cli-go/error"
)

type Expects struct {
	inner []Expect
}

func New[T any]() (*Expects, error) {
	rtype := reflect.TypeOf((*T)(nil)).Elem()

	if rtype.Kind() != reflect.Struct {
		return nil, clierr.New("Non-struct Target", fmt.Sprintf("arguments cannot be parsed for non-struct types (got: %s)", rtype.Kind()))
	}

	for i := 0; i < rtype.NumField(); i++ {
		rfield := rtype.Field(i)

		if unicode.IsLower(rune(rfield.Name[0])) {
			// field is private
			continue
		}
	}
}

type Expect interface {
	Set(value any) bool
	Position() int
}

type ExpectOption struct {
	target uintptr
}

func (e *ExpectOption) Set(value any) bool {
	e.value
}

type ExpectType uint

func (e ExpectType) String() string {
	switch uint(e) {
	case 0:
		return "bool"
	case 1:
		return "string"
	case 2:
		return "int"
	}
}

const (
	BOOL   ExpectType = iota
	STRING ExpectType = iota
	INT    ExpectType = iota
)
