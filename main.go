package cli

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func Parse[T any](args []string) (*T, error) {
	rt := reflect.TypeOf((*T)(nil)).Elem()
	rv := reflect.New(rt).Elem()

	// check is T struct
	if rv.Kind().String() != "struct" {
		return nil, NewNonStructTargetErr(rv.Kind().String())
	}

	// get expect names/types
	expects := make([]Expect, 0, rt.NumField())
	for i := 0; i < rt.NumField(); i++ {
		rf := rt.Field(i)

		// check is private
		if rf.PkgPath == "" {
			continue
		}

		// check expected type
		switch rf.Type.String() {
		case "string", "int", "uint", "bool":
		default:
			return nil, errors.New("unsupported expected type detected")
		}
		etype := rf.Type.String()

		// check position
		tag, ok := rf.Tag.Lookup("pos")
		var position int
		if ok {
			num, err := strconv.Atoi(tag)
			if err != nil {
				return nil, errors.New("tag(\"pos\") must be parseable as a integer")
			}
			position = num
		} else {
			position = -1
		}

		// check long option name
		tag, ok = rf.Tag.Lookup("long")
		var long string
		if ok {
			long = tag
		} else {
			long = strings.ToLower(rf.Name)
		}

		// check short option name
		tag, ok = rf.Tag.Lookup("short")
		var short string
		if ok {
			short = tag
		} else {
			short = ""
		}

		// add
		expects = append(expects, Expect{position, long, short, etype})
	}

	return nil, errors.New("Unimplemented")
}

type Expect struct {
	position int
	long     string
	short    string
	etype    string // string, int, uint, bool
}

type NonStructTargetErr struct {
	actual string
}

func NewNonStructTargetErr(actual string) NonStructTargetErr {
	return NonStructTargetErr{actual}
}

func (e NonStructTargetErr) Error() string {
	return fmt.Sprintf("Non-Struct Target Error: arguments cannot be parsed for non-struct types (got: %s)", e.actual)
}
