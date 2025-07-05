package internal

import (
	"errors"
	"os"
	"testing"
	"time"
)

func TestExtractBinlogPositionFromOutput(t *testing.T) {
	data, err := os.ReadFile("testdata/sample.binlog")
	if err != nil {
		t.Fatalf("failed to read sample: %v", err)
	}
	target := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	file, pos, ts, err := ExtractBinlogPositionFromOutput(string(data), target)
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
}

func TestExtractBinlogPositionFromOutputNotFound(t *testing.T) {
	data, err := os.ReadFile("testdata/sample_no_event.binlog")
	if err != nil {
		t.Fatalf("failed to read sample: %v", err)
	}
	target := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)
	_, _, _, err = ExtractBinlogPositionFromOutput(string(data), target)
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
	file, pos, ts, err := ExtractBinlogPositionFromOutput(string(data), target)
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
}
