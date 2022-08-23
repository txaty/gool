package gool

import (
	"testing"
)

func BenchmarkNormalGoroutine(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		finish := make(chan struct{})
		go func() {
			for k := 0; k < 100; k++ {
				_ = k
			}
			finish <- struct{}{}
		}()
		<-finish
	}
}

func BenchmarkPool(b *testing.B) {
	p := NewPool(4, 100)
	defer p.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.Submit(func(interface{}) interface{} {
			for k := 0; k < 100; k++ {
				_ = k
			}
			return nil
		}, nil)
	}
}

func TestPool_Submit(t *testing.T) {
	type args struct {
		handler func(interface{}) interface{}
		args    interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test",
			args: args{
				handler: func(arg interface{}) interface{} {
					for k := 0; k < 100; k++ {
						_ = k
					}
					return nil
				},
				args: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPool(10, 100)
			p.Submit(tt.args.handler, tt.args.args)
		})
	}
}

func TestPool_Map(t *testing.T) {
	type args struct {
		handler func(interface{}) interface{}
		args    []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test",
			args: args{
				handler: func(arg interface{}) interface{} {
					for k := 0; k < 100; k++ {
						_ = k
					}
					return nil
				},
				args: []interface{}{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPool(10, 100)
			defer p.Close()
			p.Map(tt.args.handler, tt.args.args)
		})
	}
}
