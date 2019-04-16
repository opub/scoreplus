package util

import (
	"testing"
)

func TestNewID(t *testing.T) {
	for i := 0; i < 10000; i++ {
		id, err := RandomID()
		if err != nil {
			t.Errorf("error in NewID: %s", err)
		}
		if len(id) != 10 {
			t.Errorf("id length invalid: %d", len(id))
		}
	}
}
