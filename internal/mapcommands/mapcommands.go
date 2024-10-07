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

type LocationExploreResponse struct {
	PokemeonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
	} `json:"pokemon_encounters"`
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

	var locationAreaResponse LocationAreaResponse
	entry, found := c.Get(mapUrl)

	if found {
		err := json.Unmarshal(entry, &locationAreaResponse)
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

		err = json.Unmarshal(data, &locationAreaResponse)
		if err != nil {
			return fmt.Errorf("Error decoding response: %v", err)
		}

		c.Add(mapUrl, data)
	}

	for _, area := range locationAreaResponse.Results {
		fmt.Println(area.Name)
	}

	cfg.Next = locationAreaResponse.Next
	cfg.Previous = locationAreaResponse.Previous

	return nil
}

func CommandMap(n string, c *pokecache.Cache, cfg *Config) error {
	err := fetchLocationAreas(c, cfg, false)
	if err != nil {
		return err
	}

	return nil
}

func CommandMapb(n string, c *pokecache.Cache, cfg *Config) error {
	err := fetchLocationAreas(c, cfg, true)
	if err != nil {
		return err
	}

	return nil
}

func CommandExplore(n string, c *pokecache.Cache, cfg *Config) error {
	mapUrl := "https://pokeapi.co/api/v2/location-area/" + n
	fmt.Printf("Exploring %s...\n", n)

	var locationExploreResponse LocationExploreResponse
	entry, found := c.Get(mapUrl)

	if found {
		err := json.Unmarshal(entry, &locationExploreResponse)
		if err != nil {
			return fmt.Errorf("Error decoding cached entry: %v", err)
		}

	} else {
		res, err := http.Get(mapUrl)
		if err != nil {
			return fmt.Errorf("Error calling explore url: %v", err)
		}
		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("Error reading response: %v", err)
		}

		err = json.Unmarshal(data, &locationExploreResponse)
		if err != nil {
			return fmt.Errorf("Error decoding response: %v", err)
		}

		c.Add(mapUrl, data)
	}

	fmt.Println("Found Pokemon:")
	for _, encounter := range locationExploreResponse.PokemeonEncounters {
		fmt.Println(encounter.Pokemon.Name)
	}

	return nil
}
