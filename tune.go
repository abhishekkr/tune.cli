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

type TunefindSearcResults struct {
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

func GoqueryDocument(url string) *goquery.Document {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}
	return doc
}

func GoqueryAttrsFrom(url string, goquerySelector string, attr string) (results TunefindSearcResults) {
	doc := GoqueryDocument(url)

	a_nodes := doc.Find(goquerySelector)
	results.Results = make([]TunefindSearcResult, a_nodes.Size())
	a_nodes.Each(func(i int, s *goquery.Selection) {
		var attrValue string
		var attrPresent bool
		if attr == "text" {
			attrValue = s.Text()
			attrPresent = (attrValue != "")
		} else {
			attrValue, attrPresent = s.Attr(attr)
		}
		if !attrPresent {
			s_html, _ := s.Html()
			log.Printf("[warn] %s\n", s_html)
		}

		results.Results[i] = TunefindSearcResult{RelUrl: attrValue}
	})
	return
}

func GoqueryHrefsFrom(url string, goquerySelector string) (results TunefindSearcResults) {
	return GoqueryAttrsFrom(url, goquerySelector, "href")
}

func GoqueryTextFrom(url string, goquerySelector string) (results TunefindSearcResults) {
	return GoqueryAttrsFrom(url, goquerySelector, "text")
}

func TunefindSearch(searchQuery string, searchType string) {
	fullUrl := TunefindUrlFor("search", searchQuery)
	goquerySelector := "div.row.tf-search-results a"
	for _, result := range GoqueryHrefsFrom(fullUrl, goquerySelector).Results {
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

func TunefindMovie(relUrl string) {
	fullUrl := fmt.Sprintf("%s%s", tunefindBaseUrl, relUrl)

	goquerySelector := "div.Tunefind__Content div.SongRow__center___1I0Cg h4.SongTitle__heading___3kxXK a"
	for _, result := range GoqueryHrefsFrom(fullUrl, goquerySelector).Results {
		fmt.Println(result.RelUrl)
	}
}

func TunefindTv(relUrl string) {
	fullUrl := fmt.Sprintf("%s%s", tunefindBaseUrl, relUrl)

	goquerySelector := "div.Tunefind__Content ul[aria-labelledby='season-dropdown'] a[role='menuitem']"
	for _, result := range GoqueryHrefsFrom(fullUrl, goquerySelector).Results {
		fmt.Println(result.RelUrl)
	}
}

func TunefindArtist(relUrl string) {
	fullUrl := fmt.Sprintf("%s%s", tunefindBaseUrl, relUrl)

	goquerySelector := "div.Tunefind__Content div.AppearanceRow__songInfoTitleBlock___3woDL div.AppearanceRow__songInfoTitle___38aKt"
	for _, result := range GoqueryTextFrom(fullUrl, goquerySelector).Results {
		fmt.Println(result.RelUrl)
	}
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
