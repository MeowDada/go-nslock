package nslock

import (
	"context"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestMap(t *testing.T) {
	testcases := []struct {
		description  string
		scenario     func() bool
		expectResult bool
	}{
		{
			"valid lock, it will prevent all the other locks to acquire it with the namespace",
			func() bool {
				ns := "123"
				m := NewMap()
				l1 := m.New(context.Background(), ns)
				if err := l1.Lock(time.Second); err != nil {
					return false
				}

				final := false
				var wg sync.WaitGroup
				for i := 0; i < 100; i++ {
					wg.Add(1)
					go func() {
						defer wg.Done()
						l2 := m.New(context.Background(), ns)
						err := l2.Lock(time.Microsecond)
						if err == nil {
							final = true
						}
					}()
				}
				wg.Wait()
				return final == false
			},
			true,
		},
		{
			"valid lock, there should no conflicts when locking in different namespaces",
			func() bool {
				m := NewMap()
				for i := 0; i < 100; i++ {
					key := strconv.Itoa(i)
					lock := m.New(context.Background(), key)
					err := lock.Lock(time.Millisecond)
					if err != nil {
						return false
					}
				}
				return true
			},
			true,
		},
		{
			"valid lock and unlock",
			func() bool {
				m := NewMap()
				for i := 0; i < 100; i++ {
					key := strconv.Itoa(i)
					lock := m.New(context.Background(), key)
					err := lock.Lock(10 * time.Millisecond)
					if err != nil {
						return false
					}
					lock.Unlock()
				}
				return true
			},
			true,
		},
		{
			"valid lock, expect lock blocks rlock method",
			func() bool {
				ns := "123"
				m := NewMap()
				lock := m.New(context.Background(), ns)
				if err := lock.Lock(time.Second); err != nil {
					return false
				}
				for i := 0; i < 100; i++ {
					rlock := m.New(context.Background(), ns)
					err := rlock.RLock(time.Millisecond)
					if err == nil {
						return false
					}
				}
				return true
			},
			true,
		},
		{
			"valid rlock",
			func() bool {
				ns := "123"
				m := NewMap()
				for i := 0; i < 100; i++ {
					lock := m.New(context.Background(), ns)
					err := lock.RLock(time.Millisecond)
					if err != nil {
						return false
					}
				}
				return true
			},
			true,
		},
		{
			"valid rlock, expect rlock blocks lock method",
			func() bool {
				ns := "123"
				m := NewMap()
				rlock := m.New(context.Background(), ns)
				if err := rlock.RLock(time.Second); err != nil {
					return false
				}
				lock := m.New(context.Background(), ns)
				return lock.Lock(time.Second) != nil
			},
			true,
		},
		{
			"valid rlock with unlock",
			func() bool {
				ns := "123"
				m := NewMap()
				for i := 0; i < 100; i++ {
					rlock := m.New(context.Background(), ns)
					if err := rlock.RLock(time.Second); err != nil {
						return false
					}
					rlock.RUnlock()
				}
				lock := m.New(context.Background(), ns)
				return lock.Lock(time.Millisecond) == nil
			},
			true,
		},
		{
			"valid LockFn",
			func() bool {
				ns := []string{"123", "456", "789"}
				m := NewMap()
				for i := 0; i < 100; i++ {
					lock := m.New(context.Background(), ns...)
					if err := lock.LockFn(time.Second, func() error { return nil }); err != nil {
						return false
					}
				}
				return true
			},
			true,
		},
		{
			"valid RLockFn",
			func() bool {
				ns := []string{"123", "456", "789"}
				m := NewMap()
				for i := 0; i < 100; i++ {
					lock := m.New(context.Background(), ns...)
					if err := lock.RLockFn(time.Second, func() error { return nil }); err != nil {
						return false
					}
				}
				return true
			},
			true,
		},
	}

	for _, tc := range testcases {
		result := tc.scenario()
		if result != tc.expectResult {
			t.Errorf("%s: expect %v, but get %v", tc.description, tc.expectResult, result)
		}
	}
}

func TestGetLockErr(t *testing.T) {
	err := newGetLockErr(time.Now(), nil, time.Second)
	if err.Error() != timeoutErr {
		t.Errorf("expect %v, get %v", timeoutErr, err.Error())
	}
}
