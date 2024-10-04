package mapcommands

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type LocationAreaResponse struct {
	Results []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
}

type Config struct {
	Next     *string
	Previous *string
}

func fetchLocationAreas(cfg *Config, usePrevious bool) error {
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

	res, err := http.Get(mapUrl)
	if err != nil {
		return fmt.Errorf("Error calling map: %v", err)
	}
	defer res.Body.Close()

	var locationAreaReponse LocationAreaResponse
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&locationAreaReponse)
	if err != nil {
		return fmt.Errorf("Error decoding response: %v", err)
	}

	for _, area := range locationAreaReponse.Results {
		fmt.Println(area.Name)
	}

	cfg.Next = locationAreaReponse.Next
	cfg.Previous = locationAreaReponse.Previous

	return nil
}

func CommandMap(cfg *Config) error {
	err := fetchLocationAreas(cfg, false)
	if err != nil {
		return err
	}

	return nil
}

func CommandMapb(cfg *Config) error {
	err := fetchLocationAreas(cfg, true)
	if err != nil {
		return err
	}

	return nil
}
