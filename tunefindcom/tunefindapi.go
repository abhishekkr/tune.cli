package tunecli_tunefindcom

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/abhishekkr/gol/golgoquery"
	"github.com/abhishekkr/gol/golhttpclient"
)

var (
	tunefindBaseUrl = "https://www.tunefind.com"
)

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

func TunefindSearch(searchQuery string, searchType string) {
	fullUrl := TunefindUrlFor("search", searchQuery)
	goquerySelector := "div.row.tf-search-results a"
	var tunefindSearchResults TunefindSearchResults
	tunefindSearchResults.GoqueryResultsToTunefindSearchResults(golgoquery.GoqueryHrefsFrom(fullUrl, goquerySelector))
	for _, result := range tunefindSearchResults.Results {
		if searchType != result.LinkType() && searchType != "all" {
			continue
		}

		fmt.Println("[*]", result.RelUrl)
		if result.LinkType() == "tv" {
			TunefindTv(result.RelUrl)
		} else if result.LinkType() == "movie" {
			TunefindMovie(result.RelUrl)
		} else if result.LinkType() == "artist" {
			TunefindArtist(result.RelUrl)
		}
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
	for _, result := range golgoquery.GoqueryHrefsFromParents(fullUrl, artistSelector).Results {
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

func TunefindTvEpisodeSongs(relUrl string) {
	fullUrl := fmt.Sprintf("%s%s", tunefindBaseUrl, relUrl)

	goquerySelector := "div.Tunefind__Content div.SongRow__center___1I0Cg h4.SongTitle__heading___3kxXK a"
	for _, result := range golgoquery.GoqueryHrefsFrom(fullUrl, goquerySelector).Results {
		song := TunefindSong{RelUrl: result}
		song.TunefindSongsDetails(relUrl)
		fmt.Println(song)
	}
}

func TunefindTvEpisodes(relUrl string) {
	fullUrl := fmt.Sprintf("%s%s", tunefindBaseUrl, relUrl)

	goquerySelector := "div.Tunefind__Content li.MainList__item___fZ13_ h3.EpisodeListItem__title___32XUR a"
	for _, result := range golgoquery.GoqueryHrefsFrom(fullUrl, goquerySelector).Results {
		TunefindTvEpisodeSongs(result)
	}
}

func TunefindTv(relUrl string) {
	fullUrl := fmt.Sprintf("%s%s", tunefindBaseUrl, relUrl)

	goquerySelector := "div.Tunefind__Content ul[aria-labelledby='season-dropdown'] a[role='menuitem']"
	for _, result := range golgoquery.GoqueryHrefsFrom(fullUrl, goquerySelector).Results {
		TunefindTvEpisodes(result)
	}
}

func TunefindMovie(relUrl string) {
	fullUrl := fmt.Sprintf("%s%s", tunefindBaseUrl, relUrl)

	goquerySelector := "div.Tunefind__Content div.SongRow__center___1I0Cg h4.SongTitle__heading___3kxXK a"
	for _, result := range golgoquery.GoqueryHrefsFrom(fullUrl, goquerySelector).Results {
		song := TunefindSong{RelUrl: result}
		song.TunefindSongsDetails(relUrl)
		fmt.Println(result, ">>>", song)
	}
}

func TunefindArtist(relUrl string) {
	fullUrl := fmt.Sprintf("%s%s", tunefindBaseUrl, relUrl)

	goquerySelector := "div.Tunefind__Content div.AppearanceRow__songInfoTitleBlock___3woDL div.AppearanceRow__songInfoTitle___38aKt"
	for _, result := range golgoquery.GoqueryTextFrom(fullUrl, goquerySelector).Results {
		song := TunefindSong{RelUrl: result}
		song.TunefindSongsDetails(relUrl)
		fmt.Println(result, ">>>", song)
	}
}
