// MIT License
//
// Copyright (c) 2022 Tommy TIAN
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package gool

import (
	"runtime"
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
	p := NewPool[any, any](4, 100)
	defer p.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.Submit(func(any) any {
			for k := 0; k < 100; k++ {
				_ = k
			}
			return nil
		}, nil)
	}
}

func TestPool_Submit(t *testing.T) {
	type args struct {
		handler func(any) any
		args    any
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test",
			args: args{
				handler: func(arg any) any {
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
			p := NewPool[any, any](10, 100)
			p.Submit(tt.args.handler, tt.args.args)
		})
	}
}

func TestPool_AsyncSubmit(t *testing.T) {
	type args struct {
		handler func(any) any
		args    any
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test",
			args: args{
				handler: func(arg any) any {
					for k := 0; k < 100; k++ {
						_ = k
					}
					return nil
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPool[any, any](10, 100)
			p.AsyncSubmit(tt.args.handler, tt.args.args)
		})
	}
}
func TestPool_Map(t *testing.T) {
	type args struct {
		handler func(any) any
		args    []any
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test",
			args: args{
				handler: func(arg any) any {
					for k := 0; k < 100; k++ {
						_ = k
					}
					return nil
				},
				args: []any{nil, nil, nil, nil, nil, nil, nil, nil, nil, nil},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPool[any, any](10, 100)
			defer p.Close()
			p.Map(tt.args.handler, tt.args.args)
		})
	}
}

func TestNewPool(t *testing.T) {
	type args struct {
		numWorkers int
		cap        int
	}
	type testCase[A any, R any] struct {
		name string
		args args
		want *Pool[A, R]
	}
	tests := []testCase[any, any]{
		{
			name: "test_num_workers_0",
			args: args{
				numWorkers: 0,
				cap:        100,
			},
			want: &Pool[any, any]{
				numWorkers: runtime.NumCPU(),
				taskChan:   make(chan task[any, any], 100),
			},
		},
		{
			name: "test_cap_0",
			args: args{
				numWorkers: 10,
				cap:        0,
			},
			want: &Pool[any, any]{
				numWorkers: 10,
				taskChan:   make(chan task[any, any], 20),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPool[any, any](tt.args.numWorkers, tt.args.cap); cap(got.taskChan) != cap(tt.want.taskChan) ||
				got.numWorkers != tt.want.numWorkers {
				t.Errorf("NewPool() = %v, want %v", got, tt.want)
			}
		})
	}
}
