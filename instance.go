package nslock

import (
	"context"
	"time"
)

// Instance is a lock instance which is able to lock/unlock multiple namespaces at once.
//
// Note that you can only create a valid instance by invoking Map.New(ctx, keys...).
type Instance struct {
	ctx  context.Context
	keys []string
	m    *Map
}

// LockFn locks all underlying namespaces with a write lock and automatically
// unlock this locker when LockFn returned. It will keep contesting for the
// successful lock until it succeed or being timeout.
func (ins *Instance) LockFn(timeout time.Duration, fn func() error) error {
	// Records the starting timestamp to determine if it should
	// throw timeout error.
	start := time.Now()

	// This integer array stores the index of the locked namespace.
	// When this lock instance failed to lock all the namespaces, it
	// have to do some rollback for freeing previous locked namespaces.
	var success []int

	// Iterating all namespaces under this locker instance and try lock
	// them all at once. If it failed to lock anyone of the namespace. The
	// lock operation will be treated as failure.
	for i, key := range ins.keys {
		if !ins.m.lock(ins.ctx, key, timeout) {
			for _, j := range success {
				ins.m.unlock(ins.keys[j])
			}
			return newGetLockErr(start, ins.keys, timeout)
		}
		success = append(success, i)
	}

	// Defer unlocks all locked namespaces by this locker.
	defer ins.Unlock()

	// Return the execution result of the callback function.
	return fn()
}

// RLockFn locks all underlying namespaces with a read lock and automatically
// unlock this locker when LockFn returned. It will keep contesting for the
// successful lock until it succeed or being timeout.
func (ins *Instance) RLockFn(timeout time.Duration, fn func() error) error {
	// Records the starting timestamp to determine if it should
	// throw timeout error.
	start := time.Now()

	// This integer array stores the index of the locked namespace.
	// When this lock instance failed to lock all the namespaces, it
	// have to do some rollback for freeing previous locked namespaces.
	var success []int

	// Iterating all namespaces under this locker instance and try lock
	// them all at once. If it failed to lock anyone of the namespace. The
	// lock operation will be treated as failure.
	for i, key := range ins.keys {
		if !ins.m.rlock(ins.ctx, key, timeout) {
			for _, j := range success {
				ins.m.unlock(ins.keys[j])
			}
			return newGetLockErr(start, ins.keys, timeout)
		}
		success = append(success, i)
	}

	// Defer unlocks all locked namespaces by this locker.
	defer ins.RUnlock()

	// Return the execution result of the callback function.
	return fn()
}

// Lock locks all underlying namespaces with a write lock.
func (ins *Instance) Lock(timeout time.Duration) error {
	// Records the starting timestamp to determine if it should
	// throw timeout error.
	start := time.Now()

	// This integer array stores the index of the locked namespace.
	// When this lock instance failed to lock all the namespaces, it
	// have to do some rollback for freeing previous locked namespaces.
	var success []int

	// Iterating all namespaces under this locker instance and try lock
	// them all at once. If it failed to lock anyone of the namespace. The
	// lock operation will be treated as failure.
	for i, key := range ins.keys {
		if !ins.m.lock(ins.ctx, key, timeout) {
			for _, j := range success {
				ins.m.unlock(ins.keys[j])
			}
			return newGetLockErr(start, ins.keys, timeout)
		}
		success = append(success, i)
	}
	return nil
}

// Unlock unlocks the write lock of underlying namespace.
func (ins *Instance) Unlock() {
	for _, key := range ins.keys {
		ins.m.unlock(key)
	}
}

// RLock locks underlying namespace with a read lock.
func (ins *Instance) RLock(timeout time.Duration) error {
	// Records the starting timestamp to determine if it should
	// throw timeout error.
	start := time.Now()

	// This integer array stores the index of the locked namespace.
	// When this lock instance failed to lock all the namespaces, it
	// have to do some rollback for freeing previous locked namespaces.
	var success []int

	// Iterating all namespaces under this locker instance and try lock
	// them all at once. If it failed to lock anyone of the namespace. The
	// lock operation will be treated as failure.
	for i, key := range ins.keys {
		if !ins.m.rlock(ins.ctx, key, timeout) {
			for _, j := range success {
				ins.m.unlock(ins.keys[j])
			}
			return newGetLockErr(start, ins.keys, timeout)
		}
		success = append(success, i)
	}
	return nil
}

// RUnlock unlocks the read lock of underlying namespace.
func (ins *Instance) RUnlock() {
	for _, key := range ins.keys {
		ins.m.runlock(key)
	}
}
