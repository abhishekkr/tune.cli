package tunecli_tunefindcom

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/abhishekkr/gol/golgoquery"
	"github.com/abhishekkr/gol/golhttpclient"
)

var (
	tunefindBaseUrl = "https://www.tunefind.com"
)

type TunefindFilter struct {
	SearchQuery, SearchType              string
	SeasonIndex, EpisodeIndex, SongIndex int
}

type TunefindSong struct {
	Title             string
	RelUrl            string
	Artist            string
	ArtistUrl         string
	YoutubeForwardUrl string
}

type TunefindSearcResult struct {
	RelUrl string
	Songs  []TunefindSong
}

type TunefindSearchResults struct {
	Results []TunefindSearcResult
}

func (searchResult TunefindSearcResult) LinkType() string {
	urlTypeRegex, _ := regexp.Compile("^/([A-Za-z]*)/")
	urlType := urlTypeRegex.FindStringSubmatch(searchResult.RelUrl)[1]
	if urlType == "show" {
		return "tv"
	}
	return urlType
}

func TunefindUrlFor(urlType string, queryItem string) string {
	urlType = strings.ToLower(urlType)
	queryItem = strings.ToLower(queryItem)

	if urlType == "movie" {
		return fmt.Sprintf("%s/movies/%s", tunefindBaseUrl, queryItem)
	} else if urlType == "tv" {
		return fmt.Sprintf("%s/show/%s", tunefindBaseUrl, queryItem)
	} else if urlType == "artist" {
		return fmt.Sprintf("%s/artist/%s", tunefindBaseUrl, queryItem)
	} else if urlType == "search" {
		return fmt.Sprintf("%s/search/site?q=%s", tunefindBaseUrl, queryItem)
	}
	return ""
}

func IsMovieOnTunefind(movie string) bool {
	movieUrl := TunefindUrlFor("movie", movie)
	return golhttpclient.LinkExists(movieUrl)
}

func IsTvOnTunefind(show string) bool {
	showUrl := TunefindUrlFor("tv", show)
	return golhttpclient.LinkExists(showUrl)
}

func IsArtistOnTunefind(artist string) bool {
	artistUrl := TunefindUrlFor("artist", artist)
	return golhttpclient.LinkExists(artistUrl)
}

func (results *TunefindSearchResults) GoqueryResultsToTunefindSearchResults(goqueryResults golgoquery.GoqueryResults) {
	results.Results = make([]TunefindSearcResult, len(goqueryResults.Results))
	for idx, goqueryResult := range goqueryResults.Results {
		results.Results[idx] = TunefindSearcResult{RelUrl: goqueryResult}
	}
}

func (song *TunefindSong) TunefindSongsDetailsArtist(fullUrl string) {
	artistSelector := []string{
		fmt.Sprintf("div.Tunefind__Content div.SongRow__center___1I0Cg h4.SongTitle__heading___3kxXK a[href='%s']", song.RelUrl),
		"..",
		"..",
		"..",
		"a.Tunefind__Artist",
	}
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
	for _, result := range golgoquery.GoqueryHrefsFromParents(fullUrl, youtubeUrlSelector).Results {
		song.YoutubeForwardUrl = result
	}
}

func (song *TunefindSong) TunefindSongsDetails(listPageUrl string) {
	fullUrl := fmt.Sprintf("%s%s", tunefindBaseUrl, listPageUrl)

	goquerySelector := fmt.Sprintf("div.Tunefind__Content div.SongRow__center___1I0Cg h4.SongTitle__heading___3kxXK a[href='%s']", song.RelUrl)
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
	}
	return
}

func (searchFilter TunefindFilter) TunefindTvEpisodeSongs(relUrl string) (songs []TunefindSong) {
	fullUrl := fmt.Sprintf("%s%s", tunefindBaseUrl, relUrl)
	goquerySelector := "div.Tunefind__Content div.SongRow__center___1I0Cg h4.SongTitle__heading___3kxXK a"
	songResults := golgoquery.GoqueryHrefsFrom(fullUrl, goquerySelector).Results
	songs = make([]TunefindSong, len(songResults))

	songs = searchFilter.SongsResults(songResults, relUrl)
	return
}

