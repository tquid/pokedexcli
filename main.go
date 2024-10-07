package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/tquid/pokedexcli/internal/pokeapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func([]string) error
}

func initCommands(client *pokeapi.Client) map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"explore": {
			name:        "explore",
			description: "Explore an area (use 'explore <area>')",
			callback:    func(params []string) error { return commandExplore(client, params) },
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "show next 20 map entries",
			callback:    func([]string) error { return commandMap(client) },
		},
		"mapb": {
			name:        "mapb",
			description: "show previous 20 map entries",
			callback:    func([]string) error { return commandMapb(client) },
		},
	}
}

func commandHelp([]string) error {
	fmt.Println("help: print this helpful message\nexit: exit the pokedex")
	return nil
}

func commandExit([]string) error {
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

func commandExplore(c *pokeapi.Client, params []string) error {
	var areaName string
	if len(params) == 0 {
		return fmt.Errorf("'explore' command requires an area name, e.g. 'explore canalave-city-area'")
	}
	areaName = params[0]
	pokemonList, err := c.ExploreArea(areaName)
	if err != nil {
		return fmt.Errorf("exploring area %s: %v\n", areaName, err)
	}
	fmt.Printf("Exploring %s...\n", areaName)
	if len(pokemonList) == 0 {
		fmt.Println("No pokemon found!")
		return nil
	}
	fmt.Println("Found pokemon:")
	for _, pokemon := range pokemonList {
		fmt.Printf(" - %s\n", pokemon)
	}
	return nil
}

func promptAndRead() ([]string, error) {
	fmt.Print("pokedex > ")
	s := bufio.NewScanner(os.Stdin)
	s.Scan()
	err := s.Err()
	if err != nil {
		return nil, fmt.Errorf("error trying to scan input: %w", err)
	}
	return strings.Fields(s.Text()), nil
}

func main() {
	c := pokeapi.NewClient()
	cmds := initCommands(c)

	for {
		fields, err := promptAndRead()
		var params []string
		if len(fields) == 0 {
			continue
		}
		command := fields[0]
		if len(fields) > 1 {
			params = fields[1:]
		}
		if err != nil {
			fmt.Printf("Command error: %v\n", err)
		}
		if _, exists := cmds[command]; exists {
			err = cmds[command].callback(params)
			if err != nil {
				fmt.Printf("Error trying command: %v\n", err)
			}
		} else {
			fmt.Printf("unknown command '%s'\n", command)
		}
	}
}
