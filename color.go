package dsflags

var colorer Colorer = &DefaultColorer{}

// UseColorer sets the colorer to be used when printing the usage of flags.
func UseColorer(c Colorer) {
	colorer = c
}

// Colorer is used to define the behaviour of the colorer.
type Colorer interface {
	// UsageOf is used to color the initial, "Usage of ./program:" string.
	UsageOf(string) string
	// Flag is used to color the flags themselves. For example, the -v or --verbose.
	Flag(string) string
	// Comma is used for processing the comma that separates a short and long
	// flag alias.
	Comma(string) string
	// Usage is used to color the usage string of a flag.
	Usage(string) string
	// FlagType colors the type of the flag. For booleans this will be an empty string.
	FlagType(string) string
	// DefaultValue is used to color the default value of a flag. This also includes
	// parenthesis.
	DefaultValue(string) string
}

// DefaultColorer is the default colorer that doesn't color anything.
type DefaultColorer struct{}

// UsageOf returns the string as is.
func (c *DefaultColorer) UsageOf(s string) string {
	return s
}

// Flag returns the string as is.
func (c *DefaultColorer) Flag(s string) string {
	return s
}

// Comma returns the string as is.
func (c *DefaultColorer) Comma(s string) string {
	return s
}

// Usage returns the string as is.
func (c *DefaultColorer) Usage(s string) string {
	return s
}

// FlagType returns the string as is.
func (c *DefaultColorer) FlagType(s string) string {
	return s
}

// DefaultValue returns the string as is.
func (c *DefaultColorer) DefaultValue(s string) string {
	return s
}
