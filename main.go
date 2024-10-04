package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*Config) error
}

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

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Exits the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:        "map",
			description: "Displays the next 20 locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 locations",
			callback:    commandMapb,
		},
	}
}

func commandHelp(*Config) error {
	commands := getCommands()
	fmt.Println("")
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	for _, cmd := range commands {
		fmt.Printf("%s: %s\n", cmd.name, cmd.description)
	}
	fmt.Println("")
	return nil
}

func commandExit(*Config) error {
	fmt.Println("Exiting program...")
	os.Exit(0)
	return nil
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

func commandMap(cfg *Config) error {
	err := fetchLocationAreas(cfg, false)
	if err != nil {
		return err
	}

	return nil
}

func commandMapb(cfg *Config) error {
	err := fetchLocationAreas(cfg, true)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	cfg := &Config{}
	commands := getCommands()
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("pokedex > ")

		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "Error reading command:", err)
			continue
		}

		if scanner.Scan() {
			input := scanner.Text()
			if cmd, exists := commands[input]; exists {
				err := cmd.callback(cfg)
				if err != nil {
					fmt.Println("Error executing command:", err)
				}
			} else {
				fmt.Println("Unknown command! Type 'help' for a list of commands.")
			}
		}

	}
}
