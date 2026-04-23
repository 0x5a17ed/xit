package xit

import (
	"errors"
	"iter"
	"sync/atomic"
)

var (
	ErrNotReady      = errors.New("iterator has not completed yet")
	ErrConcurrentRun = errors.New("iterator already running")
)

type result struct {
	err   error
	ready bool
}

func Perform[T any](
	fn func(yield func(T) bool) error,
) (iter.Seq[T], func() error) {
	var stored atomic.Value
	stored.Store(result{err: nil, ready: false})

	var running atomic.Bool

	seq := func(yield func(T) bool) {
		if !running.CompareAndSwap(false, true) {
			panic(ErrConcurrentRun)
		}
		defer running.Store(false)

		stored.Store(result{err: nil, ready: false})
		err := fn(yield)
		stored.Store(result{err: err, ready: true})
	}

	doneFn := func() error {
		r := stored.Load().(result)
		if !r.ready {
			return ErrNotReady
		}
		return r.err
	}

	return seq, doneFn
}
