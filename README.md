# CoSoRadio
This'll get filled out more another day.

TL;DR: Eventualyl this will scrape [counter.social](https://counter.social) firehose for posts using the #CoSoMusic hashtag. For each one, ff they have a YouTube link as well, get the Title from the YouTube video and then add the song to a public Spotify playlist and share it out on CoSo.

## TechStack
Language: Golang

APIs: Spotify [Future: counter.social, YouTube]

Runs: Locally. Due to Spotify's authentication scheme, you're forced through a web login. There's no way to do a "headless" API integration with spotify.