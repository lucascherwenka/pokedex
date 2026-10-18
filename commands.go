package main

import (
	"fmt" 
	"os"
	"errors"
	"math/rand"
	"strings"
	"regexp"
)

const (
	colorReset   = "\033[0m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorBlue    = "\033[34m"
	colorCyan    = "\033[36m"
	colorMagenta = "\033[35m"
	colorYellow  = "\033[33m"
)

func commandExit(cfg *config, _ []string) error {
	fmt.Printf("\n    %sClosing the Pokedex... Goodbye!%s\n\n", colorYellow, colorReset)
	os.Exit(0)
	return nil
}

func commandHelp(cfg *config, _ []string) error {
	fmt.Println()
	fmt.Printf("	%s--Welcome to the Pokedex!--%s\n\n", colorYellow, colorReset)
	fmt.Printf("  Usage:\n")
	fmt.Println()
	for _, cmd := range cfg.commands {
		fmt.Printf("  - %s: 	%s\n", cmd.name, cmd.description)
	}
	fmt.Println()
	return nil
}

func commandMapf(cfg *config, _ []string) error {
	locationsResp, err := cfg.pokeapiClient.ListLocations(cfg.nextLocationsURL)
	if err != nil {
		return err
	}

	cfg.nextLocationsURL = locationsResp.Next
	cfg.prevLocationsURL = locationsResp.Previous

	fmt.Printf("\n %sMapping locations!%s\n\n", colorYellow, colorReset)
	for _, loc := range locationsResp.Results {
		fmt.Printf("   %s\n", loc.Name)
	}
	fmt.Printf("\n %sTo check a location for pokemon use %s<explore><area name>%s\n\n", colorYellow, colorBlue, colorReset)
	return nil
}

func commandMapb(cfg *config, _ []string) error {
	if cfg.prevLocationsURL == nil {
		return errors.New("\nyou're on the first page\n")
	}

	locationResp, err := cfg.pokeapiClient.ListLocations(cfg.prevLocationsURL)
	if err != nil {
		return err
	}

	cfg.nextLocationsURL = locationResp.Next
	cfg.prevLocationsURL = locationResp.Previous

	fmt.Printf("\n %sMapping locations!%s\n\n", colorYellow, colorReset)
	for _, loc := range locationResp.Results {
		fmt.Printf("   %s\n", loc.Name)
	}
	fmt.Printf("\n %sTo check a location for pokemon use %s<explore><area name>%s\n\n", colorYellow, colorBlue, colorReset)
	return nil
}

func commandExplore(cfg *config, areaName []string) error {
	if len(areaName) == 0 {
		fmt.Printf("\n%sNo area selected!%s\n", colorRed, colorReset)
		return nil
	}
	locationResp, err := cfg.pokeapiClient.GetLocation(areaName[0])
	if err != nil {
		fmt.Printf("\n%s<%s>%s is not a valid location%s\n", colorRed, areaName[0], colorYellow, colorReset)
		return nil
	}
	fmt.Printf("\n  %sExploring %s%s...%s\n", colorMagenta, colorYellow, areaName[0], colorReset)
	fmt.Printf("\n  Found Pokemon:\n")
	for i, _ := range locationResp.PokemonEncounters {
		fmt.Printf("    - %s\n", locationResp.PokemonEncounters[i].Pokemon.Name)
	}
	return nil
}

func commandCatch(cfg *config, pokemon []string) error {
	if len(pokemon) == 0 {
		fmt.Printf("%sNo Pokemon selected!%s\n", colorRed, colorReset)
		return nil
	}
	pokemonResp, err := cfg.pokeapiClient.GetPokemon(pokemon[0])
	if err != nil {
		fmt.Printf("<%s> is not a valid Pokemon\n", pokemon[0])
		return nil
	}
	catchChance := 1
	if 100 - (pokemonResp.BaseExperience / 3) <= 0 {
		catchChance = 1
	} else {
		catchChance = 100 - (pokemonResp.BaseExperience / 3)
	}
	fmt.Printf("\n   %sThrowing a Pokeball at %s...%s\n", colorMagenta, pokemon[0], colorReset)
	if rand.Intn(100) < catchChance {
		fmt.Printf("   %s%s was caught!%s\n", colorGreen, pokemon[0], colorReset)
		fmt.Printf("   %sYou may now inspect it with the %s<inspect>%s command!%s\n\n", colorGreen, colorYellow, colorGreen, colorReset)
		cfg.caughtPokemon[pokemon[0]] = pokemonResp
	} else {
		fmt.Printf("   %s%s escaped!%s\n\n",colorRed, pokemon[0], colorReset)
	}
	return nil
}

func commandInspect(cfg *config, pokemon []string) error {
	if len(pokemon) != 1 {
		fmt.Println("\n   You must provide a pokemon name\n")
		return nil
	}

	mon, ok := cfg.caughtPokemon[pokemon[0]]
	if !ok {
		fmt.Println("\n   You have not caught that pokemon yet.\n")
		return nil
	}

	const spriteWidth = 64
	const spriteHeight = 64
	const leftColWidth = 30

	// 1. Get sprite rendered into lines (32x32 pixels = 16 terminal lines)
	spriteLines := cfg.pokeapiClient.GetSpriteLines(mon.Sprites.FrontDefault, spriteWidth, spriteHeight)

	for len(spriteLines) > 0 && strings.TrimSpace(stripANSI(spriteLines[0])) == "" {
    	spriteLines = spriteLines[1:]
	}
	// 2. Prepare the right-hand text lines
	textLines := []string{
		fmt.Sprintf("%sName: %s%s", colorYellow, mon.Name, colorReset),
		fmt.Sprintf("Height: %d", mon.Height),
		fmt.Sprintf("Weight: %d", mon.Weight),
		"Stats:",
	}
	for _, stat := range mon.Stats {
		textLines = append(textLines, fmt.Sprintf("  -%s: %v", stat.Stats.Name, stat.BaseStat))
	}
	textLines = append(textLines, "Types:")
	for _, tizype := range mon.Types {
		textLines = append(textLines, fmt.Sprintf("  -%s", tizype.Type.Name))
	}

	// 3. Print side-by-side
	maxLines := len(spriteLines)
	if len(textLines) > maxLines {
		maxLines = len(textLines)
	}

	fmt.Println()
	for i := 0; i < maxLines; i++ {
		text := ""
		if i < len(textLines) {
			text = textLines[i]
		}

		sprite := ""
		if i < len(spriteLines) {
			sprite = spriteLines[i]
		}

		// Calculate visible padding (simple padding check)
		padding := leftColWidth - len(stripANSI(text))
		if padding < 2 {
			padding = 2
		}

		spaces := strings.Repeat(" ", padding)
		fmt.Printf("   %s%s%s\n", text, spaces, sprite)
	}
	fmt.Println()

	return nil
}

func commandPokedex(cfg *config, _ []string) error {
	if len(cfg.caughtPokemon) <= 0 {
		fmt.Println("\n	Your pokedex is empty! Go catch some 'mons!\n")
	} else {
		fmt.Println(" \n Your Pokedex:")
		for mon, _ := range cfg.caughtPokemon {
			fmt.Printf("   - %v\n", mon)
		}
	}
	return nil
}

var ansiRegex = regexp.MustCompile(`\x1b\[[0-9;]*m`)

func stripANSI(str string) string {
	return ansiRegex.ReplaceAllString(str, "")
}