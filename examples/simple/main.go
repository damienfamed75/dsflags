package main

import (
	"fmt"
	"time"

	"github.com/damienfamed75/dsflags"
)

func main() {
	name := dsflags.Flag('n', "name", "John Doe", "Name to greet")
	age := dsflags.Flag('a', "age", 25, "Age of the person")
	delay := dsflags.LongFlag("delay", time.Second*0, "Delay before printing")
	verbose := dsflags.ShortFlag('v', false, "Verbose output")

	dsflags.Parse()

	if *delay > 0 {
		if *verbose {
			fmt.Printf("Waiting for %s\n", *delay)
		}
		time.Sleep(*delay)
	}

	fmt.Printf("Hello, %s! You are %d years old.\n", *name, *age)
}
