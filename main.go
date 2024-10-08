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
		"catch": {
			name:        "catch",
			description: "Try to catch a Pokemon",
			callback:    func(params []string) error { return commandCatch(client, params) },
		},
		"inspect": {
			name:        "inspect",
			description: "Show Pokemon details",
			callback:    func(params []string) error { return commandInspect(client, params) },
		},
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
		"pokedex": {
			name:        "pokedex",
			description: "show your pokedex",
			callback:    func([]string) error { return commandPokedex(client) },
		},
	}
}

func commandCatch(c *pokeapi.Client, params []string) error {
	if len(params) == 0 {
		return fmt.Errorf("'catch' command requires a pokemon name, e.g. 'catch pikachu'")
	}
	pokemonName := params[0]
	pokemon, err := c.GetPokemon(pokemonName)
	if err != nil {
		return fmt.Errorf("error getting pokemon info: %w", err)
	}
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)
	if pokemon.Catch() {
		fmt.Printf("%s was caught!\n", pokemonName)
		fmt.Println("You may now inspect it with the inspect command.")
		c.AddPokedexEntry(pokemon)
	} else {
		fmt.Printf("%s escaped!\n", pokemonName)
	}
	return nil
}

func commandInspect(c *pokeapi.Client, params []string) error {
	if len(params) == 0 {
		return fmt.Errorf("'inspect' command requires a pokemon name, e.g. 'inspect pikachu'")
	}
	pokemonName := params[0]
	pokemon, ok := c.GetPokedexEntry(pokemonName)
	if !ok {
		fmt.Println("you have not caught that pokemon (or it doesn't exist)")
		return nil
	}
	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("  -%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, pokemonType := range pokemon.Types {
		fmt.Printf("  - %s\n", pokemonType.Type.Name)
	}
	return nil
}

func commandPokedex(c *pokeapi.Client) error {
	fmt.Println("Your Pokedex:")
	pokedex := c.ListPokedex()
	if len(pokedex) == 0 {
		fmt.Println(" Nothing yet!")
		return nil
	}
	for _, name := range pokedex {
		fmt.Printf(" - %s\n", name)
	}
	return nil
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
