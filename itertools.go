// Package itertools implements common utilities to work with Go 1.23's iterators.
package itertools

import (
	"iter"
)

// FromSlice returns an iterator yielding all the values from vs.
func FromSlice[V any](vs []V) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range vs {
			if !yield(v) {
				return
			}
		}
	}
}

// Map returns an iterator that will yield values from seq after transforming them using f.
func Map[V any, W any](seq iter.Seq[V], f func(V) W) iter.Seq[W] {
	return func(yield func(W) bool) {
		for v := range seq {
			if !yield(f(v)) {
				return
			}
		}
	}
}

// Filter returns an iterator that will yield values from seq only if they pass p.
func Filter[V any](seq iter.Seq[V], p func(V) bool) iter.Seq[V] {
	return func(yield func(V) bool) {
		for v := range seq {
			if p(v) {
				if !yield(v) {
					return
				}
			}
		}
	}
}
