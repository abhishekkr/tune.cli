package tunecli_tunefindcom

import (
	"fmt"

	golgoquery "github.com/abhishekkr/gol/golgoquery"
	golhttpclient "github.com/abhishekkr/gol/golhttpclient"
)

func IsMovieOnTunefind(movie string) bool {
	movieUrl := TunefindUrlFor("movie", movie)
	return golhttpclient.LinkExists(movieUrl)
}

func (searchFilter TunefindFilter) TunefindMovie(relUrl string) (songs []TunefindSong) {
	golgoquery.ReloadCache = searchFilter.RefreshCache

	fullUrl := fmt.Sprintf("%s%s", tunefindBaseUrl, relUrl)
	goquerySelector := "div.Tunefind__Content div.SongRow__center___1I0Cg h4.SongTitle__heading___3kxXK a"
	golgoquery.CacheGoquery = true
	songResults := golgoquery.GoqueryHrefsFrom(fullUrl, goquerySelector).Results
	songs = make([]TunefindSong, len(songResults))

	songs = searchFilter.SongsResults(songResults, relUrl)
	return
}
