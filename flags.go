package dsflags

import (
	"flag"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"golang.org/x/exp/constraints"
)

var (
	registeredFlags []Usager
)

var (
	_ flag.Value = &genericValue[int]{}
	_ Usager     = &genericValue[int]{}
)

type genericValue[T Flaggable] struct {
	value       any
	pointer     *T
	short, long string
	usage       string
}

// Flaggable is a constraint that specifies all the supported types for flags.
type Flaggable interface {
	constraints.Integer | constraints.Float | time.Duration | bool | string
}

// convertFromStr tries to convert the string value into the type provided.
// Even though there are no constraints the expected type is any integer or
// floating pointer number that can be natively casted from a float64.
//
// An error is returned if the types cannot be natively casted.
func convertFromStr[To Flaggable](to *To, from string) error {
	// Parse the string to the highest precision possible.
	parsed, err := strconv.ParseFloat(from, 64)
	if err != nil {
		return fmt.Errorf("cannot convert %q to %T: %w", from, to, err)
	}
	parsedVal := reflect.ValueOf(parsed)
	// Since the parameter is a pointer we can use Elem to get the pointed value
	// without having to check if it is a pointer.
	typ := reflect.TypeOf(to).Elem()
	if !parsedVal.CanConvert(typ) {
		return fmt.Errorf("casting %T to %T: %w", parsed, to, ErrCannotCastToType)
	}
	// Convert the float64 to the type this is natively.
	*to = parsedVal.Convert(typ).Interface().(To)

	return nil
}

// Short implements the Usager interface and returns the short alias for this flag.
func (v *genericValue[T]) Short() string {
	return v.short
}

// Long implements the Usager interface and returns the long alias for this flag.
func (v *genericValue[T]) Long() string {
	return v.long
}

// Usage implements the Usager interface and returns the usage for this flag.
func (v *genericValue[T]) Usage() string {
	return v.usage
}

// DefaultValue implements the Usager interface and returns the default value for this flag.
func (v *genericValue[T]) DefaultValue() any {
	return v.value
}

// FlagType implements the Usager interface and returns the type of the flag.
// Booleans are not printed.
func (v *genericValue[T]) FlagType() string {
	typStr := reflect.TypeOf(v.value).String()
	if typStr == "bool" { // omit boolean types
		return ""
	}
	return typStr
}

// IsZeroValue implements the Usager interface and is true if the default value
// is zero value for its type.
func (v *genericValue[T]) IsZeroValue() bool {
	var empty T
	return empty == v.value
}

// Set implements the flag.Value interface and sets the value of the flag's pointer
// from a string.
func (v *genericValue[T]) Set(from string) error {
	switch v.value.(type) {
	case string:
		*(any(v.pointer).(*string)) = from // pointer fuckery
	case uint, uint8, uint16, uint32, int, int8, int16, int32, int64, float32, float64:
		if err := convertFromStr(v.pointer, from); err != nil {
			return err
		}
	case bool:
		b, err := strconv.ParseBool(from)
		if err != nil {
			return err
		}
		*(any(v.pointer).(*bool)) = b
	case time.Duration:
		dur, err := time.ParseDuration(from)
		if err != nil {
			return err
		}
		*(any(v.pointer).(*time.Duration)) = dur
	}

	return nil
}

// String returns the current value of this flag as a string.
func (v *genericValue[T]) String() string {
	if v.pointer == nil {
		var empty T
		return fmt.Sprint(empty)
	}
	return fmt.Sprint(*v.pointer)
}

// IsBoolFlag returns true if this flag is of boolean type.
func (v *genericValue[T]) IsBoolFlag() bool {
	_, ok := v.value.(bool)
	return ok
}

// Flag creates and registers a new flag of the specified type.
// The type is either inferred from the input or specified via the generic syntax.
// After creating the flag it is registered as a commandline flag with a short and long alias.
func Flag[T Flaggable](short byte, long string, value T, usage string) *T {
	newFlag := value
	FlagVar[T](&newFlag, string(short), long, value, usage)
	return &newFlag
}

// LongFlag creates and registers a new flag with only a long alias.
func LongFlag[T Flaggable](long string, value T, usage string) *T {
	newFlag := value
	FlagVar[T](&newFlag, "", long, value, usage)
	return &newFlag
}

// ShortFlag creates and registers a new flag with only a short alias.
func ShortFlag[T Flaggable](short byte, value T, usage string) *T {
	newFlag := value
	FlagVar[T](&newFlag, string(short), "", value, usage)
	return &newFlag
}

// FlagVar creates and registers a new flag with the given pointer.
func FlagVar[T Flaggable](pointer *T, short string, long string, value T, usage string) {
	f := &genericValue[T]{
		value:   any(value),
		pointer: pointer,
		short:   short,
		long:    long,
		usage:   usage,
	}
	// Register both variations of this flag.
	if short != "" {
		flag.CommandLine.Var(f, short, usage)
	}
	if long != "" {
		flag.CommandLine.Var(f, long, usage)
	}
	// Append the flag to the registered flags global var.
	registeredFlags = append(registeredFlags, f)
}
