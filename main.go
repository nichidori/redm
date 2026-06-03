package main

import (
	"fmt"
	"os"

	"github.com/nichidori/redm/internal/rest"
)

func main() {
	c := rest.NewClient(
		"http://localhost:3000",
		"e7315305dc0f377b1f5d0a37ac353e47436a5b2a",
	)

	cli := &CLI{
		Name: "redm",
		State: &state{
			client: c,
		},
	}

	if err := cli.Execute(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
