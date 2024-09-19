package main

import (
	"bufio"
	"fmt"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
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
	}

}

func commandHelp() error {
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

func commandExit() error {
	fmt.Println("Exiting program...")
	os.Exit(0)
	return nil
}

func main() {
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
				err := cmd.callback()
				if err != nil {
					fmt.Println("Error executing command:", err)
				}
			} else {
				fmt.Println("Unknown command! Type 'help' for a list of commands.")
			}
		}

	}
}
