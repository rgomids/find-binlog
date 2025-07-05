package cmd

import (
	"fmt"
	"os"
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
			if _, err := time.Parse("2006-01-02", dateStr); err != nil {
				return fmt.Errorf("invalid date format (expected YYYY-MM-DD)")
			}

			fmt.Printf("host: %s\n", host)
			fmt.Printf("port: %d\n", port)
			fmt.Printf("user: %s\n", user)
			fmt.Printf("password: %s\n", password)
			fmt.Printf("date: %s\n", dateStr)
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

func Execute() {
	if err := NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
