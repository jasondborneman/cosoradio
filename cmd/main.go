package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"

	csm "github.com/jasondborneman/cosoradio/CoSoRadioApp"
	spotifyapi "github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	oauth2 "golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	yt "google.golang.org/api/youtube/v3"
)

const spotifyRedirectURI = "http://localhost:8080/callback"
const youtubeRedirectURI = "http://localhost:8080/youtube"

var (
	spotifyAuth = spotifyauth.New(spotifyauth.WithRedirectURL(spotifyRedirectURI), spotifyauth.WithScopes(
		spotifyauth.ScopeUserReadPrivate,
		spotifyauth.ScopePlaylistModifyPublic,
		spotifyauth.ScopeImageUpload))
	chs   = make(chan *spotifyapi.Client)
	chg   = make(chan string)
	state = "abc123"
)

func main() {
	scrapeCoSo := os.Getenv("CSM_SCRAPE_COSO") == "true"
	doToot := os.Getenv("CSM_DO_TOOT") == "true"

	log.Println("     ████                              ")
	log.Println("   ██░░░░/█                            ")
	log.Println("  ██░░░░/░░░██                          ")
	log.Println("  ██░░░/░░░░██                          ")
	log.Println("██░#CoSoRadio░██                        ")
	log.Println("██░░░░░░░░░░░░██                        ")
	log.Println("██░░░░░░░░░░░░██                        ")
	log.Println("  ██░░░░░░░░██                          ")
	log.Println("    ████████")
	log.Println()
	log.Println(fmt.Sprintf("----------] Scrape CoSo: %t", scrapeCoSo))
	log.Println(fmt.Sprintf("----------] Do Toot: %t", doToot))
	log.Println("-----------------------------------------")
	log.Println("Starting HTTP server for auth callbacks...")
	http.HandleFunc("/callback", completeSpotifyAuth)
	http.HandleFunc("/youtube", completeYoutubeAuth)
	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	var spotifyClient *spotifyapi.Client
	spotifyClient = StartSpotifyAuthProcess()

	os.Mkdir("Temp/", 0755)
	defer os.RemoveAll("Temp/")
	googleService, err := StartGoogleAuthProcess()
	if err != nil {
		log.Printf("Error Creating Google API Service: %v", err)
	}

	os.Mkdir("Images/Temp/", 0755)
	defer os.RemoveAll("Images/Temp/")
	err = csm.Run(*spotifyClient, *googleService, scrapeCoSo, doToot)
	if err != nil {
		log.Printf("Error running the application: %v", err)
	}
}

func StartSpotifyAuthProcess() *spotifyapi.Client {
	url := spotifyAuth.AuthURL(state)
	log.Println("SPOTIFY AUTH")
	log.Println("\tTaking you to Spotify for authentication!")
	exec.Command("open", url).Start()

	// wait for auth to complete
	client := <-chs
	return client
}

func completeSpotifyAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := spotifyAuth.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "\tCouldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}

	// use the token to get an authenticated client
	client := spotifyapi.New(spotifyAuth.Client(r.Context(), tok))
	chs <- client
}

func StartGoogleAuthProcess() (*yt.Service, error) {
	log.Println("GOOGLE AUTH")
	ctx := context.Background()
	WriteToClientSecretJson()

	b, err := ioutil.ReadFile("Temp/client_secret.json")
	if err != nil {
		log.Printf("\tError reading client_secret.json for google auth: %v", err)
		return nil, err
	}

	config, err := google.ConfigFromJSON(b, yt.YoutubeReadonlyScope)
	client, err := getGoogleClient(ctx, config)
	service, err := yt.New(client)

	if err != nil {
		log.Printf("\tError creating YouTube Service: %v", err)
		return nil, err
	}
	return service, nil
}

func WriteToClientSecretJson() error {
	googleAuthJson := `{"installed":{"client_id":"THEID","project_id":"cosoradio","auth_uri":"https://accounts.google.com/o/oauth2/auth","token_uri":"https://oauth2.googleapis.com/token","auth_provider_x509_cert_url":"https://www.googleapis.com/oauth2/v1/certs","client_secret":"THESECRET","redirect_uris":["THEREDIRECT"]}}`
	googleAuthJson = strings.Replace(googleAuthJson, "THEID", os.Getenv("GOOGLE_ID"), 1)
	googleAuthJson = strings.Replace(googleAuthJson, "THESECRET", os.Getenv("GOOGLE_SECRET"), 1)
	googleAuthJson = strings.Replace(googleAuthJson, "THEREDIRECT", youtubeRedirectURI, 1)
	f, err := os.Create("Temp/client_secret.json")
	defer f.Close()
	if err != nil {
		log.Printf("Error creating client_secret.json for google auth: %v", err)
		return err
	}
	_, err = f.WriteString(googleAuthJson)
	if err != nil {
		log.Printf("Error writing to client_secret.json for google auth: %v", err)
		return err
	}

	return nil
}

// getgetGoogleClientClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getGoogleClient(ctx context.Context, config *oauth2.Config) (*http.Client, error) {
	tok := getGoogleTokenFromWeb(config)
	return config.Client(ctx, tok), nil
}

func completeYoutubeAuth(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	chg <- code
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getGoogleTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	log.Println("\tTaking you to Google for authentication!")
	exec.Command("open", authURL).Start()

	code := <-chg
	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}
