package nslock

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/avast/retry-go"
)

func init() {
	retry.DefaultDelay = time.Millisecond
	retry.DefaultMaxJitter = 10 * time.Millisecond
}

var errRetry = fmt.Errorf("retry error")

type rwmutex struct {
	sync.Mutex
	isWriteLock bool
	ref         int
}

func newRWMutex() *rwmutex {
	return &rwmutex{}
}

func (m *rwmutex) IncrRef() {
	m.ref++
}

func (m *rwmutex) DecrRef() {
	m.ref--
}

func (m *rwmutex) GetLock(ctx context.Context, timeout time.Duration) (locked bool) {
	retryCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return retry.Do(m.retryLock(retryCtx)) == nil
}

func (m *rwmutex) GetRLock(ctx context.Context, timeout time.Duration) (locked bool) {
	retryCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return retry.Do(m.retryRLock(retryCtx)) == nil
}

func (m *rwmutex) retryLock(ctx context.Context) func() error {
	return func() error {
		select {
		case <-ctx.Done():
			return retry.Unrecoverable(ctx.Err())
		default:
		}

		locked := m.lock()
		if locked {
			return nil
		}
		return errRetry
	}
}

func (m *rwmutex) retryRLock(ctx context.Context) func() error {
	return func() error {
		select {
		case <-ctx.Done():
			return retry.Unrecoverable(ctx.Err())
		default:
		}

		locked := m.rlock()
		if locked {
			return nil
		}
		return errRetry
	}
}

func (m *rwmutex) lock() (locked bool) {
	m.Mutex.Lock()
	if m.ref == 0 && !m.isWriteLock {
		m.ref = 1
		m.isWriteLock = true
		locked = true
	}
	m.Mutex.Unlock()
	return locked
}

func (m *rwmutex) rlock() (locked bool) {
	m.Mutex.Lock()
	if !m.isWriteLock {
		m.IncrRef()
		locked = true
	}
	m.Mutex.Unlock()
	return locked
}

func (m *rwmutex) unlock() (unlocked bool) {
	m.Mutex.Lock()
	if m.isWriteLock && m.ref == 1 {
		m.ref = 0
		m.isWriteLock = false
		unlocked = true
	}
	m.Mutex.Unlock()
	return unlocked
}

func (m *rwmutex) runlock() (unlocked bool) {
	m.Mutex.Lock()
	if !m.isWriteLock {
		if m.ref > 0 {
			m.DecrRef()
			unlocked = true
		}
	}
	m.Mutex.Unlock()
	return unlocked
}
