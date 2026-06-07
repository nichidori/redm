package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/nichidori/redm/internal/config"
	"github.com/nichidori/redm/internal/redmineapi"
)

type state struct {
	client *redmineapi.Client
	config *config.Config
}

type Command struct {
	Name        string
	Description string
	Setup       func(fs *flag.FlagSet, s *state) func() error
}

type CLI struct {
	Name         string
	State        *state
	Commands     map[string]Command
	CommandNames []string
}

func (c *CLI) Register(cmd Command) {
	if c.Commands == nil {
		c.Commands = make(map[string]Command)
	}

	if _, exists := c.Commands[cmd.Name]; !exists {
		c.CommandNames = append(c.CommandNames, cmd.Name)
	}

	c.Commands[cmd.Name] = cmd
}

func (c *CLI) Execute(args []string) error {
	if len(args) < 1 || args[0] == "help" {
		c.PrintGlobalHelp()
		return nil
	}

	cmdName := args[0]
	cmd, exists := c.Commands[cmdName]
	if !exists {
		return fmt.Errorf("unknown command: %s", cmdName)
	}

	fs := flag.NewFlagSet(cmd.Name, flag.ExitOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s %s [options]\n\n", c.Name, cmd.Name)
		fmt.Fprintf(os.Stderr, "%s\n\nOptions:\n", cmd.Description)
		fs.PrintDefaults()
	}

	run := cmd.Setup(fs, c.State)

	if err := fs.Parse(args[1:]); err != nil {
		return err
	}

	return run()
}

func (c *CLI) PrintGlobalHelp() {
	fmt.Printf("Usage: %s <command> [options]\n", c.Name)
	fmt.Println("\nAvailable commands:")

	for _, name := range c.CommandNames {
		cmd := c.Commands[name]
		fmt.Printf("  %s %s\n", FixLength(cmd.Name, 12), cmd.Description)
	}
}

func SelectOption[T any](reader *bufio.Reader, label string, options []T, optionFormatter func(T) string) (T, error) {
	fmt.Printf("Select %s:\n", label)
	for i, o := range options {
		fmt.Printf("(%v) %s\n", i+1, optionFormatter(o))
	}

	fmt.Printf("Enter %s number: ", label)

	input, err := reader.ReadString('\n')
	if err != nil {
		var t T
		return t, err
	}

	idx, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil {
		var t T
		return t, err
	}
	idx--

	if idx < 0 || idx >= len(options) {
		var t T
		return t, fmt.Errorf("invalid choice: must be between 1 and %d", len(options))
	}

	return options[idx], nil
}
