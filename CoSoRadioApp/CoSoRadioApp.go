package cosoradioapp

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	coso "github.com/jasondborneman/cosoradio/CoSo"
	spotify "github.com/jasondborneman/cosoradio/Spotify"
	"github.com/mpvl/unique"
	spotifyapi "github.com/zmb3/spotify/v2"
)

func Run(client spotifyapi.Client, scrapeCoSo bool, doToot bool) error {
	fmt.Println("     ████                              ")
	fmt.Println("   ██░░░░/█                            ")
	fmt.Println("  ██░░░░/░░░██                          ")
	fmt.Println("  ██░░░/░░░░██                          ")
	fmt.Println("██░#CoSoRadio░██                        ")
	fmt.Println("██░░░░░░░░░░░░██                        ")
	fmt.Println("██░░░░░░░░░░░░██                        ")
	fmt.Println("  ██░░░░░░░░██                          ")
	fmt.Println("    ████████")
	fmt.Println()
	fmt.Println(fmt.Sprintf("----------] Scrape CoSo: %t", scrapeCoSo))
	fmt.Println(fmt.Sprintf("----------] Do Toot: %t", doToot))
	fmt.Println("-----------------------------------------")

	var songs []spotify.Song
	var err error
	if scrapeCoSo {
		fmt.Println("Scraping CoSo for #cosoradio...")
		fmt.Println("-------------------------------")
		songs, err = coso.GetSongsFromCoSo()
		if err != nil {
			return err
		}
	} else {
		fmt.Println("Pulling Music From Fixture...")
		fmt.Println("-----------------------------")
		songs, err = coso.ReadSongsFromFixture()
		if err != nil {
			return err
		}
	}

	if songs == nil || len(songs) == 0 {
		return errors.New("No songs found!")
	}

	t := time.Now()
	dateString := t.Format("01-02-2006")
	var uniqueRecommenders []string
	for _, recommendedBy := range songs {
		uniqueRecommenders = append(uniqueRecommenders, recommendedBy.RecommendedBy)
	}
	unique.Strings(&uniqueRecommenders)
	recommendersString := ""
	moreThanFive := false
	for i, recommendedBy := range uniqueRecommenders {
		if i > 5 {
			moreThanFive = true
			break
		}
		recommendersString = fmt.Sprintf("%s\n%s", recommendersString, recommendedBy)
	}
	if moreThanFive {
		recommendersString = fmt.Sprintf("%s\n%s", recommendersString, "...and more!")
	}
	tootMessage := fmt.Sprintf("#CoSoRadio for %s!\n\nFeaturing music recommendations from: %s\n\n%s",
		dateString,
		recommendersString,
		"XXXX",
	)

	ctx := context.Background()
	playlistUrl, err := spotify.CreatePlaylist(client, ctx, songs, recommendersString)
	if err != nil {
		return err
	}
	if doToot {
		tootMessage = strings.Replace(tootMessage, "XXXX", playlistUrl, 1)
		err = coso.TootSongs(songs, tootMessage)
		if err != nil {
			return err
		}
	} else {
		fmt.Println("Made a playlist, but not tooting to CoSO")
	}
	return nil
}
