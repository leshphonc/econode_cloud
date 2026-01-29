package timeutil

import "time"

func NormalizeReportedAt(ms *int64, now time.Time, maxFuture, maxPast time.Duration) *time.Time {
	if ms == nil {
		return nil
	}
	t := time.UnixMilli(*ms)

	if maxFuture > 0 && t.After(now.Add(maxFuture)) {
		return nil
	}
	if maxPast > 0 && t.Before(now.Add(-maxPast)) {
		return nil
	}
	return &t
}
