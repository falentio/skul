package xrand

import "testing"

func TestSmol(t *testing.T) {
	for i := 0; i < 100; i++ {
		t.Log(Smol.GeneratePassword(8))
	}
}
