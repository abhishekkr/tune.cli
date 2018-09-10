package tunecli_tunefindcom

import (
	"fmt"

	golgoquery "github.com/abhishekkr/gol/golgoquery"
	golhttpclient "github.com/abhishekkr/gol/golhttpclient"
)

func IsArtistOnTunefind(artist string) bool {
	artistUrl := TunefindUrlFor("artist", artist)
	return golhttpclient.LinkExists(artistUrl)
}

func (searchFilter TunefindFilter) TunefindArtist(relUrl string) (songs []TunefindSong) {
	golgoquery.ReloadCache = searchFilter.RefreshCache

	fullUrl := fmt.Sprintf("%s%s", tunefindBaseUrl, relUrl)
	goquerySelector := selectorTunefind["artist"]

	songResults := GoqueryTextFrom(fullUrl, goquerySelector)
	songs = make([]TunefindSong, len(songResults))

	songs = searchFilter.SongsResults(songResults, relUrl)
	return
}
