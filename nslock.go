package nslock

import (
	"context"
	"fmt"
	"time"

	cmap "github.com/orcaman/concurrent-map"
)

// RWLocker is a locker interface that provides read-write lock capabilities.
type RWLocker interface {
	LockFn(timeout time.Duration, fn func() error) error
	RLockFn(timeout time.Duration, fn func() error) error
	Lock(timeout time.Duration) error
	Unlock()
	RLock(timeout time.Duration) error
	RUnlock()
}

// Map is a namespace locks manager.
type Map struct {
	shardMap cmap.ConcurrentMap
}

// NewMap returns an initialized map instance.
func NewMap() *Map {
	return &Map{
		shardMap: cmap.New(),
	}
}

// New returns a namespace lock instance managed by this map.
func (m *Map) New(ctx context.Context, keys ...string) RWLocker {
	return &Instance{
		ctx:  ctx,
		keys: keys,
		m:    m,
	}
}

func (m *Map) lock(ctx context.Context, key string, timeout time.Duration) (locked bool) {
	mu := m.shardMap.Upsert(key, newNSLock(1), upsertNSLock).(*nsLock)
	locked = mu.GetLock(ctx, timeout)
	if !locked {
		m.shardMap.RemoveCb(key, rollbackNSLock)
	}
	return locked
}

func (m *Map) rlock(ctx context.Context, key string, timeout time.Duration) (locked bool) {
	mu := m.shardMap.Upsert(key, newNSLock(1), upsertNSLock).(*nsLock)
	locked = mu.GetRLock(ctx, timeout)
	if !locked {
		m.shardMap.RemoveCb(key, rollbackNSLock)
	}
	return locked
}

func (m *Map) unlock(key string) {
	m.shardMap.RemoveCb(key, unlockNSLock)
}

func (m *Map) runlock(key string) {
	m.shardMap.RemoveCb(key, runlockNSLock)
}

type nsLock struct {
	ref int
	*rwmutex
}

func newNSLock(ref int) *nsLock {
	return &nsLock{
		ref:     ref,
		rwmutex: newRWMutex(),
	}
}

func (n *nsLock) IncrRef() {
	n.ref++
}

func (n *nsLock) DecrRef() {
	n.ref--
}

func (n *nsLock) GetLock(ctx context.Context, timeout time.Duration) bool {
	return n.rwmutex.GetLock(ctx, timeout)
}

func (n *nsLock) GetRLock(ctx context.Context, timeout time.Duration) bool {
	return n.rwmutex.GetRLock(ctx, timeout)
}

func upsertNSLock(exist bool, valueInMap, newValue interface{}) interface{} {
	if !exist {
		return newValue
	}
	mu := valueInMap.(*nsLock)
	mu.IncrRef()
	return mu
}

func rollbackNSLock(key string, v interface{}, exist bool) bool {
	// Decrement the reference count of the nslock with namepsace as key. Then, remove the
	// namespace lock if and only if it's reference count is zero.
	if exist {
		nslk := v.(*nsLock)
		nslk.DecrRef()
		if nslk.ref < 0 {
			panic(fmt.Errorf("nsLock (key=%s) with negative reference count", key))
		}
		if nslk.ref == 0 {
			return true
		}
	}
	return false
}

func unlockNSLock(key string, v interface{}, exist bool) bool {
	// Decrement the reference count of the nslock with namepsace as key. Then, unlock the
	// nslock and remove this nslock from the map if and only if it's reference count is zero.
	if exist {
		nslk := v.(*nsLock)
		nslk.unlock()
		nslk.DecrRef()
		if nslk.ref < 0 {
			panic(fmt.Errorf("nsLock (key=%s) with negative reference count", key))
		}
		if nslk.ref == 0 {
			return true
		}
	}
	return false
}

func runlockNSLock(key string, v interface{}, exist bool) bool {
	// Decrement the reference count of the nslock with namepsace as key. Then, unlock the
	// nslock and remove this nslock from the map if and only if it's reference count is zero.
	if exist {
		nslk := v.(*nsLock)
		nslk.runlock()
		nslk.DecrRef()
		if nslk.ref < 0 {
			panic(fmt.Errorf("nsLock (key=%s) with negative reference count", key))
		}
		if nslk.ref == 0 {
			return true
		}
	}
	return false
}
