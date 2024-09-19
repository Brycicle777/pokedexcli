package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("pokedex >")

	for {
		if err := scanner.Err(); err != nil {
			fmt.Fprintln(os.Stderr, "Error reading command:", err)
			fmt.Println("pokedex >")
			continue
		}
		if scanner.Scan() {
			enteredText := scanner.Text()
			if enteredText == "help" {
				fmt.Println("Use this pokedex with words!")
				fmt.Println("pokedex >")
				continue
			}
			if enteredText == "exit" {
				fmt.Println("Exiting program...")
				break
			}
		}

	}
}
