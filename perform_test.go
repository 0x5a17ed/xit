package xit

import (
	"errors"
	"slices"
	"testing"
)

func TestPerform(t *testing.T) {
	errBoom := errors.New("boom")

	type testCase struct {
		name    string
		fn      func(yield func(int) bool) error
		wantSeq []int
		wantErr error
	}

	tests := []testCase{
		{
			name: "yields values and returns nil",
			fn: func(yield func(int) bool) error {
				yield(1)
				yield(2)
				yield(3)
				return nil
			},
			wantSeq: []int{1, 2, 3},
			wantErr: nil,
		},
		{
			name: "yields nothing and returns nil",
			fn: func(yield func(int) bool) error {
				return nil
			},
			wantSeq: []int{},
			wantErr: nil,
		},
		{
			name: "yields values and returns error",
			fn: func(yield func(int) bool) error {
				yield(1)
				yield(2)
				return errBoom
			},
			wantSeq: []int{1, 2},
			wantErr: errBoom,
		},
		{
			name: "yields nothing and returns error",
			fn: func(yield func(int) bool) error {
				return errBoom
			},
			wantSeq: []int{},
			wantErr: errBoom,
		},
		{
			name: "stops early when yield returns false",
			fn: func(yield func(int) bool) error {
				if !yield(1) {
					return nil
				}
				yield(2)
				yield(3)
				return nil
			},
			wantSeq: []int{1},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seq, done := Perform(tt.fn)

			if err := done(); !errors.Is(err, ErrNotReady) {
				t.Fatalf("before iteration: got %v, want ErrNotReady", err)
			}

			var got []int
			if tt.name == "stops early when yield returns false" {
				seq(func(v int) bool {
					got = append(got, v)
					return false
				})
			} else {
				got = slices.Collect(seq)
			}

			if !slices.Equal(got, tt.wantSeq) {
				t.Fatalf("seq: got %v, want %v", got, tt.wantSeq)
			}

			if err := done(); !errors.Is(err, tt.wantErr) {
				t.Fatalf("err: got %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestPerformReusable(t *testing.T) {
	var calls int
	seq, getErr := Perform(func(yield func(int) bool) error {
		calls++
		yield(calls)
		return nil
	})

	slices.Collect(seq)
	if err := getErr(); err != nil {
		t.Fatal(err)
	}

	slices.Collect(seq)
	if err := getErr(); err != nil {
		t.Fatal(err)
	}

	if calls != 2 {
		t.Fatalf("calls: got %d, want 2", calls)
	}
}

func TestPerformConcurrentRunPanics(t *testing.T) {
	var seq func(func(int) bool)
	seq, _ = Perform(func(yield func(int) bool) error {
		defer func() {
			r := recover()
			if r == nil {
				t.Fatal("expected panic on concurrent run")
			}
			if !errors.Is(r.(error), ErrConcurrentRun) {
				t.Fatalf("got %v, want ErrConcurrentRun", r)
			}
		}()
		seq(func(int) bool { return true })
		return nil
	})

	slices.Collect(seq)
}
