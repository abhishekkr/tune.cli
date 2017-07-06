package tunecli_youtubecom

import (
	"fmt"

	"github.com/abhishekkr/gol/golgoquery"
)

var (
	youtubeBaseUrl = "https://www.youtube.com"
)

func Links(url string) []string {
	goquerySelector := "div#page-container div.branded-page-v2-body.branded-page-v2-primary-column-content ol li h3.yt-lockup-title a"
	results := golgoquery.GoqueryHrefsFrom(url, goquerySelector).Results
	for idx, result := range results {
		results[idx] = fmt.Sprintf("%s%s", youtubeBaseUrl, result)
	}
	return results
}

func FirstLink(url string) string {
	return Links(url)[0]
}
