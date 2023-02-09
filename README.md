# CoSoRadio
TL;DR: Eventualyl this will scrape [counter.social](https://counter.social) firehose for posts using the `#CoSoMusic` hashtag. For each one, ff they have a YouTube link as well, get the Title from the YouTube video and then add the song to a public Spotify playlist and share it out on CoSo.

## TechStack
Language: Golang

APIs: Spotify, YouTube [Future: counter.social]

Runs: Locally (for now). Due to Spotify's authentication scheme, you're forced through a web login. There's no way to do a "headless" API integration with spotify that I've found yet. More research needed into that for both Spotify and Google APIs.

## Build and Run
```
$ cd cmd
$ go build .
$ GOOGLE_ID=<GOOGLE ID> \
  GOOGLE_SECRET=<GOOGLE SECRET> \
  SPOTIFY_ID=<SPOTIFY ID> \
  SPOTIFY_SECRET=<SPOTIFY SECRET> \
  COSO_CLIENT_KEY=<COSO CLIENT KEY/ID> \
  COSO_CLIENT_SECRET=<COSO CLIENT SECRET> \
  CSM_MAKE_PLAYLIST=true \
  CSM_DO_TOOT=true \
  CSM_SCRAPE_COSO=true \
  ./cmd
```

### Workflow
1. Auth to Spotify and create Spotify Client
2. Auth to Google and create Google API Service
3. Pull music from source
    * Currently just uses the CoSo `timeline` public API. The `search` API isn't returning results when searching for statuses. 
    * Currently only getting the latest 40. Eventually if I can get the `search` API working, I can page through results since the last time/status ID and get more.
4. Given the slice of Song objects, take the YouTubeUrl, look it up with the YouTube API, and get the title and tags from the video and add them to the correct Song object.
5. Create an empty Spotify playlist
6. For each Song, search for the Title + First Tag in Spotify
7. Add the resulting songs to the playlist
8. Post the Spotify Playlist to CoSo as `#CoSoRadio`

Run this once a day. Eventually I'll see if I can't get it running once a day as a cloud function in GCP.

![](Readme_Img/cosoradio_cli.png)

### Required Envs

* SPOTIFY_ID : The spotify ID [App set up with a redirect url of `http://localhost:8080/callback`]
* SPOTIFY_SECRET : The spotify Secret
* GOOGLE_ID : Google API Key ID [Oauth2 with scope to read only from youtube api]
* GOOGLE_SECRET: Google API Secret
* COSO_CLIENT_KEY : Counter.social API Key (In documentation referred to as client id)
* COSO_CLIENT_SECRET : Counter.social API Secret

### Options

* CSM_MAKE_PLAYLIST : when =true will make the spotify playlist. Defaults to false
* CSM_DO_TOOT : when =true will send the toot to coso. Only will post to CoSo if the playlist is crated too. Defaults to false. 
* CSM_SCRAPE_COSO : when =true will scrape coso for #CoSoMusic hashtag for the previous day. when =false will pull from /cmd/Fixtures/songs.json. Defaults to false.