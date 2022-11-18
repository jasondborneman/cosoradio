package coso

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"

	spotify "github.com/jasondborneman/cosoradio/Spotify"
	tools "github.com/jasondborneman/cosoradio/Tools"
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

func TootSongs(songs []spotify.Song, content string) error {
	content = tools.TruncateString(content, 500)
	log.Println(content)
	return errors.New("TootSongs: Not Yet Implemented")
}
