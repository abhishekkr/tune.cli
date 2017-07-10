package tunecli_tunefindcom

import (
	"fmt"
	"log"
	"strings"

	golbin "github.com/abhishekkr/gol/golbin"
)

func ShowSongs(songsMap map[string][]TunefindSong) {
	for relUrl, songs := range songsMap {
		fmt.Printf("[ %s ]\n", relUrl)
		for _, song := range songs {
			fmt.Printf("## %s\n", song.Title)
			fmt.Printf("    [url](%s%s)\n", tunefindBaseUrl, song.RelUrl)
			fmt.Printf("    by [%s](%s%s)\n", song.Artist, tunefindBaseUrl, song.ArtistUrl)
			fmt.Printf("    listen at [youtube](%s)\n", song.FirstYoutubeLink())
		}
	}
}

func PlayOrNot() bool {
	var choice string
	fmt.Printf(" | play (y|n): ")
	_, err := fmt.Scanf("%s", &choice)
	if err != nil {
		return false
	}
	choice = strings.ToLower(choice)
	if choice == "y" || choice == "yes" {
		return true
	}
	return false
}

func PlaySongs(songsMap map[string][]TunefindSong) {
	for relUrl, songs := range songsMap {
		fmt.Printf("[ %s ]\n", relUrl)
		for _, song := range songs {
			fmt.Printf("* %s", song.Title)
			fmt.Printf("[by %s]", song.Artist)
			if !PlayOrNot() {
				continue
			}

			cmdOutput := golbin.RunWithAssignedApp(song.FirstYoutubeLink())
			log.Println(cmdOutput)
		}
	}
}
