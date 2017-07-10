package tunecli_tunefindcom

import (
	"fmt"
	"log"

	golgoquery "github.com/abhishekkr/gol/golgoquery"
	golhttpclient "github.com/abhishekkr/gol/golhttpclient"

	youtubecom "github.com/abhishekkr/tune.cli/youtubecom"
)

func (song TunefindSong) FirstYoutubeLink() string {
	return youtubecom.FirstLink(song.YoutubeUrl)
}

func (song *TunefindSong) TunefindSongsDetailsArtist(fullUrl string) {
	artistSelector := []string{
		fmt.Sprintf("div.Tunefind__Content div.SongRow__center___1I0Cg h4.SongTitle__heading___3kxXK a[href='%s']", song.RelUrl),
		"..",
		"..",
		"..",
		"a.Tunefind__Artist",
	}
	golgoquery.CacheGoquery = true
	for _, result := range golgoquery.GoqueryTextFromParents(fullUrl, artistSelector).Results {
		song.Artist = result
	}
}

func (song *TunefindSong) TunefindSongsDetailsArtistLink(fullUrl string) {
	artistUrlSelector := []string{
		fmt.Sprintf("div.Tunefind__Content div.SongRow__center___1I0Cg h4.SongTitle__heading___3kxXK a[href='%s']", song.RelUrl),
		"..",
		"..",
		"..",
		"a.Tunefind__Artist",
	}
	golgoquery.CacheGoquery = true
	for _, result := range golgoquery.GoqueryHrefsFromParents(fullUrl, artistUrlSelector).Results {
		song.ArtistUrl = result
	}
}

func (song *TunefindSong) TunefindSongsDetailsYoutube(fullUrl string) {
	youtubeUrlSelector := []string{
		fmt.Sprintf("div.Tunefind__Content div.SongRow__center___1I0Cg h4.SongTitle__heading___3kxXK a[href='%s']", song.RelUrl),
		"..",
		"..",
		"..",
		"..",
		"a.StoreLinks__youtube___2MHoI",
	}
	golgoquery.CacheGoquery = true
	for _, result := range golgoquery.GoqueryHrefsFromParents(fullUrl, youtubeUrlSelector).Results {
		song.YoutubeUrl = golhttpclient.UrlRedirectTo(fmt.Sprintf("%s%s", tunefindBaseUrl, result))
	}
}

func (song *TunefindSong) TunefindSongsDetails(listPageUrl string) {
	fullUrl := fmt.Sprintf("%s%s", tunefindBaseUrl, listPageUrl)

	goquerySelector := fmt.Sprintf("div.Tunefind__Content div.SongRow__center___1I0Cg h4.SongTitle__heading___3kxXK a[href='%s']", song.RelUrl)
	golgoquery.CacheGoquery = true
	for _, result := range golgoquery.GoqueryTextFrom(fullUrl, goquerySelector).Results {
		song.Title = result
	}

	song.TunefindSongsDetailsArtist(fullUrl)
	song.TunefindSongsDetailsArtistLink(fullUrl)
	song.TunefindSongsDetailsYoutube(fullUrl)
}

func (searchFilter TunefindFilter) SongsResults(songResults []string, relUrl string) (songs []TunefindSong) {
	if searchFilter.SongIndex > len(songResults) {
		log.Printf("[warn] song#%d not found, it only has %d songs",
			(searchFilter.SongIndex + 1),
			(len(songResults) + 1),
		)
		return
	} else if searchFilter.SongIndex >= 0 {
		songs = make([]TunefindSong, 1)
		songs[0] = TunefindSong{RelUrl: songResults[searchFilter.SongIndex]}
		songs[0].TunefindSongsDetails(relUrl)
		return
	}

	songs = make([]TunefindSong, len(songResults))
	for idx, result := range songResults {
		songs[idx] = TunefindSong{RelUrl: result}
		songs[idx].TunefindSongsDetails(relUrl)
		searchFilter.TunefindSongOutput(songs[idx])
	}
	return
}
