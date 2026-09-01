package main

import "testing"

//бенчмарки
/*
➜  1_ValuesVsPointers git:(main) ✗ go test -bench=. -benchmem
goos: darwin
goarch: amd64
pkg: my-module/studyGo/golangroadmap2026/course/01_syntax/1.6_ValuesAndPointers/1_ValuesVsPointers
cpu: Intel(R) Core(TM) i5-4570R CPU @ 2.70GHz
BenchmarkValueReceiver-4     	 4998042	       278.4 ns/op	       0 B/op	       0 allocs/op
BenchmarkPointerReceiver-4   	532183563	         2.169 ns/op	       0 B/op	       0 allocs/op
PASS
ok  	my-module/studyGo/golangroadmap2026/course/01_syntax/1.6_ValuesAndPointers/1_ValuesVsPointers	3.693s
*/

// BenchmarkValueReceiver тестирует вызов метода по значению
func BenchmarkValueReceiver(b *testing.B) {
	instance := BigStruct{}
	b.ResetTimer()
	for i:=0; i<b.N; i++ {
		_ = instance.ValueReceiverMethod()
	}
}

// BenchmarkPointerReceiver тестирует вызов метода по указателю
func BenchmarkPointerReceiver(b *testing.B) {
	instance := BigStruct{}
	b.ResetTimer()
	for i:=0; i<b.N; i++ {
		_ = instance.PointerReceiverMethod()
	}
}