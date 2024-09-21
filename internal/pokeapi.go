package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type LocationArea struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Config struct {
	Count    int            `json:"count"`
	Next     string         `json:"next"`
	Previous string         `json:"previous"`
	Results  []LocationArea `json:"results"`
}

const apiUrl = "https://pokeapi.co/api/v2/"

func (cfg *Config) NextLocationAreas() error {
	url := fmt.Sprintf("%s/location-area", apiUrl)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("can't get %s: %v", url, err)
	}
	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API call to %s failed: %v", url, err)
	}

	err = json.Unmarshal(body, cfg)
	if err != nil {
		return fmt.Errorf("can't read response body: %v", err)
	}
	return nil
}
