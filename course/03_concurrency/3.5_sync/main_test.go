package main

import (
	"testing"
)

func BenchmarkMutexSerial(b *testing.B) {
	counter := &MutexCounter{}
	for i := 0; i < b.N; i++ {
		counter.Increment()
	}
}

func BenchmarkMutexParallel(b *testing.B) {
	counter := &MutexCounter{}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			counter.Increment()
		}
	})
}

func BenchmarkAtomicSerial(b *testing.B) {
	counter := &AtomicCounter{}
	for i := 0; i < b.N; i++ {
		counter.Increment()
	}
}

func BenchmarkAtomicParallel(b *testing.B) {
	counter := &AtomicCounter{}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			counter.Increment()
		}
	})
}