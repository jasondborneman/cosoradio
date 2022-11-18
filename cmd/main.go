package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	csm "github.com/jasondborneman/cosoradio/CoSoRadioApp"
	spotifyapi "github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

func main() {
	scrapeCoSo := os.Getenv("CSM_SCRAPE_COSO") == "true"
	doToot := os.Getenv("CSM_DO_TOOT") == "true"
	var client *spotifyapi.Client
	client = StartSpotifyAuthProcess()
	os.Mkdir("Images/Temp/", 0755)
	defer os.RemoveAll("Images/Temp/")
	err := csm.Run(*client, scrapeCoSo, doToot)
	if err != nil {
		log.Printf("Error running the application: %v", err)
	}
}

const redirectURI = "http://localhost:8080/callback"

var (
	auth = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI), spotifyauth.WithScopes(
		spotifyauth.ScopeUserReadPrivate,
		spotifyauth.ScopePlaylistModifyPublic,
		spotifyauth.ScopeImageUpload))
	ch    = make(chan *spotifyapi.Client)
	state = "abc123"
)

func StartSpotifyAuthProcess() *spotifyapi.Client {
	// first start an HTTP server
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	url := auth.AuthURL(state)
	fmt.Println("Taking you to Spotify for authentication!")
	exec.Command("open", url).Start()

	// wait for auth to complete
	client := <-ch

	// use the client to make calls that require authorization
	user, err := client.CurrentUser(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("You are logged in as:", user.ID)
	return client
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}

	// use the token to get an authenticated client
	client := spotifyapi.New(auth.Client(r.Context(), tok))
	fmt.Fprintf(w, "Login Completed!")
	ch <- client
}
