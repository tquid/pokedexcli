package pokeapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Direction int

const (
	Forward Direction = iota
	Backward
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

type Client struct {
	config *Config
	apiUrl string
}

func NewClient() *Client {
	client := &Client{
		config: &Config{
			Count:    0,
			Next:     "",
			Previous: "",
			Results:  nil,
		},
		apiUrl: "https://pokeapi.co/api/v2/",
	}
	return client
}

func (c *Client) IsNew() bool {
	if c.config.Results == nil {
		return true
	}
	return false
}

func (c *Client) NextLocationAreas() error {
	var url string
	if c.config.Next != "" {
		url = c.config.Next
	} else {
		url = fmt.Sprintf("%s/location-area", c.apiUrl)
	}
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("can't get %s: %v", url, err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API call to %s failed: %v", url, err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("can't read response body: err")
	}
	defer resp.Body.Close()
	err = json.Unmarshal(body, c.config)
	if err != nil {
		return fmt.Errorf("can't read response body: %v", err)
	}
	return nil
}

func (c *Client) PreviousLocationAreas() error {
	var url string
	if c.config.Previous != "" {
		url = c.config.Previous
	} else {
		return fmt.Errorf("can't go back at beginning of map")
	}
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("can't get %s: %v", url, err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API call to %s failed: %v", url, err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("can't read response body: err")
	}
	defer resp.Body.Close()
	err = json.Unmarshal(body, c.config)
	if err != nil {
		return fmt.Errorf("can't read response body: %v", err)
	}
	return nil
}

func (c *Client) GetLocationNames() []string {
	var names []string
	for _, result := range c.config.Results {
		names = append(names, result.Name)
	}
	return names
}
