package snowflake

import (
	"testing"
)

func TestNextId(t *testing.T) {
	SetMachineID(0)
	for i := 0; i < 10; i++ {
		t.Log(ToTimeUnix(NextId()))
	}
}
