package cmd

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	host     string
	dateStr  string
	port     int
	user     string
	password string
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "find-binlog",
		Short: "Find first binlog event after a specific date",
		RunE: func(cmd *cobra.Command, args []string) error {
			if host == "" {
				return fmt.Errorf("host is required")
			}
			if dateStr == "" {
				return fmt.Errorf("date is required")
			}
			targetDate, err := time.Parse("2006-01-02", dateStr)
			if err != nil {
				return fmt.Errorf("invalid date format (expected YYYY-MM-DD)")
			}

			ctx := context.Background()

			binlogs, err := listBinlogs(ctx)
			if err != nil {
				return err
			}
			for _, f := range binlogs {
				pos, ts, found, err := scanBinlog(ctx, f, targetDate)
				if err != nil {
					return err
				}
				if found {
					fmt.Printf("Binlog: %s\nPosition: %s\nTimestamp: %s\n", f, pos, ts.Format("2006-01-02"))
					return nil
				}
			}
			fmt.Printf("Nenhum evento encontrado a partir de %s\n", dateStr)
			return nil
		},
	}

	cmd.Flags().StringVarP(&host, "host", "H", "", "MySQL host")
	cmd.Flags().StringVarP(&dateStr, "date", "d", "", "Target date (YYYY-MM-DD)")
	cmd.Flags().IntVarP(&port, "port", "P", 3306, "MySQL port")
	cmd.Flags().StringVarP(&user, "user", "u", "", "MySQL user")
	cmd.Flags().StringVarP(&password, "password", "p", "", "MySQL password")
	cmd.MarkFlagRequired("host")
	cmd.MarkFlagRequired("date")

	return cmd
}

func listBinlogs(ctx context.Context) ([]string, error) {
	args := []string{
		"-h", host,
		"-P", strconv.Itoa(port),
		"-u", user,
		"-p" + password,
		"-N",
		"-e", "SHOW BINARY LOGS;",
	}
	out, err := exec.CommandContext(ctx, "mysql", args...).Output()
	if err != nil {
		return nil, err
	}
	var logs []string
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) > 0 {
			logs = append(logs, fields[0])
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return logs, nil
}

func scanBinlog(ctx context.Context, file string, target time.Time) (string, time.Time, bool, error) {
	args := []string{
		"--read-from-remote-server",
		"--host=" + host,
		"--user=" + user,
		"--password=" + password,
		"--port=" + strconv.Itoa(port),
		"--verbose",
		"--base64-output=DECODE-ROWS",
		file,
	}

	cmd := exec.CommandContext(ctx, "./pkg/bin/mysqlbinlog", args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", time.Time{}, false, err
	}
	if err := cmd.Start(); err != nil {
		return "", time.Time{}, false, err
	}

	rPos := regexp.MustCompile(`^# at\s+(\d+)`)
	rTS := regexp.MustCompile(`^###\s*SET\s+TIMESTAMP=(\d+)`)

	scanner := bufio.NewScanner(stdout)
	var (
		pos string
		ts  time.Time
	)
	for scanner.Scan() {
		line := scanner.Text()
		if m := rPos.FindStringSubmatch(line); m != nil {
			pos = m[1]
			continue
		}
		if m := rTS.FindStringSubmatch(line); m != nil && pos != "" {
			unixVal, err := strconv.ParseInt(m[1], 10, 64)
			if err != nil {
				continue
			}
			ts = time.Unix(unixVal, 0)
			dateOnly, _ := time.Parse("2006-01-02", ts.Format("2006-01-02"))
			if !dateOnly.Before(target) {
				cmd.Process.Kill()
				cmd.Wait()
				return pos, dateOnly, true, nil
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return "", time.Time{}, false, err
	}
	cmd.Wait()
	return "", time.Time{}, false, nil
}

func Execute() {
	if err := NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
