package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/tquid/pokedexcli/internal/pokecache"
)

type Direction int

const (
	Forward Direction = iota
	Backward
)

type LocationAreaPage struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type LocationArea struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

type Config struct {
	Count    int                `json:"count"`
	Next     string             `json:"next"`
	Previous string             `json:"previous"`
	Results  []LocationAreaPage `json:"results"`
}

type Client struct {
	config *Config
	apiUrl string
	cache  *pokecache.Cache
}

func NewClient() *Client {
	client := &Client{
		config: &Config{
			Count:    0,
			Next:     "",
			Previous: "",
			Results:  nil,
		},
		apiUrl: "https://pokeapi.co/api/v2",
		cache:  pokecache.NewCache(time.Minute * 5),
	}
	return client
}

func (c *Client) IsNew() bool {
	if c.config.Results == nil {
		return true
	}
	return false
}

func (c *Client) callAPI(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("can't get %s: %w", url, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("call to %s failed: %s", url, resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("can't read response body: err")
	}
	defer resp.Body.Close()
	return body, nil
}

func (c *Client) NextLocationAreas() error {
	var url string
	if c.config.Next != "" {
		url = c.config.Next
	} else {
		url = fmt.Sprintf("%s/location-area", c.apiUrl)
	}
	if body, hit := c.cache.Get(url); hit {
		err := json.Unmarshal(body, c.config)
		if err != nil {
			return fmt.Errorf("can't unmarshal cache result: %w", err)
		}
		return nil
	}
	data, err := c.callAPI(url)
	err = json.Unmarshal(data, c.config)
	if err != nil {
		return fmt.Errorf("can't read response body: %w", err)
	}
	c.cache.Add(url, data)
	return nil
}

func (c *Client) PreviousLocationAreas() error {
	var url string
	if c.config.Previous != "" {
		url = c.config.Previous
	} else {
		return fmt.Errorf("can't go back at beginning of map")
	}
	if data, hit := c.cache.Get(url); hit {
		err := json.Unmarshal(data, c.config)
		if err != nil {
			return fmt.Errorf("can't unmarshal cache result: %w", err)
		}
		return nil
	}
	data, err := c.callAPI(url)
	if err != nil {
		return fmt.Errorf("API error: %w", err)
	}
	err = json.Unmarshal(data, c.config)
	if err != nil {
		return fmt.Errorf("can't unmarshal response body: %w", err)
	}
	c.cache.Add(url, data)
	return nil
}

func PokemonListFromLocationArea(data []byte) ([]string, error) {
	var location LocationArea
	err := json.Unmarshal(data, &location)
	if err != nil {
		return nil, fmt.Errorf("can't unmarshal LocationArea data: %w", err)
	}
	var pokemonList []string
	for _, encounter := range location.PokemonEncounters {
		pokemonList = append(pokemonList, encounter.Pokemon.Name)
	}
	return pokemonList, nil
}

func (c *Client) ExploreArea(areaName string) ([]string, error) {
	url := fmt.Sprintf("%s/location-area/%s", c.apiUrl, areaName)
	if body, hit := c.cache.Get(url); hit {
		pokemonList, err := PokemonListFromLocationArea(body)
		if err != nil {
			return nil, fmt.Errorf("cache read error: %w", err)
		}
		return pokemonList, nil
	}
	body, err := c.callAPI(url)
	if err != nil {
		return nil, fmt.Errorf("API error: %w", err)
	}
	pokemonList, err := PokemonListFromLocationArea(body)
	if err != nil {
		return nil, fmt.Errorf("can't read location area data: %w", err)
	}
	c.cache.Add(url, body)
	return pokemonList, nil
}

func (c *Client) GetLocationNames() []string {
	var names []string
	for _, result := range c.config.Results {
		names = append(names, result.Name)
	}
	return names
}
