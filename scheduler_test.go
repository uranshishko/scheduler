package scheduler

import (
	"errors"
	"sync"
	"testing"
	"time"
)

// Test New Scheduler
func TestNewScheduler(t *testing.T) {
	start := time.Now()
	s := New(start)
	if s == nil {
		t.Fatal("Expected instance, got nil")
	}
}

// Test valid scheduling
func TestScheduleValid(t *testing.T) {
	s := New(time.Now())

	var wg sync.WaitGroup
	wg.Add(1) // Ensure we increment before task execution

	handler := func(event Event) error {
		defer wg.Done()
		return nil
	}

	cancel, err := s.Schedule("@every 1s", handler)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	time.Sleep(2 * time.Second)
	cancel()

	wg.Wait() // Ensure the test waits for completion
}

// Test invalid scheduling expression
func TestScheduleInvalidExpression(t *testing.T) {
	s := New(time.Now())

	handler := func(event Event) error {
		return nil
	}

	_, err := s.Schedule("invalid", handler)
	if err == nil {
		t.Fatal("Expected error for invalid expression, got nil")
	}
}

// Test handler returning error stops execution
func TestHandlerErrorStopsExecution(t *testing.T) {
	s := New(time.Now())

	var count int
	handler := func(event Event) error {
		count++
		return errors.New("stop execution")
	}

	cancel, err := s.Schedule("@every 1s", handler)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	time.Sleep(2 * time.Second)
	cancel()

	if count != 1 {
		t.Fatalf("Expected handler to run once, ran %d times", count)
	}
}

// Test custom duration parsing
func TestParseCustomDuration(t *testing.T) {
	s := New(time.Now())
	_, err := s.Schedule("@every 10h20m5s100ms1200ns", func(event Event) error {
		return nil
	})

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}
