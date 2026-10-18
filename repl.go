package main

import (
	"strings"
	"fmt"
	"bufio"
	"os"
	"github.com/lucascherwenka/pokedex/internal/pokeapi"
)

type cliCommand struct {
	name		string
	description string
	callback	func(*config, []string) error
}

type config struct {
	commands  	     map[string]cliCommand
	pokeapiClient    pokeapi.Client
	nextLocationsURL *string
	prevLocationsURL *string
	caughtPokemon    map[string]pokeapi.CurrentPokemon
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
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"map": {
			name:		 "map",
			description: "Displays 20 locations",
			callback:	 commandMapf,
		},
		"mapb": {
			name:		 "mapb",
			description: "Displays previous 20 locations",
			callback:	 commandMapb,
		},
		"explore": {
			name: 		 "explore",
			description: "Explore a named area to find the pokemon within",
			callback: 	 commandExplore,
		},
		"catch": {
			name: 		 "catch",
			description: "Attempt to catch a pokemon of your choice within the given area",
			callback: 	 commandCatch,
		},
		"inspect": {
			name:		 "inspect",
			description: "Inspect a pokemon you have caught",
			callback:    commandInspect,
		},
		"pokedex": {
			name: 		 "pokedex",
			description: "List all of the pokemon you have caught",
			callback:    commandPokedex,
		},
	}
}


func startRepl(cfg *config) {
	reader := bufio.NewScanner(os.Stdin)
	cfg.caughtPokemon = make(map[string]pokeapi.CurrentPokemon)
	for {
		fmt.Printf("%s\nPokedex > %s", colorCyan, colorReset)
		reader.Scan()

		words := cleanInput(reader.Text())
		if len(words) == 0 {
			continue
		}

		commandName := words[0]

		command, exists := cfg.commands[commandName]
		if exists {
			err := command.callback(cfg, words[1:])
			if err != nil {
				fmt.Println(err)
			}
			continue
		} else {
			fmt.Println("Unknown command")
			continue
		}
	}
}

func cleanInput(text string) []string {
	output := strings.ToLower(text)
	words := strings.Fields(output)
	return words
}

type pokeStruct struct {
	Count    int       `json:"count"`
	Next     *string   `json:"next"`
	Previous *string   `json:"previous"`
	Results  []Results `json:"results"`
}
type Results struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}



