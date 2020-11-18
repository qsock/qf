package timewheel_test

import (
	"fmt"
	"github.com/qsock/qf/util/timewheel"
	"testing"
	"time"
)

var (
	timer = timewheel.New(5, 10*time.Millisecond)
)

func TestAdd(t *testing.T) {
	timewheel.Add(time.Millisecond*30, func() {
		fmt.Println(time.Now().String())
	})
	time.Sleep(time.Minute)
}

func Benchmark_Add(b *testing.B) {
	for i := 0; i < b.N; i++ {
		timer.Add(time.Minute, func() {

		})
	}
}

func Benchmark_StartStop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		timer.Start()
		timer.Stop()
	}
}
