package bricks

import "time"

// Timespan is a half-open instant range in the sense that Contains is exclusive of both ends:
// an instant equal to Start or End is outside the span.
type Timespan struct {
	Start time.Time
	End   time.Time
}

// Contains reports whether e is strictly after Start and strictly before End.
func (ts Timespan) Contains(e time.Time) bool {
	return e.After(ts.Start) && e.Before(ts.End)
}
