package tunecli_youtubecom

import (
	"fmt"
	"log"

	"github.com/abhishekkr/gol/golgoquery"
)

var (
	youtubeBaseUrl = "https://www.youtube.com"
)

func Links(url string) []string {
	goquerySelector := "div#page-container div.branded-page-v2-body.branded-page-v2-primary-column-content ol li h3.yt-lockup-title a"
	results, err := golgoquery.GoqueryHrefsFrom(url, goquerySelector)
	if err != nil {
		log.Println("[warn] error go-querying", url)
	}
	for idx, result := range results.Results {
		results.Results[idx] = fmt.Sprintf("%s%s", youtubeBaseUrl, result)
	}
	return results.Results
}

func FirstLink(url string) string {
	links := Links(url)
	if len(links) == 0 {
		log.Println("[error] no youtube url found for", url)
		return ""
	}
	return links[0]
}
