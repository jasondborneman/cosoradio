package coso

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	spotify "github.com/jasondborneman/cosoradio/Spotify"
	"github.com/mpvl/unique"
)

func ReadSongsFromFixture() ([]spotify.Song, error) {
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
		return songs, nil
	}
}

func GetSongsFromCoSo() ([]spotify.Song, error) {
	log.Printf("GetSongsFrmCoSo: Not Yet Implemented")
	return nil, errors.New("GetSongsFromCoSo: Not Yet Implemented")
}

func TootSongs(spotifyPlaylist string, songs []spotify.Song) error {
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
	tootMessage := fmt.Sprintf("#CoSoRadio for %s!\n%s\n\nFeaturing music recommendations from: %s",
		dateString,
		spotifyPlaylist,
		recommendersString)
	log.Println(tootMessage)
	return errors.New("TootSongs: Not Yet Implemented")
}
