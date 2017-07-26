package tunecli_tunefindcom

import (
	"fmt"
	"log"

	golgoquery "github.com/abhishekkr/gol/golgoquery"
	golhttpclient "github.com/abhishekkr/gol/golhttpclient"
)

func IsTvOnTunefind(show string) bool {
	showUrl := TunefindUrlFor("tv", show)
	return golhttpclient.LinkExists(showUrl)
}

func (searchFilter TunefindFilter) TunefindTvEpisodeSongs(relUrl string) (songs []TunefindSong) {
	fullUrl := fmt.Sprintf("%s%s", tunefindBaseUrl, relUrl)
	goquerySelector := "div.Tunefind__Content div.SongRow__center___1I0Cg h4.SongTitle__heading___3kxXK a"

	songResults := GoqueryHrefsFrom(fullUrl, goquerySelector)
	songs = make([]TunefindSong, len(songResults))

	songs = searchFilter.SongsResults(songResults, relUrl)
	return
}

func (searchFilter TunefindFilter) TunefindTvEpisodes(relUrl string) (songs []TunefindSong) {
	fullUrl := fmt.Sprintf("%s%s", tunefindBaseUrl, relUrl)
	goquerySelector := "div.Tunefind__Content li.MainList__item___fZ13_ h3.EpisodeListItem__title___32XUR a"

	episodeResults := GoqueryHrefsFrom(fullUrl, goquerySelector)

	if searchFilter.EpisodeIndex > len(episodeResults) {
		log.Printf("[warn] episode#%d not found, it only has %d episodes",
			(searchFilter.EpisodeIndex + 1),
			(len(episodeResults) + 1),
		)
		return
	} else if searchFilter.EpisodeIndex >= 0 {
		songs = searchFilter.TunefindTvEpisodeSongs(episodeResults[searchFilter.EpisodeIndex])
		return
	}

	for _, result := range episodeResults {
		songs = searchFilter.TunefindTvEpisodeSongs(result)
	}
	return
}

func (searchFilter TunefindFilter) TunefindTv(relUrl string) (songs []TunefindSong) {
	golgoquery.ReloadCache = searchFilter.RefreshCache

	fullUrl := fmt.Sprintf("%s%s", tunefindBaseUrl, relUrl)
	goquerySelector := "div.Tunefind__Content ul[aria-labelledby='season-dropdown'] a[role='menuitem']"

	seasonResults := GoqueryHrefsFrom(fullUrl, goquerySelector)

	if searchFilter.SeasonIndex > len(seasonResults) {
		log.Printf("[warn] season#%d not found, it only has %d seasons",
			(searchFilter.SeasonIndex + 1),
			(len(seasonResults) + 1),
		)
		return
	} else if searchFilter.SeasonIndex >= 0 {
		songs = searchFilter.TunefindTvEpisodes(seasonResults[searchFilter.SeasonIndex])
		return
	}

	for _, result := range seasonResults {
		songs = searchFilter.TunefindTvEpisodes(result)
	}
	return
}
