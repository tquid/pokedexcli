package main

import (
	"bufio"
	"fmt"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func initCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
	}
}

func commandHelp() error {
	fmt.Println("help: print this helpful message\nexit: exit the pokedex")
	return nil
}

func commandExit() error {
	os.Exit(0)
	return nil
}

func promptAndRead() (string, error) {
	fmt.Print("pokedex > ")
	s := bufio.NewScanner(os.Stdin)
	s.Scan()
	err := s.Err()
	if err != nil {
		return "", fmt.Errorf("error trying to scan input: %v", err)
	}
	return s.Text(), nil
}

func main() {
	c := initCommands()
	for {
		command, err := promptAndRead()
		if err != nil {
			fmt.Printf("unable to read input: %v\n", err)
		}
		if _, exists := c[command]; exists {
			c[command].callback()
		} else {
			fmt.Printf("unknown command '%s'\n", command)
		}
	}
}
