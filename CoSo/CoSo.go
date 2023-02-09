package coso

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/uuid"
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

func GetSongsFromCoSo(service yt.Service, cosoToken string) ([]spotify.Song, error) {
	retval, err := GetSongsFromCoSoTimeline(cosoToken)
	if err != nil {
		return nil, err
	}
	if len(retval) == 0 {
		err = errors.New("0 songs returned from CoSo")
		return nil, err
	}
	songs, err := AddYouTubeTitlesToSongs(service, retval)
	if err != nil {
		log.Printf("Error adding youtube title to song struct: %v", err)
		return nil, err
	}
	return songs, nil
}

func GetSongsFromCoSoSearch(cosoToken string) ([]spotify.Song, error) {
	search_url := "https://counter.social/api/v2/search?type=statuses&q=cosomusic&limit=1"
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
	// Make sure to only pick ones that have a youtube.com or youtu.be link in them.
	return []spotify.Song{}, nil
}

func GetSongsFromCoSoTimeline(cosoToken string) ([]spotify.Song, error) {
	timeline_url := "https://counter.social/api/v1/timelines/tag/cosomusic?limit=100"
	req, err := http.NewRequest("GET", timeline_url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", cosoToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error on response from GetSongsFromCoSo Search.\n[ERROR] - %s\n", err)
		return nil, err
	}
	defer resp.Body.Close()
	// Make sure to only pick ones that have a youtube.com or youtu.be link in them.
	var statuses_unfiltered TimelineStatus
	var retval []spotify.Song
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&statuses_unfiltered)
	if err != nil {
		message := fmt.Sprintf("Error Parsing CoSo Timeline Response: %v", err)
		log.Println(message)
		return nil, err
	}
	for _, status := range statuses_unfiltered {
		if status.Card.URL != "" && (strings.Contains(status.Card.URL, "youtube.com") || strings.Contains(status.Card.URL, "youtu.be")) {
			s := spotify.Song{}
			s.RecommendedBy = status.Account.Username
			s.YouTubeUrl = status.Card.URL
			retval = append(retval, s)
		}
	}
	return retval, nil
}

func TootSongs(songs []spotify.Song, content string, cosoToken string) error {
	content = tools.TruncateString(content, 500)
	fmt.Println(content)
	fmt.Println("TOOTING...")
	toot_url := fmt.Sprintf("https://counter.social/api/v1/statuses?status=%s", url.QueryEscape(content))
	req, err := http.NewRequest("POST", toot_url, nil)
	if err != nil {
		log.Printf("Error creating new POST request to toot.\n[ERROR] - %s\n", err)
		return err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", cosoToken))
	req.Header.Add("Idempotency-Key", uuid.New().String())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error on response from Post Toot.\n[ERROR] - %s\n", err)
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading response body.\n[ERROR] - %s\n", err)
			return err
		}
		log.Printf("Non-200 Status Code posting Toot\n[%d] [RESPONSE]\n%s\n", resp.StatusCode, string(body))
		return errors.New(fmt.Sprintf("Non-200 Status Code:%d", resp.StatusCode))
	}
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	return nil
}
