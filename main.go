package main

import (
	"bufio"
	"fmt"
	"internal/mapcommands"
	"internal/pokecache"
	"internal/pokecollect"
	"os"
	"strings"
	"time"
)

type cliCommand struct {
	callback    func(string, *pokecache.Cache, *mapcommands.Config) error
	name        string
	description string
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
			callback:    mapcommands.CommandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Displays the previous 20 locations",
			callback:    mapcommands.CommandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Explores the specified area",
			callback:    mapcommands.CommandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Attempts to catch the specified pokemon",
			callback:    pokecollect.CommandCatch,
		},
	}
}

func commandHelp(string, *pokecache.Cache, *mapcommands.Config) error {
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

func commandExit(string, *pokecache.Cache, *mapcommands.Config) error {
	fmt.Println("Exiting program...")
	os.Exit(0)
	return nil
}

func main() {
	cfg := &mapcommands.Config{}
	c := pokecache.NewCache(5 * time.Second)
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
			splitInput := strings.Split(input, " ")
			command := ""
			parameter := ""
			if len(splitInput) > 1 {
				command = splitInput[0]
				parameter = splitInput[1]
			} else {
				command = splitInput[0]
			}
			if cmd, exists := commands[command]; exists {
				err := cmd.callback(parameter, c, cfg)
				if err != nil {
					fmt.Println("Error executing command:", err)
				}
			} else {
				fmt.Println("Unknown command! Type 'help' for a list of commands.")
			}
		}

	}
}
