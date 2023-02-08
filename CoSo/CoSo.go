package coso

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	spotify "github.com/jasondborneman/cosoradio/Spotify"
	tools "github.com/jasondborneman/cosoradio/Tools"
	youtube "github.com/jasondborneman/cosoradio/YouTube"
	yt "google.golang.org/api/youtube/v3"
)

func ReadSongsFromFixture(service yt.Service) ([]spotify.Song, error) {
	data, err := ioutil.ReadFile("Fixtures/songs.json")
	if err != nil {
		fmt.Println(fmt.Sprintf("Error reading song fixture: [%s]", err.Error()))
		return nil, err
	}
	stringData := string(data)
	var songs []spotify.Song
	err = json.Unmarshal([]byte(stringData), &songs)
	if err != nil {
		log.Printf("Error unmarshalling song fixture: [%v]", err)
		return nil, err
	} else {
		log.Printf("Success %d reading songs from fixture.", len(songs))
	}
	songs, err = AddYouTubeTitlesToSongs(service, songs)
	if err != nil {
		log.Printf("Error adding youtube title to song struct: %v", err)
		return nil, err
	}
	return songs, nil
}

func AddYouTubeTitlesToSongs(googleService yt.Service, songs []spotify.Song) ([]spotify.Song, error) {
	var retval []spotify.Song
	for i, song := range songs {
		title, tags, err := youtube.GetTitleFromVideo(googleService, song.YouTubeUrl)
		if err != nil {
			log.Printf("Error getting title from video %v", err)
		}
		log.Printf("%d : %s : %s", i, title, song.YouTubeUrl)
		updatedSong := spotify.Song{
			YouTubeUrl:    song.YouTubeUrl,
			YouTubeTitle:  title,
			RecommendedBy: song.RecommendedBy,
			YouTubeTags:   tags,
		}
		retval = append(retval, updatedSong)
	}
	return retval, nil
}

func GetSongsFromCoSo(cosoToken string) ([]spotify.Song, error) {
	// log.Printf("GetSongsFrmCoSo: Not Yet Implemented")
	// return nil, errors.New("GetSongsFromCoSo: Not Yet Implemented")
	search_url := "https://counter.social/api/v1/search?type=statuses&q=cosomusic"
	req, err := http.NewRequest("GET", search_url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", cosoToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response from GetSongsFromCoSo Search.\n[ERROR] -", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
		return nil, err
	}
	log.Println(string([]byte(body)))
	return []spotify.Song{}, nil
}

func TootSongs(songs []spotify.Song, content string, cosoToken string) error {
	content = tools.TruncateString(content, 500)
	return errors.New("TootSongs: Not Yet Implemented")
}
