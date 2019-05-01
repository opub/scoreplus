package util

import (
	"testing"
)

func TestNewID(t *testing.T) {
	for i := 0; i < 10000; i++ {
		id := RandomString(20)
		if len(id) != 20 {
			t.Errorf("id length invalid: %d", len(id))
		}
	}
}
