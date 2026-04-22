package cli

import (
	"fmt"
	"reflect"

	clierr "github.com/mentai-mayo/cli-go/error"
	expect "github.com/mentai-mayo/cli-go/expect"
)

type Parser[T any] struct {
	//
}

func Parse[T any](list []string) (*T, error) {

	expects, err := expect.New[T]()
	if err != nil {
		return nil, err
	}

	// reflect.Type[T]
	rtype := reflect.TypeOf((*T)(nil)).Elem()

	// reflect.Value[*T]
	target := reflect.New(rtype)

	index := 0
	next := func() (string, bool) {
		if index >= len(list) {
			return "", false
		}
		arg := list[index]
		index += 1
		return arg, true
	}

	optend := false
	for {
		arg, ok := next()
		if !ok {
			break
		}

		// end of option
		if arg == "--" {
			optend = true
			continue
		}

		// long option
		if !optend && arg[0:2] == "--" {
			info := expects.GetOpt(arg[2:])
			field := target.FieldByName(info.FieldName())
			switch info.EType() {
			case expects.BOOL:
				field.SetBool(true)
			case expects.STRING:
				value, ok := next()
				if !ok {
					return nil, clierr.New("No Value Provided", fmt.Sprintf("key \"%s\" expects string value, but no one provided", arg[2:]))
				}
				field.SetString(value)
			}
			continue
		}

		// short option
		if !optend && rune(arg[0]) == '-' {
			info := expects.GetOpt(arg[1:])

			continue
		}

		// argument
	}

	return nil, clierr.New("Not Implemented", "func Parse[T]([]string) is not implemented")
}

type Arguments struct {
	inner []KV
}

func Normalize(list []string, expects expect.Expects) Arguments {
	array := make([]KV, 0, 0)
}

type KV struct {
	key   string
	value []string
}

func NewKV(key string, value []string) KV {
	return KV{key, value}
}

func (kv KV) Get() (string, []string) {
	return kv.key, kv.value
}
