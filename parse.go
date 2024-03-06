package dsflags

import "flag"

// Parse goes through and read all the registered flags from the command line.
func Parse() {
	flag.Parse()
}
