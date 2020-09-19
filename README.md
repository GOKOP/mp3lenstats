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
mp3length music/*.mp3
```
will give you length of each track, mean average length of all tracks, median of all tracks, the largest length and the shortest length.