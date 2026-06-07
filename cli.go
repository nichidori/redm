package main

import (
	"flag"
	"fmt"
	"os"

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
