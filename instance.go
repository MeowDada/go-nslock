package nslock

import (
	"context"
	"time"
)

// Instance is a lock instance which is able to lock/unlock a namepsace.
type Instance struct {
	ctx  context.Context
	keys []string
	m    *Map
}

// Lock locks all underlying namespaces with a write lock.
func (ins *Instance) Lock(timeout time.Duration) error {
	start := time.Now()
	var success []int
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
	start := time.Now()
	var success []int
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
