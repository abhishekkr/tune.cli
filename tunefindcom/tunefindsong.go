package tunecli_tunefindcom

import (
	"fmt"
	"log"

	golhttpclient "github.com/abhishekkr/gol/golhttpclient"

	youtubecom "github.com/abhishekkr/tune.cli/youtubecom"
)

func (song TunefindSong) FirstYoutubeLink() string {
	return youtubecom.FirstLink(song.YoutubeUrl)
}

func (song *TunefindSong) TunefindSongsDetailsArtist(fullUrl string) {
	for _, result := range GoqueryTextFromParents(fullUrl, selectorSongsDetailsArtist(song.RelUrl)) {
		song.Artist = result
	}
}

func (song *TunefindSong) TunefindSongsDetailsArtistLink(fullUrl string) {
	for _, result := range GoqueryHrefsFromParents(fullUrl, selectorSongsDetailsArtistLink(song.RelUrl)) {
		song.ArtistUrl = result
	}
}

func (song *TunefindSong) TunefindSongsDetailsYoutube(fullUrl string) {
	for _, result := range GoqueryHrefsFromParents(fullUrl, selectorSongsDetailsYoutube(song.RelUrl)) {
		song.YoutubeUrl = golhttpclient.UrlRedirectTo(fmt.Sprintf("%s%s", tunefindBaseUrl, result))
	}
}

func (song *TunefindSong) TunefindSongsDetails(listPageUrl string) {
	fullUrl := fmt.Sprintf("%s%s", tunefindBaseUrl, listPageUrl)

	for _, result := range GoqueryTextFrom(fullUrl, selectorSongDetails(song.RelUrl)) {
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
