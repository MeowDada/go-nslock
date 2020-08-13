package nslock

import (
	"time"
)

const timeoutErr = "get lock operation timeout"

// GetLockErr implements error interface. And it appends
// the information of the get lock request.
type GetLockErr struct {
	start   time.Time
	end     time.Time
	timeout time.Duration
	keys    []string
}

func newGetLockErr(start time.Time, keys []string, timeout time.Duration) GetLockErr {
	return GetLockErr{
		start:   start,
		end:     time.Now(),
		timeout: timeout,
		keys:    keys,
	}
}

// Error implements error interface.
func (e GetLockErr) Error() string {
	return timeoutErr
}
