xiter
=====

[![Go Reference](https://pkg.go.dev/badge/deedles.dev/xiter.svg)](https://pkg.go.dev/deedles.dev/xiter)

xiter is a very simple implementation of iterator utility functions for Go's iterators that were introduced in 1.23. Although the module's functionality is compatible with Go 1.23, all of its features should work just fine with any older version of Go that has support for generics (1.18+).

Note that due to the lack of generic type aliases, this package's `Seq` type and the standard library's `iter.Seq` type need to be manually converted between unless using `GOEXPERIMENT=aliastypeparams`. This should hopefully be resolved in Go 1.24. See https://github.com/golang/go/issues/46477.
