package tunecli_tunefindcom

import (
	"strings"

	golgoquery "github.com/abhishekkr/gol/golgoquery"
)

func (results *TunefindSearchResults) GoqueryResultsToTunefindSearchResults(searchResults []string) {
	results.Results = make([]TunefindSearcResult, len(searchResults))
	for idx, searchResult := range searchResults {
		results.Results[idx] = TunefindSearcResult{RelUrl: searchResult}
	}
}

func (searchFilter TunefindFilter) TunefindSearch() (songs map[string][]TunefindSong) {
	var tunefindSearchResults TunefindSearchResults
	goquerySelector := "div.row.tf-search-results a"

	golgoquery.ReloadCache = searchFilter.RefreshCache

	searchFilter.SearchQuery = strings.Replace(searchFilter.SearchQuery, " ", "+", -1)
	fullUrl := TunefindUrlFor("search", searchFilter.SearchQuery)

	searchResults := GoqueryHrefsFrom(fullUrl, goquerySelector)
	tunefindSearchResults.GoqueryResultsToTunefindSearchResults(searchResults)

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
