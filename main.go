package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/tquid/pokedexcli/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func initCommands(client *pokeapi.Client) map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "show next 20 map entries",
			callback:    func() error { return commandMap(client) },
		},
		"mapb": {
			name:        "mapb",
			description: "show previous 20 map entries",
			callback:    func() error { return commandMapb(client) },
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

func commandMap(c *pokeapi.Client) error {
	err := c.NextLocationAreas()
	if err != nil {
		fmt.Printf("Error getting next map chunk: %v\n", err)
	}
	for _, name := range c.GetLocationNames() {
		fmt.Println(name)
	}
	return nil
}

func commandMapb(c *pokeapi.Client) error {
	err := c.PreviousLocationAreas()
	if err != nil {
		fmt.Printf("Error getting previous map chunk: %v\n", err)
	}
	for _, name := range c.GetLocationNames() {
		fmt.Println(name)
	}
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
	c := pokeapi.NewClient()
	cmds := initCommands(c)

	for {
		command, err := promptAndRead()
		if err != nil {
			fmt.Printf("unable to read input: %v\n", err)
		}
		if _, exists := cmds[command]; exists {
			cmds[command].callback()
		} else {
			fmt.Printf("unknown command '%s'\n", command)
		}
	}
}
