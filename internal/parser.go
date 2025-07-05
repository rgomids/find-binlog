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
	// regex helpers to capture the binlog file, the position and the event timestamp
	rFile := regexp.MustCompile(`processing log events from (\S+),`)
	rPos := regexp.MustCompile(`^# at\s+(\d+)`)
	rTS := regexp.MustCompile(`^###\s*SET\s+TIMESTAMP=(\d+)`)

	scanner := bufio.NewScanner(strings.NewReader(output))

	var (
		currentPos int64
		havePos    bool
		binlogFile string
	)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// try to capture the binlog file name (only once)
		if binlogFile == "" {
			if m := rFile.FindStringSubmatch(line); m != nil {
				binlogFile = m[1]
				continue
			}
		}

		// capture the current position
		if m := rPos.FindStringSubmatch(line); m != nil {
			p, err := strconv.ParseInt(m[1], 10, 64)
			if err != nil {
				return "", 0, time.Time{}, err
			}
			currentPos = p
			havePos = true
			continue
		}

		// when a timestamp line appears, associate it with the previously
		// found position and compare with the target date
		if m := rTS.FindStringSubmatch(line); m != nil && havePos {
			unixVal, err := strconv.ParseInt(m[1], 10, 64)
			if err != nil {
				return "", 0, time.Time{}, err
			}

			eventTime := time.Unix(unixVal, 0).UTC()
			if !eventTime.Before(targetDate) {
				return binlogFile, currentPos, eventTime, nil
			}

			// event is before the target, discard current position
			havePos = false
		}
	}

	if err := scanner.Err(); err != nil {
		return "", 0, time.Time{}, err
	}

	return "", 0, time.Time{}, ErrNoEventFound
}
