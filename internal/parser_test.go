package internal

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

func TestExtractBinlogPositionFromOutput(t *testing.T) {
	data, err := os.ReadFile("testdata/sample.binlog")
	if err != nil {
		t.Fatalf("failed to read sample: %v", err)
	}
	target := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	file, pos, ts, lineNum, err := ExtractBinlogPositionFromOutput(string(data), target)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if file != "binlog.000001" {
		t.Errorf("expected file binlog.000001, got %s", file)
	}
	if pos != 128 {
		t.Errorf("expected pos 128, got %d", pos)
	}
	expected := target
	if !ts.Equal(expected) {
		t.Errorf("expected ts %v, got %v", expected, ts)
	}
	if lineNum != 5 {
		t.Errorf("expected line 5, got %d", lineNum)
	}
}

func TestExtractBinlogPositionFromOutputNotFound(t *testing.T) {
	data, err := os.ReadFile("testdata/sample_no_event.binlog")
	if err != nil {
		t.Fatalf("failed to read sample: %v", err)
	}
	target := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)
	_, _, _, _, err = ExtractBinlogPositionFromOutput(string(data), target)
	if !errors.Is(err, ErrNoEventFound) {
		t.Fatalf("expected ErrNoEventFound, got %v", err)
	}
}

func TestExtractBinlogPositionFromOutputSameTimestamp(t *testing.T) {
	data, err := os.ReadFile("testdata/sample_same_ts.binlog")
	if err != nil {
		t.Fatalf("failed to read sample: %v", err)
	}
	target := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	file, pos, ts, lineNum, err := ExtractBinlogPositionFromOutput(string(data), target)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if file != "binlog.000002" {
		t.Errorf("expected file binlog.000002, got %s", file)
	}
	if pos != 4 {
		t.Errorf("expected pos 4, got %d", pos)
	}
	if !ts.Equal(target) {
		t.Errorf("expected ts %v, got %v", target, ts)
	}
	if lineNum != 3 {
		t.Errorf("expected line 3, got %d", lineNum)
	}
}

func TestExtractBinlogPositionFromOriginalCommitTimestamp(t *testing.T) {
	data, err := os.ReadFile("testdata/sample_original_commit.binlog")
	if err != nil {
		t.Fatalf("failed to read sample: %v", err)
	}
	target := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	file, pos, ts, lineNum, err := ExtractBinlogPositionFromOutput(string(data), target)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if file != "binlog.000003" {
		t.Errorf("expected file binlog.000003, got %s", file)
	}
	if pos != 128 {
		t.Errorf("expected pos 128, got %d", pos)
	}
	expected := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
	if !ts.Equal(expected) {
		t.Errorf("expected ts %v, got %v", expected, ts)
	}
	if lineNum != 5 {
		t.Errorf("expected line 5, got %d", lineNum)
	}
}

func TestExtractBinlogPositionOriginalCommitTimestampMultiple(t *testing.T) {
	data, err := os.ReadFile("testdata/sample_original_commit_multi.binlog")
	if err != nil {
		t.Fatalf("failed to read sample: %v", err)
	}
	target := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	file, pos, ts, lineNum2, err := ExtractBinlogPositionFromOutput(string(data), target)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if file != "binlog.000004" {
		t.Errorf("expected file binlog.000004, got %s", file)
	}
	if pos != 128 {
		t.Errorf("expected pos 128, got %d", pos)
	}
	if !ts.Equal(target) {
		t.Errorf("expected ts %v, got %v", target, ts)
	}
	if lineNum2 != 5 {
		t.Errorf("expected line 5, got %d", lineNum2)
	}
}

func TestSaveFrameShot(t *testing.T) {
	var lines []string
	for i := 1; i <= 300; i++ {
		lines = append(lines, fmt.Sprintf("line %d", i))
	}
	lines[0] = "mysqlbinlog: [Note] Start processing log events from binlog.000010, position 4."
	lines[148] = "# at 999"
	lines[149] = "### SET TIMESTAMP=1704067200"

	output := strings.Join(lines, "\n")
	target := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	binlog, pos, ts, lineNum, err := ExtractBinlogPositionFromOutput(output, target)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lineNum != 150 {
		t.Fatalf("expected line 150, got %d", lineNum)
	}
	if pos != 999 {
		t.Fatalf("expected pos 999, got %d", pos)
	}
	if !ts.Equal(target) {
		t.Fatalf("expected ts %v, got %v", target, ts)
	}

	path, err := SaveFrameShot(output, lineNum, binlog)
	if err != nil {
		t.Fatalf("failed to save frameshot: %v", err)
	}
	defer os.Remove(path)

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read frameshot file: %v", err)
	}
	saved := strings.Split(strings.TrimSuffix(string(data), "\n"), "\n")
	if len(saved) != 200 {
		t.Fatalf("expected 200 lines, got %d", len(saved))
	}
	if saved[0] != lines[49] {
		t.Errorf("first line mismatch: %s", saved[0])
	}
	if saved[len(saved)-1] != lines[249] {
		t.Errorf("last line mismatch: %s", saved[len(saved)-1])
	}
}
