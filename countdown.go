package main

import (
	"flag"
	"fmt"
	"github.com/davidscholberg/go-durationfmt"
	"os"
	"time"
)

var (
	targetTime     time.Time
	update         time.Duration
	prefix         string
	postfix        string
	durationFormat string
	doneMessage    string
	outputFile     string
)

// Cached last write to make sure we don't update the file unless we need to.
var lastTimeWritten string

func main() {
	if err := handleFlags(); err != nil {
		panic(err)
	}
	fmt.Printf("counting down to: %s\n", targetTime.String())
	fmt.Printf("output writing to: %s (format: %s)\n", outputFile, durationFormat)

	if err := countdown(); err != nil {
		panic(err)
	}
	fmt.Printf("target time reached!\n")
}

func countdown() error {
	file, err := os.OpenFile(outputFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	ticker := time.NewTicker(update)
	for {
		select {
		case t := <-ticker.C:
			if t.After(targetTime) {
				if doneMessage != "" {
					if err := writeFile(file, doneMessage); err != nil {
						return err
					}
				} else {
					if err := writeDuration(file, targetTime); err != nil {
						return err
					}
				}
				return nil
			}
			if err := writeDuration(file, t); err != nil {
				return err
			}
		}
	}
}

func writeDuration(file *os.File, t time.Time) error {
	dur := targetTime.Sub(t)
	fmtDur, err := durationfmt.Format(dur, durationFormat)
	if err != nil {
		return err
	}

	return writeFile(file, prefix+fmtDur+postfix)
}

func writeFile(file *os.File, message string) error {
	if lastTimeWritten == message {
		// No need to write again.
		return nil
	}
	if err := file.Truncate(0); err != nil {
		return err
	}
	if _, err := file.Seek(0, 0); err != nil {
		return err
	}
	if _, err := file.WriteString(message); err != nil {
		return err
	}
	return nil
}

func handleFlags() error {
	target := flag.String("target", "10:00:00", "the time to count down to in hh:mm:ss format (24-hr)")
	up := flag.String("update", "100ms", "the amount of time to wait until checking to see if the file needs updating")
	flag.StringVar(&prefix, "prefix", "", "a prefix to add beginning of the output")
	flag.StringVar(&postfix, "postfix", "", "a postfix to add to end of the output")
	flag.StringVar(&durationFormat, "format", "%00h:%00m:%00s", "the duration format to output (see https://github.com/davidscholberg/go-durationfmt)")
	flag.StringVar(&outputFile, "output", "countdown.txt", "the output file to update")
	flag.StringVar(&doneMessage, "complete", "", "after the target time is hit this will be written to the file (if set, no prefix/postfix)")
	flag.Parse()

	upd, err := time.ParseDuration(*up)
	if err != nil {
		return err
	}
	update = upd

	parsed, err := time.Parse("15:04:05", *target)
	if err != nil {
		return err
	}
	targetTime = time.Now()
	targetTime = time.Date(targetTime.Year(), targetTime.Month(), targetTime.Day(), parsed.Hour(), parsed.Minute(), parsed.Second(), 0, targetTime.Location())
	if targetTime.Before(time.Now()) {
		// Must be counting down until a time tomorrow
		targetTime = targetTime.Add(24 * time.Hour)
	}
	return nil
}
