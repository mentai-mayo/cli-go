package cli

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"github.com/mentai-mayo/cli-go/array"
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
		if unicode.IsLower(rune(rf.Name[0])) {
			continue
		}

		// get field name
		name := rf.Name

		// check expected type
		switch rf.Type.String() {
		case "string", "int", "uint", "bool":
		default:
			return nil, errors.New(fmt.Sprintf("unsupported expected type \"%s\" detected", rf.Type.String()))
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
		expects = append(expects, Expect{name, position, long, short, etype})
	}

	// create struct
	parsed := reflect.New(rt).Interface().(*T)

	// copy command-line arguments
	arguments := array.FromSlice(args)

	fmt.Printf("args: %#v\n", args)
	fmt.Printf("arguments: %#v\n", arguments)

	var remain []string
	{
		// parse command-line arguments
		argsarr := array.New[string](uint(arguments.Len()))
		for {
			elem, ok := arguments.Dequeue()
			fmt.Printf("elem, ok = %#v, %t\n", elem, ok)
			if !ok {
				break
			}
			if *elem == "-" {
				argsarr.Push(*elem)
				continue
			}
			if *elem == "--" {
				argsarr.Push(*elem)
				for {
					elem, ok := arguments.Dequeue()
					if !ok {
						break
					}
					argsarr.Push(*elem)
				}
				break
			}
			for _, expect := range expects {
				if expect.position < 0 {
					continue
				}
				if *elem == fmt.Sprintf("--%s", expect.long) || *elem == fmt.Sprintf("-%s", expect.short) {
					switch expect.etype {
					case "string", "int":
						value, ok := arguments.Dequeue()
						if !ok {
							return nil, NewNoOptionValueSetErr()
						}
						if expect.etype == "string" {
							reflect.ValueOf(parsed).FieldByName(expect.name).SetString(*value)
						} else {
							value, err := strconv.ParseInt(*value, 10, 32)
							if err != nil {
								return nil, err
							}
							reflect.ValueOf(parsed).FieldByName(expect.name).SetInt(value)
						}
					case "bool":
						reflect.ValueOf(parsed).FieldByName(expect.name).SetBool(true)
					}
				}
			}
			argsarr.Push(*elem)
		}
		remain = argsarr.Into()
		fmt.Printf("rem: %#v\n", remain)
	}

	fmt.Printf("expects: %#v\n", expects)

	for _, expect := range expects {
		if expect.position < 0 {
			continue
		}
		if expect.position >= len(remain) {
			return nil, NewPositionOutOfRangeErr(expect.position)
		}
		switch expect.etype {
		case "string":
			reflect.ValueOf(parsed).FieldByName(expect.name).SetString(remain[expect.position])
		case "int":
			num, err := strconv.ParseInt(remain[expect.position], 10, 32)
			if err != nil {
				return nil, err
			}
			reflect.ValueOf(parsed).FieldByName(expect.name).SetInt(num)
		case "bool":
			switch remain[expect.position] {
			case "true":
				reflect.ValueOf(parsed).FieldByName(expect.name).SetBool(true)
			case "false":
				reflect.ValueOf(parsed).FieldByName(expect.name).SetBool(false)
			default:
				return nil, NewParseBoolErr(remain[expect.position])
			}
		}
	}

	return parsed, nil
}

type Expect struct {
	name  string // field name
	long  string
	short string
	etype reflect.Kind // string, int, bool
}

type Expects struct {
	arguments []Expect
	options   []Expect
}

func InitExpects(rtype reflect.Type) (*Expects, error) {
	fields := rtype.NumField()
	arguments := make([]Expect, 0, fields)
	options := make([]Expect, 0, fields)
	position := 0
	for i := 0; i < fields; i++ {
		field := rtype.Field(i)

		name := field.Name

		// --- private field check ---
		if unicode.IsLower(rune(name[0])) {
			// `field` is private field
			continue
		}

		// --- expected type ---
		var etype reflect.Kind
		switch field.Type.Kind() {
		case reflect.String, reflect.Int, reflect.Bool:
			etype = field.Type.Kind()
		default:
			return nil, errors.New(fmt.Sprintf("Unsupported Expected Type: Unsupported expected type [%s (kind: %s)] detected", field.Type.String(), field.Type.Kind().String()))
		}

		// --- get long option name ---
		long := field.Tag.Get("long")

		// --- get short option name ---
		short := field.Tag.Get("short")

		if long == "" && short == "" {
			arguments = append(arguments, Expect{name, long, short, etype})
		} else {
			options = append(options, Expect{name, long, short, etype})
		}
	}
}

type CLIArguments struct {
	cursor int
	slice  []string
}

func (a *CLIArguments) Next() string {
	var next string
	if len(a.slice) <= a.cursor {
		next = ""
	} else {
		next = a.slice[a.cursor]
		a.cursor += 1
	}
	return next
}

func NewCLIArguments(osargs []string) CLIArguments {
	return CLIArguments{cursor: 0, slice: osargs}
}

// ----- errors -----

type ParseBoolErr struct {
	raw string
}

func NewParseBoolErr(raw string) ParseBoolErr {
	return ParseBoolErr{raw}
}

func (e ParseBoolErr) Error() string {
	return fmt.Sprintf("cannot parse \"%s\" as bool", e.raw)
}

type PositionOutOfRangeErr struct {
	pos int
}

func NewPositionOutOfRangeErr(pos int) PositionOutOfRangeErr {
	return PositionOutOfRangeErr{pos}
}

func (e PositionOutOfRangeErr) Error() string {
	return fmt.Sprintf("Position Out of Range Error: Position %d is out of range", e.pos)
}

type NoOptionValueSetErr struct{}

func NewNoOptionValueSetErr() NoOptionValueSetErr {
	return NoOptionValueSetErr{}
}

func (e NoOptionValueSetErr) Error() string {
	return "No Option Value Set Error: no option value set"
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
