package main

import (
	"flag"
	"os"

	tunefindcom "github.com/abhishekkr/tune.cli/tunefindcom"
)

func main() {
	searchQuery := flag.String("search", "", "Item to search. (Required)")
	searchType := flag.String("type", "all", "default:all|movie|tv|artist")
	/*
		seasonIndex := flag.Int("season", 1, "which season if it's a tv type")
		episodeIndex := flag.Int("episode", 1, "which episode if it's a movie type")
	*/

	flag.Parse()

	if *searchQuery == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	tunefindcom.TunefindSearch(*searchQuery, *searchType)
}
