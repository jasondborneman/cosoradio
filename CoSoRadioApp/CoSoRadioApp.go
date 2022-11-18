package cosoradioapp

import (
	"context"
	"errors"
	"fmt"

	coso "github.com/jasondborneman/cosoradio/CoSo"
	spotify "github.com/jasondborneman/cosoradio/Spotify"
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
	ctx := context.Background()
	playlistUrl, err := spotify.CreatePlaylist(client, ctx, songs)
	if err != nil {
		return err
	}
	if doToot {
		err = coso.TootSongs(playlistUrl, songs)
		if err != nil {
			return err
		}
	} else {
		fmt.Println("Made a playlist, but not tooting to CoSO")
	}
	return nil
}
