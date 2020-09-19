package main

import (
	"os"
	"fmt"
	"time"
	"sort"
	"strconv"
	"github.com/tcolgate/mp3"
)

func main() {
	// ignoring the command name
	filenames := getArguments()[1:]
	durations := getAndPrintDurations(filenames)

	fmt.Println("")

	fmt.Println("Mean:", formatDuration(calcMeanDur(durations)))
	fmt.Println("Median:", formatDuration(calcMedianDur(durations)))
	fmt.Println("Max:", formatDuration(getMaxDur(durations)))
	fmt.Println("Min:", formatDuration(getMinDur(durations)))
}

// I don't like how time.Duration is formatted by default
func formatDuration(duration time.Duration) string {
	minutes := int(duration.Minutes())
	seconds := int(duration.Seconds()) - minutes*60

	return strconv.Itoa(minutes) + "m " + strconv.Itoa(seconds) + "s"
}

func getArguments() []string {
	args := os.Args;
	if(len(args) <= 1) {
		panic("You must give an argument!")
	} else {
		return args
	}
}

func getAndPrintDurations(filenames []string) []time.Duration {
	var durations []time.Duration

	for _, arg := range filenames {
		reader, err := os.Open(arg)

		if err != nil {
			panic(err)
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
					panic(err)
				}
			}

			duration += frame.Duration()
		}

		durations = append(durations, duration)
		fmt.Println(arg, ":", formatDuration(duration))
	}

	return durations
}

func calcMeanDur(durations []time.Duration) time.Duration {
	// converting from Duration to int for arithmetics
	var sumNanos int64

	for _, dur := range durations {
		sumNanos += dur.Nanoseconds()
	}

	avgNanos := sumNanos / int64(len(durations))
	return time.Duration(avgNanos)
}

func calcMedianDur(durations []time.Duration) time.Duration {
	// converting from Duration to int for arithmetics
	// (not int64 because of sort.Ints(), should be the same on modern systems I think)
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

func getMaxDur(durations []time.Duration) time.Duration {
	max := time.Duration(0)

	for _,dur := range durations {
		if dur > max {
			max = dur
		}
	}

	return max
}

func getMinDur(durations []time.Duration) time.Duration {
	min := durations[0]

	for _,dur := range durations {
		if dur < min {
			min = dur
		}
	}

	return min
}
