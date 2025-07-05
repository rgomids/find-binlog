package internal

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	host      string
	dateStr   string
	port      int
	user      string
	password  string
	frameshot bool
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
				return fmt.Errorf("erro ao listar binlogs: %w", err)
			}
			for _, f := range binlogs {
				pos, ts, lineNum, out, found, err := scanBinlog(ctx, f, targetDate)
				if err != nil {
					return err
				}
				if found {
					if frameshot {
						path, err := SaveFrameShot(out, lineNum, f)
						if err != nil {
							return err
						}
						fmt.Printf("Frameshot salvo em %s\n", path)
					}
					fmt.Printf("Arquivo: %s\nPosição: %d\nData: %s\n", f, pos, ts.Format("2006-01-02"))
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
	cmd.Flags().BoolVar(&frameshot, "frameshot", false, "save 100 lines before and after the match")
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
		return nil, fmt.Errorf("erro executando mysql: %w", err)
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
		return nil, fmt.Errorf("erro processando resposta: %w", err)
	}
	return logs, nil
}

func scanBinlog(ctx context.Context, file string, target time.Time) (int64, time.Time, int, string, bool, error) {
	args := []string{
		"--read-from-remote-server",
		"--host=" + host,
		"--user=" + user,
		"--password=" + password,
		"--port=" + strconv.Itoa(port),
		"--base64-output=DECODE-ROWS",
		"--verbose",
		file,
	}

	out, err := exec.CommandContext(ctx, "./pkg/bin/mysqlbinlog", args...).CombinedOutput()
	if err != nil {
		return 0, time.Time{}, 0, "", false, fmt.Errorf("erro executando mysqlbinlog: %w", err)
	}

	_, pos, ts, lineNum, parseErr := ExtractBinlogPositionFromOutput(string(out), target)
	if parseErr != nil {
		if errors.Is(parseErr, ErrNoEventFound) {
			return 0, time.Time{}, 0, string(out), false, nil
		}
		return 0, time.Time{}, 0, string(out), false, parseErr
	}

	return pos, ts, lineNum, string(out), true, nil
}
