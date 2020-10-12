package main

import (
	"os"
	"fmt"
	"time"
	"sort"
	"strconv"
	"strings"
	"github.com/tcolgate/mp3"
)

func main() {
	// get language in the form of "en", "pl", "ru" etc. on Unix systems
	// its not perfect but the worst case scenario is that the user gets English
	// as for non-Unicode locales, I'll just hope that no one will run this on them ;)
	lang := strings.Split(os.Getenv("LANG"), "_")[0]
	locale := createLocale(lang)

	// ignoring the command name
	filenames := getArguments(locale)[1:]
	durationsMap, number := getAndPrintDurations(filenames)
	durationsArray := extractDurations(durationsMap)

	fmt.Println("")

	fmt.Println(locale["fileNumber"], number)
	fmt.Println(locale["mean"], formatDuration(calcMeanDur(durationsArray)))
	fmt.Println(locale["median"], formatDuration(calcMedianDur(durationsArray)))

	longestFile, longestDur := getMaxDur(durationsMap)
	shortestFile, shortestDur := getMinDur(durationsMap)

	fmt.Println(locale["max"], formatDuration(longestDur), longestFile)
	fmt.Println(locale["min"], formatDuration(shortestDur), shortestFile)
}

func dontPanic(err string) {
	fmt.Println(err)
	os.Exit(1)
}

// I don't like how time.Duration is formatted by default
func formatDuration(duration time.Duration) string {
	minutes := int(duration.Minutes())
	seconds := int(duration.Seconds()) - minutes*60

	return strconv.Itoa(minutes) + "m " + strconv.Itoa(seconds) + "s"
}

func getArguments(locale map[string]string) []string {
	args := os.Args;
	if(len(args) <= 1) {
		dontPanic(locale["noArgument"])
		return []string{"pancakes"} // dead code but the compiler screams at me otherwise

	} else {
		return args
	}
}

func getAndPrintDurations(filenames []string) (map[string]time.Duration, int) {
	durations := make(map[string]time.Duration)
	number := 0 // number of files

	for _, arg := range filenames {
		reader, err := os.Open(arg)

		if err != nil {
			dontPanic(err.Error())
		}

		decoder := mp3.NewDecoder(reader)

		var duration time.Duration // total duration of the song
		var frame mp3.Frame // current mp3 frame
		skipped := 0 // skipped bits

		for {
			if err := decoder.Decode(&frame, &skipped); err != nil {
				if err.Error() == "EOF" {
					break
				} else {
					dontPanic(err.Error())
				}
			}

			duration += frame.Duration()
		}

		durations[arg] = duration
		number += 1
		fmt.Println(arg, ":", formatDuration(duration))
	}

	return durations, number
}

func extractDurations(durMap map[string]time.Duration) []time.Duration {
	var durArray []time.Duration

	for _, value := range durMap {
		durArray = append(durArray, value)
	}

	return durArray
}

func calcMeanDur(durations []time.Duration) time.Duration {
	var sum time.Duration

	for _, dur := range durations {
		sum += dur
	}

	// converting length to Duration so that it can divide
	avg := sum / time.Duration(len(durations))
	return avg
}

func calcMedianDur(durations []time.Duration) time.Duration {
	// converting from Duration to int because of sort.Ints()
	var durationsNanos []int

	for _, dur := range durations {
		durationsNanos = append(durationsNanos, int(dur.Nanoseconds()))
	}

	sort.Ints(durationsNanos)

	if len(durationsNanos)%2 == 0 {
		// calculate mean of middle elements if length isn't even
		leftMid := durationsNanos[ len(durationsNanos)/2 ]
		rightMid := durationsNanos[ (len(durationsNanos)/2)-1 ]

		return time.Duration((leftMid + rightMid) / 2)
	} else {
		return time.Duration(durationsNanos[len(durationsNanos)/2])
	}
}

func getMaxDur(durations map[string]time.Duration) (string, time.Duration) {
	maxVal := time.Duration(0)
	maxKey := ""

	for key, dur := range durations {
		if dur > maxVal {
			maxVal = dur
			maxKey = key
		}
	}

	return maxKey, maxVal
}

func getMinDur(durations map[string]time.Duration) (string, time.Duration) {
	_, minVal := getMaxDur(durations) // seems dirty
	minKey := ""

	for key, dur := range durations {
		if dur < minVal {
			minVal = dur
			minKey = key
		}
	}

	return minKey, minVal
}

func createLocale(language string) map[string]string {
	switch language {
		case "pl":
			return map[string]string {
				"mean": "Średnia:",
				"median": "Mediana:",
				"max": "Max:",
				"min": "Min:",
				"fileNumber": "Liczba plików:",
				"noArgument": "Musisz podać argumenty!",
			}
		default:
			return map[string]string {
				"mean": "Mean:",
				"median": "Median:",
				"max": "Max:",
				"min": "Min:",
				"fileNumber": "Number of files:",
				"noArgument": "You must provide arguments!",
			}
	}
}
