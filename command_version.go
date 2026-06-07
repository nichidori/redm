package main

import (
	"flag"
	"fmt"
)

var CommandVersion = Command{
	Name:        "version",
	Description: "Show program version",
	Setup: func(fs *flag.FlagSet, s *state) func() error {
		return func() error {
			fmt.Println(Version)
			return nil
		}
	},
}
