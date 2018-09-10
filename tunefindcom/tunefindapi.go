package tunecli_tunefindcom

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	golgoquery "github.com/abhishekkr/gol/golgoquery"
)

var (
	tunefindBaseUrl      = "https://www.tunefind.com"
	tunefindMovieSubUrl  = "/movies/"
	tunefindShowSubUrl   = "/show/"
	tunefindArtistSubUrl = "/artist/"
	tunefindSearchSubUrl = "/search/site?q="
)

type TunefindFilter struct {
	SearchQuery, SearchType, SearchFor   string
	SeasonIndex, EpisodeIndex, SongIndex int
	RefreshCache                         bool
}

type TunefindSong struct {
	SearchTitle string
	Title       string
	RelUrl      string
	Artist      string
	ArtistUrl   string
	YoutubeUrl  string
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
		return fmt.Sprintf("%s%s%s", tunefindBaseUrl, tunefindMovieSubUrl, queryItem)
	} else if urlType == "tv" {
		return fmt.Sprintf("%s%s%s", tunefindBaseUrl, tunefindShowSubUrl, queryItem)
	} else if urlType == "artist" {
		return fmt.Sprintf("%s%s%s", tunefindBaseUrl, tunefindArtistSubUrl, queryItem)
	} else if urlType == "search" {
		return fmt.Sprintf("%s%s%s", tunefindBaseUrl, tunefindSearchSubUrl, queryItem)
	}
	return ""
}

func GoqueryHrefsFromParents(url string, selectors []string) []string {
	golgoquery.CacheGoquery = true
	results, err := golgoquery.GoqueryHrefsFromParents(url, selectors)
	if err != nil {
		log.Println("[error] fetching some attributes for", url)
	}
	return results.Results
}

func GoqueryTextFromParents(url string, selectors []string) []string {
	golgoquery.CacheGoquery = true
	results, err := golgoquery.GoqueryTextFromParents(url, selectors)
	if err != nil {
		log.Println("[error] fetching some attributes for", url)
	}
	return results.Results
}

func GoqueryHrefsFrom(url string, selector string) []string {
	golgoquery.CacheGoquery = true
	results, err := golgoquery.GoqueryHrefsFrom(url, selector)
	if err != nil {
		log.Println("[error] fetching some attributes for", url)
	}
	return results.Results
}

func GoqueryTextFrom(url string, selector string) []string {
	golgoquery.CacheGoquery = true
	results, err := golgoquery.GoqueryTextFrom(url, selector)
	if err != nil {
		log.Println("[error] fetching some attributes for", url)
	}
	return results.Results
}
