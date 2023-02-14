package youtube

import (
	"errors"
	"log"
	"net/url"
	"strings"

	yt "google.golang.org/api/youtube/v3"
)

func GetTitleFromVideo(service yt.Service, url string) (string, []string, error) {
	vId, err := GetVideoIdFromUrl(url)
	if err != nil {
		log.Printf("\nError getting Video ID from URL %s: %v\n", url, err)
		return "", nil, err
	}

	var part []string
	part = append(part, "id")
	part = append(part, "snippet")
	part = append(part, "contentDetails")
	call := service.Videos.List(part).Id(vId)
	resp, err := call.Do()
	if err != nil {
		log.Printf("\nError getting video id [%s] from yt api: %v\n", vId, err)
		return "", nil, err
	}
	if len(resp.Items) == 0 {
		message := "\nno items from GetTitleFromVideo response"
		log.Println(message)
		return "", nil, errors.New(message)
	} else {
		return resp.Items[0].Snippet.Title, resp.Items[0].Snippet.Tags, nil
	}
}

func GetVideoIdFromUrl(urlStr string) (string, error) {
	urlObj, err := url.Parse(urlStr)
	if err != nil {
		log.Printf("Error parsing url string %s: %v", urlStr, err)
		return "", err
	}
	params, err := url.ParseQuery(urlObj.RawQuery)
	if err != nil {
		log.Printf("Error Parsing query string: %v", err)
	}
	if strings.Contains(urlStr, "watch?v") {
		return params["v"][0], nil
	} else {
		splits := strings.Split(urlStr, "/")
		return splits[len(splits)-1], nil
	}
}
