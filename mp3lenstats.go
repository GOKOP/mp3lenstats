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
	filenames := getArguments()[1:]
	durations := getDurations(filenames)

	fmt.Println("")

	fmt.Println("Średnia:", formatDuration(calcMeanDur(durations)))
	fmt.Println("Mediana:", formatDuration(calcMedianDur(durations)))
	fmt.Println("Max:", formatDuration(getMaxDur(durations)))
	fmt.Println("Min:", formatDuration(getMinDur(durations)))
}

func formatDuration(duration time.Duration) string {
	minutes := int(duration.Minutes())
	seconds := int(duration.Seconds()) - minutes*60

	return strconv.Itoa(minutes) + "m " + strconv.Itoa(seconds) + "s"
}

func getArguments() []string {
	args := os.Args;
	if(len(args) <= 1) {
		panic("dsgfegsegse")
	} else {
		return args
	}
}

func getDurations(filenames []string) []time.Duration {
	var durations []time.Duration

	for _, arg := range filenames {
		reader, err := os.Open(arg)

		if err != nil {
			panic(err)
		}

		decoder := mp3.NewDecoder(reader)

		var duration time.Duration
		var frame mp3.Frame
		skipped := 0

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
	var sumNanos int64

	for _, dur := range durations {
		sumNanos += dur.Nanoseconds()
	}

	avgNanos := sumNanos / int64(len(durations))
	return time.Duration(avgNanos)
}

func calcMedianDur(durations []time.Duration) time.Duration {
	var durationsNanos []int

	for _, dur := range durations {
		durationsNanos = append(durationsNanos, int(dur.Nanoseconds()))
	}

	sort.Ints(durationsNanos)

	if len(durationsNanos)%2 == 0 {
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
