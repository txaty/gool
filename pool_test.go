package gool

import (
	"testing"
)

func BenchmarkNormalGoroutine(b *testing.B) {
	for i := 0; i < b.N; i++ {
		go func() {
			for k := 0; k < 100; k++ {
				_ = k
			}
		}()
	}
}

func BenchmarkPool(b *testing.B) {
	p := NewPool(100)
	defer p.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.Submit(func(arg interface{}) {
			for k := 0; k < 100; k++ {
				_ = k
			}
		}, nil)
	}
}
