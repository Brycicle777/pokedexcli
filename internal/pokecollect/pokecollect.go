package pokecollect

import (
	"encoding/json"
	"fmt"
	"internal/mapcommands"
	"internal/pokecache"
	"io"
	"math/rand"
	"net/http"
)

type Pokemon struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
}

var Pokedex = make(map[string]Pokemon)

func CommandCatch(n string, c *pokecache.Cache, cfg *mapcommands.Config) error {
	catchUrl := "https://pokeapi.co/api/v2/pokemon/" + n

	var pokemonToCatch Pokemon
	entry, found := c.Get(catchUrl)

	if found {
		err := json.Unmarshal(entry, &pokemonToCatch)
		if err != nil {
			return fmt.Errorf("Error decoding cached entry: %v", err)
		}

	} else {
		res, err := http.Get(catchUrl)
		if err != nil {
			return fmt.Errorf("Error calling map url: %v", err)
		}
		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("Error reading response: %v", err)
		}

		err = json.Unmarshal(data, &pokemonToCatch)
		if err != nil {
			return fmt.Errorf("Error decoding response (please check your spelling): %v", err)
		}

		c.Add(catchUrl, data)
	}

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonToCatch.Name)
	catchChance := rand.Intn(700 - pokemonToCatch.BaseExperience)
	if catchChance >= 100 {
		fmt.Printf("%s was caught!\n", pokemonToCatch.Name)
		Pokedex[pokemonToCatch.Name] = pokemonToCatch
	} else {
		fmt.Printf("%s escaped!\n", pokemonToCatch.Name)
	}

	return nil
}
