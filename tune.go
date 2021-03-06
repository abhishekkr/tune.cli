package main

import (
	"flag"
	"log"
	"os"

	tunefindcom "github.com/abhishekkr/tune.cli/tunefindcom"
)

func main() {
	searchFor := flag.String("out", "list", "default:list|play|json")
	searchFrom := flag.String("src", "tunefind", "default:tunefind")
	searchQuery := flag.String("search", "", "Item to search. (Required)")
	searchType := flag.String("type", "all", "default:all|movie|tv|artist")
	seasonIndex := flag.Int("season", 0, "which season if it's a tv type, default:0 for all")
	episodeIndex := flag.Int("episode", 0, "which episode if it's a tv type, default:0 for all")
	songIndex := flag.Int("song", 0, "which song, default:0 for all")
	refreshCache := flag.Bool("refresh", false, "reloading cache, dfault:false")

	flag.Parse()

	if *searchQuery == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *searchFrom == "tunefind" {
		tunefindFilter := tunefindcom.TunefindFilter{
			SearchQuery:  *searchQuery,
			SearchType:   *searchType,
			SearchFor:    *searchFor,
			SeasonIndex:  (*seasonIndex - 1),
			EpisodeIndex: (*episodeIndex - 1),
			SongIndex:    (*songIndex - 1),
			RefreshCache: *refreshCache,
		}

		_ = tunefindFilter.TunefindSearch() //persist if flag passed to a playlist
	} else {
		log.Fatalf("%s source isn't supported yet, try tunefind maybe.", searchFrom)
	}
}
