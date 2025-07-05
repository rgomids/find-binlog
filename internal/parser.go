package internal

import (
	"bufio"
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var ErrNoEventFound = errors.New("no compatible events")

func ExtractBinlogPositionFromOutput(output string, targetDate time.Time) (file string, pos int64, ts time.Time, err error) {
	rFile := regexp.MustCompile(`processing log events from (\S+),`)
	rPos := regexp.MustCompile(`^# at\s+(\d+)`)
	rTS := regexp.MustCompile(`^###\s*SET\s+TIMESTAMP=(\d+)`)

	scanner := bufio.NewScanner(strings.NewReader(output))
	var (
		currentPos string
		foundFile  string
	)
	for scanner.Scan() {
		line := scanner.Text()
		if foundFile == "" {
			if m := rFile.FindStringSubmatch(line); m != nil {
				foundFile = m[1]
				continue
			}
		}
		if m := rPos.FindStringSubmatch(line); m != nil {
			currentPos = m[1]
			continue
		}
		if m := rTS.FindStringSubmatch(line); m != nil && currentPos != "" {
			unixVal, err := strconv.ParseInt(m[1], 10, 64)
			if err != nil {
				continue
			}
			ts = time.Unix(unixVal, 0)
			dateOnly, _ := time.Parse("2006-01-02", ts.Format("2006-01-02"))
			if !dateOnly.Before(targetDate) {
				p, err := strconv.ParseInt(currentPos, 10, 64)
				if err != nil {
					return "", 0, time.Time{}, err
				}
				if foundFile == "" {
					foundFile = ""
				}
				return foundFile, p, dateOnly, nil
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return "", 0, time.Time{}, err
	}
	return "", 0, time.Time{}, ErrNoEventFound
}
