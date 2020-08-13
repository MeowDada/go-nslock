go-nslock
=====
[![Go Report Card](https://goreportcard.com/badge/github.com/meowdada/go-nslock)](https://goreportcard.com/report/github.com/meowdada/go-nslock)
[![codecov](https://codecov.io/gh/MeowDada/go-nslock/branch/master/graph/badge.svg)](https://codecov.io/gh/MeowDada/go-nslock)
[![Build Status](https://travis-ci.org/MeowDada/go-nslock.svg?branch=master)](https://travis-ci.org/MeowDada/go-nslock)
[![LICENSE](https://img.shields.io/github/license/meowdada/go-nslock)](https://img.shields.io/github/license/meowdada/go-nslock)

go-nslock is a pure Go implementation of namespace lock.

# Install
```bash
go get github.com/meowdada/nslock
```

# What is namespace lock ?
Namespace lock is a kind of lock which only locks on a specific namespace, which usually is a key consisting of a string or an array of strings.

The main purpose of the namespace lock is to maintenance the consistency of a system by locking resources in proper way.

For example, say we have an application that allow user to apply CRUD operations on objects which are stored in our backend (database, cloud or something else...). Once there are multiple operations, such like ADD, UPDATE or DELETE on the same object at the sametime, how can we determine the eventually result would be. By using namespace lock, we can make these operations atomic without blocking other objects with different namespace. This could improve the performance significantly compared to using a single lock.

# Usage
In brief, you have to create a `nslock.Map` instance to manages all the namespace locks. In most cases, create one would be enough.

Then, whenever you want to acquire a namespace lock, you'll need to create a lock instance and do the `Lock` or `RLock` depends on your need. And remember to reclaim these locks by `Unlock` or `RUnlock` when you're no longer need it or it might be blocked forerver.

```go
// Creates a namespace lock manager.
m := NewMap()

// Creates a new lock instance with given namespace.
ctx, namespace := context.Background(), "myNamespace"
ins := m.New(ctx, namespace)

// Defines timeout that denotes how long can this thread blocking wait for fetching this locker.
timeout := time.Second

// Peform try locking operation.
if err := ins.GetLock(timeout); err != nil {
    // Your logic for getting lock failed...
    return ...
}
defer ins.Unlock()
```
Note that the timeout value of any GetLock operations should be reasonable one. If the value is too small such as microsecond, it might lead to impossible to fetch a lock because the creation time of internal data structure when fetching lock always consume more than a microsecond. The minimum value depends on your hardware, but i think millisecond would be a minimum safe choice.
# Contributing
Any contributions are welcome.

# References
Inspired by awesome `minio` subprojects.
* [lsync](https://github.com/minio/lsync)
* [dysnc](https://github.com/minio/dsync)

And we use these two packages internally for sharding and retry capabilities.
* [concurrent-map](https://github.com/orcaman/concurrent-map)
* [retry-go](https://github.com/avast/go-retry)