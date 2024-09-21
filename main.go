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

func initCommands() map[string]cliCommand {
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
			callback:    commandMap,
		},
		// "mapb": {
		// 	name:        "mapb",
		// 	description: "show previous 20 map entries",
		// 	callback:    commandMapb,
		// },
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

func commandMap() error {
	cfg := pokeapi.Config{
		Count:    0,
		Next:     "",
		Previous: "",
		Results:  []pokeapi.LocationArea,
	}
	err := cfg.NextLocationAreas()
	if err != nil {
		fmt.Printf("Error getting next map chunk: %v\n", err)
	}
	for _, result := range cfg.Results {
		fmt.Println(result.Name)
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
