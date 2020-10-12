Simple tool to show stats about the length of given mp3 files.
Made for fun cause I wanted to know average length of my music.

## Build
You need Go
```
go get github.com/tcolgate/mp3
go build
```

## Usage
```
mp3lenstats music/*.mp3
```
will give you the number of given files, length of each track, mean average length of all tracks, median of all tracks, the longest track and the shortest track.
