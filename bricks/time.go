package bricks

import "time"

type Timespan struct {
	Start time.Time
	End   time.Time
}

func (ts Timespan) Contains(e time.Time) bool {
	return e.After(ts.Start) && e.Before(ts.End)
}
