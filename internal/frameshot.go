package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SaveFrameShot writes 100 lines before and after the given line number (1-based)
// from the provided output and returns the created file path.
func SaveFrameShot(output string, line int, binlog string) (string, error) {
	lines := strings.Split(strings.ReplaceAll(output, "\r\n", "\n"), "\n")
	idx := line - 1
	if idx < 0 || idx >= len(lines) {
		return "", fmt.Errorf("line out of range")
	}
	start := idx - 100
	if start < 0 {
		start = 0
	}
	end := idx + 100
	if end >= len(lines) {
		end = len(lines) - 1
	}

	frame := append([]string{}, lines[start:idx]...)
	frame = append(frame, lines[idx+1:end+1]...)

	name := fmt.Sprintf("binlog-frameshot-%s-%d.log", filepath.Base(binlog), line)
	if err := os.WriteFile(name, []byte(strings.Join(frame, "\n")), 0o644); err != nil {
		return "", err
	}
	return name, nil
}