func (searchFilter TunefindFilter) TunefindTvEpisodes(relUrl string) (songs []TunefindSong) {
	fullUrl := fmt.Sprintf("%s%s", tunefindBaseUrl, relUrl)
	goquerySelector := "div.Tunefind__Content li.MainList__item___fZ13_ h3.EpisodeListItem__title___32XUR a"
	episodeResults := golgoquery.GoqueryHrefsFrom(fullUrl, goquerySelector).Results

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
	fullUrl := fmt.Sprintf("%s%s", tunefindBaseUrl, relUrl)
	goquerySelector := "div.Tunefind__Content ul[aria-labelledby='season-dropdown'] a[role='menuitem']"
	seasonResults := golgoquery.GoqueryHrefsFrom(fullUrl, goquerySelector).Results

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

func (searchFilter TunefindFilter) TunefindMovie(relUrl string) (songs []TunefindSong) {
	fullUrl := fmt.Sprintf("%s%s", tunefindBaseUrl, relUrl)
	goquerySelector := "div.Tunefind__Content div.SongRow__center___1I0Cg h4.SongTitle__heading___3kxXK a"
	songResults := golgoquery.GoqueryHrefsFrom(fullUrl, goquerySelector).Results
	songs = make([]TunefindSong, len(songResults))

	songs = searchFilter.SongsResults(songResults, relUrl)
	return
}

func (searchFilter TunefindFilter) TunefindArtist(relUrl string) (songs []TunefindSong) {
	fullUrl := fmt.Sprintf("%s%s", tunefindBaseUrl, relUrl)
	goquerySelector := "div.Tunefind__Content div.AppearanceRow__songInfoTitleBlock___3woDL div.AppearanceRow__songInfoTitle___38aKt"
	songResults := golgoquery.GoqueryTextFrom(fullUrl, goquerySelector).Results
	songs = make([]TunefindSong, len(songResults))

	songs = searchFilter.SongsResults(songResults, relUrl)
	return
}

func (searchFilter TunefindFilter) TunefindSearch() (songs map[string][]TunefindSong) {
	fullUrl := TunefindUrlFor("search", searchFilter.SearchQuery)
	goquerySelector := "div.row.tf-search-results a"
	var tunefindSearchResults TunefindSearchResults
	tunefindSearchResults.GoqueryResultsToTunefindSearchResults(golgoquery.GoqueryHrefsFrom(fullUrl, goquerySelector))
	songs = make(map[string][]TunefindSong, len(tunefindSearchResults.Results))
	for _, result := range tunefindSearchResults.Results {
		if searchFilter.SearchType != result.LinkType() && searchFilter.SearchType != "all" {
			continue
		}

		var _songs []TunefindSong
		if result.LinkType() == "tv" {
			_songs = searchFilter.TunefindTv(result.RelUrl)
		} else if result.LinkType() == "movie" {
			_songs = searchFilter.TunefindMovie(result.RelUrl)
		} else if result.LinkType() == "artist" {
			_songs = searchFilter.TunefindArtist(result.RelUrl)
		}
		songs[result.RelUrl] = _songs
	}
	return
}

func ShowTunefindSongs(songsMap map[string][]TunefindSong) {
	for relUrl, songs := range songsMap {
		fmt.Printf("[ %s ]\n", relUrl)
		for _, song := range songs {
			fmt.Printf("[*] %s\n", song.Title)
			fmt.Printf("    [url](%s%s)\n", tunefindBaseUrl, song.RelUrl)
			fmt.Printf("    by [%s](%s%s)\n", song.Artist, tunefindBaseUrl, song.ArtistUrl)
			fmt.Printf("    listen at [youtube](%s%s)\n", tunefindBaseUrl, song.YoutubeForwardUrl)
		}
	}
}
