package mapcommands

import (
	"encoding/json"
	"errors"
	"fmt"
	"internal/pokecache"
	"io"
	"net/http"
)

type LocationAreaResponse struct {
	Previous *string `json:"previous"`
	Next     *string `json:"next"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type Config struct {
	Next     *string
	Previous *string
}

func fetchLocationAreas(c *pokecache.Cache, cfg *Config, usePrevious bool) error {
	mapUrl := ""
	if usePrevious {
		if cfg.Previous != nil {
			mapUrl = *cfg.Previous
		} else {
			return errors.New("You're already on the first page!")
		}
	} else {
		if cfg.Next != nil {
			mapUrl = *cfg.Next
		} else {
			mapUrl = "https://pokeapi.co/api/v2/location-area/"
		}
	}

	var locationAreaReponse LocationAreaResponse
	entry, found := c.Get(mapUrl)

	if found {
		err := json.Unmarshal(entry, &locationAreaReponse)
		if err != nil {
			return fmt.Errorf("Error decoding cached entry: %v", err)
		}

	} else {
		res, err := http.Get(mapUrl)
		if err != nil {
			return fmt.Errorf("Error calling map url: %v", err)
		}
		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("Error reading response: %v", err)
		}

		err = json.Unmarshal(data, &locationAreaReponse)
		if err != nil {
			return fmt.Errorf("Error decoding response: %v", err)
		}

		c.Add(mapUrl, data)
	}

	for _, area := range locationAreaReponse.Results {
		fmt.Println(area.Name)
	}

	cfg.Next = locationAreaReponse.Next
	cfg.Previous = locationAreaReponse.Previous

	return nil
}

func CommandMap(c *pokecache.Cache, cfg *Config) error {
	err := fetchLocationAreas(c, cfg, false)
	if err != nil {
		return err
	}

	return nil
}

func CommandMapb(c *pokecache.Cache, cfg *Config) error {
	err := fetchLocationAreas(c, cfg, true)
	if err != nil {
		return err
	}

	return nil
}
