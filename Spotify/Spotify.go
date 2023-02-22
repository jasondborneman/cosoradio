package spotify

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	tools "github.com/jasondborneman/cosoradio/Tools"
	spotifyapi "github.com/zmb3/spotify/v2"
)

func CreatePlaylist(client spotifyapi.Client, ctx context.Context, songs []Song, recommenders string) (string, error) {
	var playlistUrl string
	user, err := client.CurrentUser(ctx)
	if err != nil {
		log.Printf("Error getting current user: %v", err)
		return "", err
	}
	t := time.Now()
	dateString := t.Format("01-02-2006")
	playlistName := fmt.Sprintf("CoSoRadio for %s", dateString)
	recommenders = strings.Replace(recommenders, "\n", ", ", -1)
	recommenders = strings.TrimPrefix(recommenders, ", ")
	playlistDescription := fmt.Sprintf("The counter.social #CoSoRadio playlist for %s.  Featuring recommendations from: %s", dateString, recommenders)
	playlistDescription = tools.TruncateString(playlistDescription, 300)
	fullPlaylist, err := client.CreatePlaylistForUser(ctx, user.ID, playlistName, playlistDescription, true, false)

	if err != nil {
		log.Printf("Error creating playlist %v", err)
		return "", err
	}
	playlistUrl = string(fullPlaylist.ExternalURLs["spotify"])

	// Playlist Thumbnail seems to be failing again for some reason. Takingit out for now.
	// png, err := os.Open("Images/cosoradio.png")
	// if err != nil {
	// 	log.Printf("Error opening thumbnail for playlist: %v", err)
	// 	return "", err
	// }
	// pngWithLabel, err := tools.AddDateToThumbnail(png, dateString, 1, 13)
	// if err != nil {
	// 	log.Printf("Error creating playslist thumbnail: %v", err)
	// 	return "", err
	// }
	// err = client.SetPlaylistImage(ctx, fullPlaylist.ID, pngWithLabel)
	// if err != nil {
	// 	log.Printf("Error adding thumbnail to playlist: %v", err)
	// 	return "", err
	// }

	var spotifyTrackIds []spotifyapi.ID
	for _, song := range songs {
		tag_query := ""
		if len(song.YouTubeTags) > 0 {
			tag_query = song.YouTubeTags[0]
		}
		searchQuery := fmt.Sprintf("%s %s", song.YouTubeTitle, tag_query)
		searchResult, err := client.Search(ctx, searchQuery, spotifyapi.SearchTypeTrack)
		if err != nil {
			log.Printf("error searching Spotify [query: %s]: %v", searchQuery, err)
			continue
		}
		if len(searchResult.Tracks.Tracks) > 0 {
			topResult := searchResult.Tracks.Tracks[0]
			spotifyTrackIds = append(spotifyTrackIds, topResult.ID)
		} else {
			log.Printf("\tNo tracks found for [%s]\n", song.YouTubeTitle)
			continue
		}
	}
	client.AddTracksToPlaylist(ctx, fullPlaylist.ID, spotifyTrackIds...)

	log.Printf("-----------------------------------------------------")
	log.Printf("Spotify Playlist Name: %s", fullPlaylist.Name)
	log.Printf("Spotify Playlist Description: %s", fullPlaylist.Description)
	log.Printf("Number Of Tracks: %d", len(spotifyTrackIds))
	log.Printf("fullPlaylist Url: %v", playlistUrl)
	log.Printf("-----------------------------------------------------")
	log.Printf("Taking you to the playlist!")
	exec.Command("open", playlistUrl).Start()
	return playlistUrl, nil
}
