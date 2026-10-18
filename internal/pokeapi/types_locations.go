package pokeapi

// RespShallowLocations -
type RespShallowLocations struct {
	Count    int     `json:"count"`
	Next     *string `json:"next"`
	Previous *string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type CurrentLocation struct {
	ID                   int                    `json:"id"`
	Name                 string                 `json:"name"`
	GameIndex            int                    `json:"game_index"`
	Location             Location               `json:"location"`
	PokemonEncounters    []PokemonEncounters    `json:"pokemon_encounters"`
	
}
type Location struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type Pokemon struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
type PokemonEncounters struct {
	Pokemon        Pokemon          `json:"pokemon"`
}

type CurrentPokemon struct {
	ID                     int              `json:"id"`
	Name                   string           `json:"name"`
	BaseExperience         int              `json:"base_experience"`
	Height				   int				`json:"height"`
	Weight				   int				`json:"weight"`
	Stats  				   []Stats		    `json:"stats"`
	Types				   []Types			`json:"types"`
	Sprites                struct {
		FrontDefault string `json:"front_default"`
	} `json:"sprites`
}

type Stats struct {
	BaseStat int  `json:"base_stat"`
	Effort   int  `json:"effort"`
	Stats	 Stat `json:"stat"`
}

type Stat struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Types struct {
	Slot int   `json:"slot"`
	Type Type  `json:"type"`
}

type Type struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}


