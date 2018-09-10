package tunecli_tunefindcom

import "fmt"

var (
	selectorTunefind = map[string]string{
		"artist": "div.Tunefind__Content div.AppearanceRow__songInfoTitleBlock___3woDL div.AppearanceRow__songInfoTitle___38aKt",
		"movie":  "div.Tunefind__Content div.SongRow__center___1I0Cg h4.SongTitle__heading___3kxXK a",
		"search": "div.row.tf-search-results a",
		"show":   "div.Tunefind__Content ul[aria-labelledby='season-dropdown'] a[role='menuitem']",
	}

	selectorShow = map[string]string{
		"episodes":     "div.Tunefind__Content li.MainList__item___2MKl8 h3.EpisodeListItem__title___32XUR a",
		"episodeSongs": "div.Tunefind__Content div.SongRow__container___3eT_L h4.SongTitle__heading___3kxXK a",
	}
)

func selectorSongDetails(relUrl string) string {
	return fmt.Sprintf("div.Tunefind__Content div.SongRow__center___1HKjk h4.SongTitle__heading___3kxXK a[href='%s']", relUrl)
}

func selectorSongsDetailsYoutube(relUrl string) []string {
	return []string{
		fmt.Sprintf("div.Tunefind__Content div.SongRow__center___1HKjk h4.SongTitle__heading___3kxXK a[href='%s']", relUrl),
		"..",
		"..",
		"..",
		"..",
		"a.StoreLinks__youtube___2meaC",
	}
}

func selectorSongsDetailsArtist(relUrl string) []string {
	return []string{
		fmt.Sprintf("div.Tunefind__Content div.SongRow__center___1HKjk h4.SongTitle__heading___3kxXK a[href='%s']", relUrl),
		"..",
		"..",
		"..",
		"a.Subtitle__subtitle___1rSyh",
	}
}

func selectorSongsDetailsArtistLink(relUrl string) []string {
	return []string{
		fmt.Sprintf("div.Tunefind__Content div.SongRow__center___1HKjk h4.SongTitle__heading___3kxXK a[href='%s']", relUrl),
		"..",
		"..",
		"..",
		"a.Subtitle__subtitle___1rSyh",
	}
}
