// Copyright 2018 The Periph Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// extra-smoketest runs all known smoke tests.
package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sort"

	"github.com/zrnsm/extra/hostextra/d2xx/d2xxsmoketest"
	"periph.io/x/periph/host"
)

// SmokeTest must be implemented by a smoke test. It will be run by this
// executable.
type SmokeTest interface {
	// Name is the name of the smoke test, it is the identifier used on the
	// command line.
	Name() string
	// Description returns a short description to be printed to the user in the
	// help page, to explain what this test does and any requirement to make it
	// work.
	Description() string
	// Run runs the test and return an error in case of failure.
	Run(f *flag.FlagSet, args []string) error
}

// tests is the list of registered smoke tests.
var tests = []SmokeTest{
	&d2xxsmoketest.SmokeTest{},
}

func usage(fs *flag.FlagSet) {
	_, _ = io.WriteString(os.Stderr, "Usage: extra-smoketest <args> <name> ...\n\n")
	fs.PrintDefaults()
	_, _ = io.WriteString(os.Stderr, "\nTests available:\n")
	names := make([]string, len(tests))
	desc := make(map[string]string, len(tests))
	l := 0
	for i := range tests {
		n := tests[i].Name()
		if len(n) > l {
			l = len(n)
		}
		names[i] = n
		desc[n] = tests[i].Description()
	}
	sort.Strings(names)
	for _, name := range names {
		fmt.Fprintf(os.Stderr, "  %-*s %s\n", l, name, desc[name])
	}
}

func mainImpl() error {
	state, err := host.Init()
	if err != nil {
		return fmt.Errorf("error loading drivers: %v", err)
	}
	fs := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	verbose := fs.Bool("v", false, "verbose mode")
	fs.Usage = func() { usage(fs) }
	if err = fs.Parse(os.Args[1:]); err == flag.ErrHelp {
		return nil
	} else if err != nil {
		return err
	}
	if fs.NArg() == 0 {
		fs.Usage()
		_, _ = io.WriteString(os.Stdout, "\n")
		return errors.New("please specify a test to run or use -help")
	}
	cmd := fs.Arg(0)
	if cmd == "help" {
		usage(fs)
		return nil
	}

	if !*verbose {
		log.SetOutput(ioutil.Discard)
	}
	log.SetFlags(log.Lmicroseconds)

	if *verbose {
		if len(state.Failed) > 0 {
			log.Print("Failed to load some drivers:")
			for _, failure := range state.Failed {
				log.Printf("- %s: %v", failure.D, failure.Err)
			}
		}
		log.Printf("Using drivers:")
		for _, driver := range state.Loaded {
			log.Printf("- %s", driver)
		}
		if len(state.Skipped) > 0 {
			log.Printf("Drivers skipped:")
			for _, failure := range state.Skipped {
				log.Printf("- %s: %v", failure.D, failure.Err)
			}
		}
	}

	for _, t := range tests {
		if t.Name() == cmd {
			f := flag.NewFlagSet("extra-smoketest "+t.Name(), flag.ExitOnError)
			u := f.Usage
			f.Usage = func() {
				fmt.Printf("%s: %s\n\n", t.Name(), t.Description())
				u()
				flags := false
				f.VisitAll(func(*flag.Flag) { flags = true })
				if !flags {
					fmt.Printf("  This smoke test doesn't have any flag.\n")
				}
			}
			if err = t.Run(f, fs.Args()[1:]); err == nil {
				log.Printf("Test %s successful", cmd)
			}
			return err
		}
	}
	return fmt.Errorf("test case %q was not found", cmd)
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "extra-smoketest: %s.\n", err)
		os.Exit(1)
	}
}
