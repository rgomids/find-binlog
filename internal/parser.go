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

func ExtractBinlogPositionFromOutput(output string, targetDate time.Time) (file string, pos int64, ts time.Time, line int, err error) {
	// regex helpers to capture the binlog file, the position and the event timestamp
	rFile := regexp.MustCompile(`processing log events from (\S+),`)
	rPos := regexp.MustCompile(`^# at\s+(\d+)`)
	rTS := regexp.MustCompile(`^###\s*SET\s+TIMESTAMP=(\d+)`)
	rOrig := regexp.MustCompile(`^#\s*original_commit_timestamp=.*\((\d{4}-\d{2}-\d{2})`)

	scanner := bufio.NewScanner(strings.NewReader(output))

	var (
		currentPos int64
		havePos    bool
		binlogFile string
		lineNumber int
	)

	for scanner.Scan() {
		lineNumber++
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
				return "", 0, time.Time{}, 0, err
			}
			currentPos = p
			havePos = true
			continue
		}

		// handle MySQL 8.0+ commit timestamp comment
		if m := rOrig.FindStringSubmatch(line); m != nil && havePos {
			eventDate, err := time.Parse("2006-01-02", m[1])
			if err != nil {
				return "", 0, time.Time{}, 0, err
			}

			if !eventDate.Before(targetDate) {
				return binlogFile, currentPos, eventDate, lineNumber, nil
			}

			// event is before the target, discard current position
			havePos = false
			continue
		}

		// when a timestamp line appears, associate it with the previously
		// found position and compare with the target date
		if m := rTS.FindStringSubmatch(line); m != nil && havePos {
			unixVal, err := strconv.ParseInt(m[1], 10, 64)
			if err != nil {
				return "", 0, time.Time{}, 0, err
			}

			eventTime := time.Unix(unixVal, 0).UTC()
			if !eventTime.Before(targetDate) {
				return binlogFile, currentPos, eventTime, lineNumber, nil
			}

			// event is before the target, discard current position
			havePos = false
		}
	}

	if err := scanner.Err(); err != nil {
		return "", 0, time.Time{}, 0, err
	}

	return "", 0, time.Time{}, 0, ErrNoEventFound
}

// ExtractClosestEventFromOutput returns the binlog event closest to the target
// date, regardless of whether it happens before or after the target. If no
// events are found, ErrNoEventFound is returned.
func ExtractClosestEventFromOutput(output string, targetDate time.Time) (file string, pos int64, ts time.Time, line int, err error) {
	rFile := regexp.MustCompile(`processing log events from (\S+),`)
	rPos := regexp.MustCompile(`^# at\s+(\d+)`)
	rTS := regexp.MustCompile(`^###\s*SET\s+TIMESTAMP=(\d+)`)
	rOrig := regexp.MustCompile(`^#\s*original_commit_timestamp=.*\((\d{4}-\d{2}-\d{2})`)

	scanner := bufio.NewScanner(strings.NewReader(output))

	var (
		currentPos int64
		havePos    bool
		binlogFile string
		lineNumber int

		bestDiff time.Duration
		bestPos  int64
		bestTS   time.Time
		bestLine int
		haveBest bool
	)

	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())

		if binlogFile == "" {
			if m := rFile.FindStringSubmatch(line); m != nil {
				binlogFile = m[1]
				continue
			}
		}

		if m := rPos.FindStringSubmatch(line); m != nil {
			p, err := strconv.ParseInt(m[1], 10, 64)
			if err != nil {
				return "", 0, time.Time{}, 0, err
			}
			currentPos = p
			havePos = true
			continue
		}

		var eventTime time.Time
		var ok bool
		if m := rOrig.FindStringSubmatch(line); m != nil && havePos {
			d, err := time.Parse("2006-01-02", m[1])
			if err != nil {
				return "", 0, time.Time{}, 0, err
			}
			eventTime = d
			ok = true
		} else if m := rTS.FindStringSubmatch(line); m != nil && havePos {
			unixVal, err := strconv.ParseInt(m[1], 10, 64)
			if err != nil {
				return "", 0, time.Time{}, 0, err
			}
			eventTime = time.Unix(unixVal, 0).UTC()
			ok = true
		}

		if ok {
			diff := targetDate.Sub(eventTime)
			if diff < 0 {
				diff = -diff
			}
			if !haveBest || diff < bestDiff {
				bestDiff = diff
				bestPos = currentPos
				bestTS = eventTime
				bestLine = lineNumber
				haveBest = true
			}
			havePos = false
		}
	}

	if err := scanner.Err(); err != nil {
		return "", 0, time.Time{}, 0, err
	}

	if !haveBest {
		return "", 0, time.Time{}, 0, ErrNoEventFound
	}

	return binlogFile, bestPos, bestTS, bestLine, nil
}
