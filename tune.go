package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var (
	tunefindBaseUrl = "https://www.tunefind.com"
)

type TunefindSearcResult struct {
	RelUrl string
}

func (searchResult TunefindSearcResult) LinkType() string {
	urlTypeRegex, _ := regexp.Compile("^/([A-Za-z]*)/")
	urlType := urlTypeRegex.FindStringSubmatch(searchResult.RelUrl)[1]
	if urlType == "show" {
		return "tv"
	}
	return urlType
}

func LinkExists(url string) bool {
	var netClient = &http.Client{
		Timeout: time.Second * 10,
	}
	response, err := netClient.Get(url)
	if err != nil || response.StatusCode > 399 {
		return false
	}
	return true
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
	return LinkExists(movieUrl)
}

func IsTvOnTunefind(show string) bool {
	showUrl := TunefindUrlFor("tv", show)
	return LinkExists(showUrl)
}

func IsArtistOnTunefind(artist string) bool {
	artistUrl := TunefindUrlFor("artist", artist)
	return LinkExists(artistUrl)
}

func TunefindSearch(searchQuery string, searchType string) {
	doc, err := goquery.NewDocument(TunefindUrlFor("search", searchQuery))
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("div.row.tf-search-results a").Each(func(i int, s *goquery.Selection) {
		href, hrefPresent := s.Attr("href")
		searchResult := TunefindSearcResult{RelUrl: href}
		if !hrefPresent {
			return
		}
		if searchType == searchResult.LinkType() || searchType == "all" {
			fmt.Println(searchResult.RelUrl)
		}
	})
}

func main() {
	searchQuery := flag.String("search", "", "Item to search. (Required)")
	searchType := flag.String("type", "all", "default:all|movie|tv|artist")
	/*
		seasonIndex := flag.Int("season", 1, "which season if it's a tv type")
		episodeIndex := flag.Int("episode", 1, "which episode if it's a movie type")
	*/

	flag.Parse()

	if *searchQuery == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	TunefindSearch(*searchQuery, *searchType)
}
