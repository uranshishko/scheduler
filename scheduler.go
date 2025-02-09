package scheduler

import (
	"errors"
	"regexp"
	"strings"
	"sync/atomic"
	"time"
)

// Regular expression to match predefined and custom scheduling expressions.
var rgxp = regexp.MustCompile(`(?P<predefined>@(yearly|monthly|weekly|daily|hourly))|(?P<custom>@every (\d+(ns|us|µs|ms|s|m|h))+)`)

// Scheduler represents a scheduling system that starts from a given time.
type Scheduler struct {
	start time.Time
}

// New creates a new Scheduler instance with a specified start time.
func New(start time.Time) *Scheduler {
	return &Scheduler{start}
}

// Handler defines a function signature that processes scheduled events.
type Handler func(event Event) error

// Event represents an occurrence of a scheduled task.
type Event struct {
	Time time.Time
}

// Schedule sets up a scheduled task based on the given expression and handler function.
// It returns a cancel function to stop the schedule, or an error if the expression is invalid.
func (s *Scheduler) Schedule(expr string, handler Handler) (func(), error) {
	// Parse the scheduling expression.
	ce, err := parse(expr)
	if err != nil {
		return nil, err
	}

	// Determine the next occurrence of the scheduled event.
	nextOccurrence := s.start
	now := time.Now()
	for nextOccurrence.Before(now) || nextOccurrence.Equal(now) {
		nextOccurrence = nextOccurrence.Add(ce.Frequency)
	}

	// Create a ticker that checks at the interval of the frequency.
	ticker := time.NewTicker(ce.Frequency)
	done := make(chan struct{})

	var closed atomic.Bool

	// Goroutine to handle scheduled execution.
	go func() {
		for {
			select {
			case <-done:
				// Cleanup and exit the goroutine.
				closed.Store(true)
				return
			case t := <-ticker.C:
				if t.Before(nextOccurrence) {
					continue
				}

				event := Event{Time: t}
				if err := handler(event); err != nil {
					ticker.Stop()
					close(done) // Close the done channel when done.
					break
				}

				// Update the next occurrence.
				nextOccurrence = ce.NextOccurrence(t)
			}
		}
	}()

	// Cancel function to stop the scheduled execution.
	cancel := func() {
		// Check if the goroutine is closed.
		if closed.Load() {
			return
		}

		ticker.Stop()
		close(done) // Close the done channel to stop the goroutine.
	}

	return cancel, nil
}

// parse analyzes the scheduling expression and returns a corresponding Schedule.
func parse(expr string) (*Schedule, error) {
	// Match the expression against the regex.
	matches := rgxp.FindStringSubmatch(expr)
	if matches == nil {
		return nil, errors.New("invalid expression")
	}

	// Map regex capture groups to their names.
	mapped := make(map[string]string)
	for i, name := range rgxp.SubexpNames() {
		if i != 0 && name != "" {
			mapped[name] = matches[i]
		}
	}

	var freq time.Duration

	// Handle predefined scheduling intervals.
	if predefined, ok := mapped["predefined"]; ok && predefined != "" {
		switch predefined {
		case "@yearly":
			freq = time.Hour * 24 * 365
		case "@monthly":
			freq = time.Hour * 24 * 30
		case "@weekly":
			freq = time.Hour * 24 * 7
		case "@daily":
			freq = time.Hour * 24
		case "@hourly":
			freq = time.Hour
		}
	}

	// Handle custom time intervals.
	if custom, ok := mapped["custom"]; ok && custom != "" {
		custom = strings.Replace(custom, "@every ", "", 1)
		var err error
		freq, err = time.ParseDuration(custom)
		if err != nil {
			return nil, err
		}
	}

	// Ensure a valid frequency was determined.
	if freq == 0 {
		return nil, errors.New("invalid expression")
	}

	return &Schedule{freq}, nil
}

// Schedule defines a recurring frequency for event execution.
type Schedule struct {
	Frequency time.Duration
}

// NextOccurrence calculates the next scheduled execution time based on the previous one.
func (s *Schedule) NextOccurrence(prev time.Time) (next time.Time) {
	next = prev.Add(s.Frequency)
	return
}
