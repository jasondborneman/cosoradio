package spotify

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	tools "github.com/jasondborneman/cosoradio/Tools"
	spotifyapi "github.com/zmb3/spotify/v2"
)

func CreatePlaylist(client spotifyapi.Client, ctx context.Context, songs []Song) (string, error) {
	var playlistUrl string
	user, err := client.CurrentUser(ctx)
	if err != nil {
		log.Printf("Error getting current user: %v", err)
		return "", err
	}
	t := time.Now()
	dateString := t.Format("01-02-2006")
	playlistName := fmt.Sprintf("CoSoRadio [Testing] %s", dateString)
	playlistDescription := fmt.Sprintf("The counter.social #CoSoRadio playlist for %s", dateString)
	fullPlaylist, err := client.CreatePlaylistForUser(ctx, user.ID, playlistName, playlistDescription, true, false)
	if err != nil {
		log.Printf("Error creating playlist %v", err)
		return "", err
	}
	playlistUrl = string(fullPlaylist.ExternalURLs["spotify"])

	png, err := os.Open("Images/cosoradio.png")
	if err != nil {
		log.Printf("Error opening thumbnail for playlist: %v", err)
		return "", err
	}
	pngWithLabel, err := tools.AddDateToThumbnail(png, dateString, 1, 13)
	err = client.SetPlaylistImage(ctx, fullPlaylist.ID, pngWithLabel)
	if err != nil {
		log.Printf("Error adding thumbnail to playlist: %v", err)
		return "", err
	}

	for i, song := range songs {
		log.Printf("%d : %s, %s, %s", i, song.Artist, song.Name, song.YouTubeTitle)
		var searchQuery string
		if song.YouTubeTitle == "" {
			searchQuery = fmt.Sprintf("%s %s", song.Artist, song.Name)
		} else {
			searchQuery = fmt.Sprintf("%s", song.YouTubeTitle)
		}
		searchResult, err := client.Search(ctx, searchQuery, spotifyapi.SearchTypeTrack)
		if err != nil {
			log.Printf("error searching Spotify: %v", err)
			return "", err
		}
		topResult := searchResult.Tracks.Tracks[0]
		client.AddTracksToPlaylist(ctx, fullPlaylist.ID, topResult.ID)
	}

	log.Printf("-------------------------------------------------")
	log.Printf("Spotify Playlist Name: %s", fullPlaylist.Name)
	log.Printf("Spotify Playlist Description: %s", fullPlaylist.Name)
	log.Printf("fullPlaylist Url: %v", playlistUrl)
	log.Printf("-------------------------------------------------")
	log.Printf("Taking you to the playlist!")
	exec.Command("open", playlistUrl).Start()
	return playlistUrl, nil
}
