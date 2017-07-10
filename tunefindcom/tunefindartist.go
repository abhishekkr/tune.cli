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
	goquerySelector := "div.Tunefind__Content div.AppearanceRow__songInfoTitleBlock___3woDL div.AppearanceRow__songInfoTitle___38aKt"
	golgoquery.CacheGoquery = true
	songResults := golgoquery.GoqueryTextFrom(fullUrl, goquerySelector).Results
	songs = make([]TunefindSong, len(songResults))

	songs = searchFilter.SongsResults(songResults, relUrl)
	return
}
