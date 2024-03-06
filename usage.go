package dsflags

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// Usager is used to define the behavior needed to print out the usage of flags.
type Usager interface {
	Short() string
	Long() string
	Usage() string
	DefaultValue() any
	IsZeroValue() bool
	FlagType() string
}

func init() {
	flag.Usage = flagUsage
}

// A custom implementation of the flag.Usage function to print each flag as
// long and short aliases instead of as its own flag.
var flagUsage = func() {
	fmt.Fprintf(flag.CommandLine.Output(), colorer.UsageOf("Usage of %s:\n"), os.Args[0])

	for _, f := range registeredFlags {
		var b strings.Builder
		if f.Short() != "" {
			fmt.Fprintf(&b, "  %s", colorer.Flag("-"+f.Short()))
			// If there's a verbose version of this flag then print a comma.
			if f.Long() != "" {
				b.WriteString(colorer.Comma(","))
			}
		} else {
			b.WriteByte(' ')
		}

		if f.Long() != "" {
			fmt.Fprintf(&b, " %s", colorer.Flag("--"+f.Long()))
		}

		flagType := f.FlagType()
		if len(flagType) > 0 { // Don't print boolean flag types
			b.WriteByte(' ')
			b.WriteString(colorer.FlagType(flagType))
		}
		// Since character boolean flags can print usage on the same line
		if b.Len() <= 4 {
			b.WriteByte('\t')
		} else {
			// Four spaces before the tab triggers good alighment for both
			// 4- and 8-space tab stops.
			b.WriteString("\n    \t")
		}
		b.WriteString(colorer.Usage(strings.ReplaceAll(f.Usage(), "\n", "\n    \t")))

		// Only display the default value if it's not the zero value.
		if !f.IsZeroValue() {
			b.WriteByte(' ')
			if _, ok := f.DefaultValue().(string); ok {
				fmt.Fprintf(&b, colorer.DefaultValue("(default %q)"), f.DefaultValue())
			} else {
				fmt.Fprintf(&b, colorer.DefaultValue("(default %v)"), f.DefaultValue())
			}
		}
		// fmt.Fprintf(flag.CommandLine.Output(), b.String(), spaceBetweenFlags)
		fmt.Fprint(flag.CommandLine.Output(), b.String(), spaceBetweenFlags)
	}
}

var spaceBetweenFlags = "\n"

// UseSpacedOutUsage adds an extra newline character between each flag in the
// usage. This can be useful for when there are a lot of flags.
func UseSpacedOutUsage() {
	spaceBetweenFlags = "\n\n"
}

// Usage prints out the usage of all the registered flags.
func Usage() {
	flag.Usage()
}
