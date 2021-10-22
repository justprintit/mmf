package library

import (
	"fmt"
	"time"
)

func (w *Worker) canRefresh(s string, args ...interface{}) bool {
	if len(args) > 0 {
		s = fmt.Sprintf(s, args...)
	}

	w.mu.Lock()
	defer w.mu.Unlock()

	if limit, ok := w.refresh[s]; ok {
		if !limit.After(time.Now()) {
			// not now
			return false
		}
	}

	// go ahead
	return true
}

func (w *Worker) nextRefresh(d time.Duration, s string, args ...interface{}) time.Time {
	if len(args) > 0 {
		s = fmt.Sprintf(s, args...)
	}

	limit := time.Now().Add(d)

	w.mu.Lock()
	defer w.mu.Unlock()

	w.refresh[s] = limit
	return limit
}
