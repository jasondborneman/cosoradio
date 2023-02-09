package cosoradioapp

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	coso "github.com/jasondborneman/cosoradio/CoSo"
	spotify "github.com/jasondborneman/cosoradio/Spotify"
	spotifyapi "github.com/zmb3/spotify/v2"
	yt "google.golang.org/api/youtube/v3"
)

func Run(spotifyClient spotifyapi.Client, googleService yt.Service, cosoToken string, scrapeCoSo bool, doToot bool) error {
	var songs []spotify.Song
	var err error
	if scrapeCoSo {
		fmt.Println("Scraping CoSo for #cosoradio...")
		fmt.Println("-------------------------------")
		songs, err = coso.GetSongsFromCoSo(googleService, cosoToken)
		if err != nil {
			return err
		}
	} else {
		fmt.Println("Pulling Music From Fixture...")
		fmt.Println("-----------------------------")
		songs, err = coso.ReadSongsFromFixture(googleService)
		if err != nil {
			return err
		}
	}

	if songs == nil || len(songs) == 0 {
		return errors.New("No songs found!")
	}

	t := time.Now()
	dateString := t.Format("01-02-2006")

	urs := make(map[string]bool)
	for _, r := range songs {
		urs[r.RecommendedBy] = true
	}

	recommendersString := ""
	moreThanFive := false
	i := 0
	for recommendedBy, _ := range urs {
		if i > 5 {
			moreThanFive = true
			break
		}
		recommendersString = fmt.Sprintf("%s, %s", recommendersString, "@"+recommendedBy)
		i++
	}

	if moreThanFive {
		recommendersString = fmt.Sprintf("%s %s", recommendersString, "...and more!")
	}
	tootMessage := fmt.Sprintf("#CoSoRadio for %s! Featuring music recommendations from: %s. %s",
		dateString,
		recommendersString,
		"XXXX",
	)

	ctx := context.Background()
	playlistUrl, err := spotify.CreatePlaylist(spotifyClient, ctx, songs, recommendersString)
	if err != nil {
		return err
	}
	if doToot {
		tootMessage = strings.Replace(tootMessage, "XXXX", playlistUrl, 1)
		err = coso.TootSongs(songs, tootMessage, cosoToken)
		if err != nil {
			return err
		}
	} else {
		fmt.Println("Made a playlist, but not tooting to CoSO")
	}
	return nil
}
