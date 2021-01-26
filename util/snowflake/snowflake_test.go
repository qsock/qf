package snowflake

import (
	"testing"
)

func TestNextId(t *testing.T) {
	SetMachineID(0)
	for i := 0; i < 10; i++ {
		id := NextId()
		t.Log(id, ToTimeUnix(NextId()))
	}
}

func BenchmarkNextId(b *testing.B) {
	for i := 0; i < b.N; i++ {
		id := NextId()
		b.Log(id, ToTimeUnix(NextId()))
	}
}

func TestMachineID(t *testing.T) {
	t.Log(machineId())
}
