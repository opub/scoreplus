package util

import (
	"fmt"
	"testing"
)

func TestNewID(t *testing.T) {

	for i := 0; i < 10000; i++ {
		id, err := NewID()
		if err != nil {
			t.Errorf("error in NewID: %s", err)
		}
		if len(id) != 12 {
			t.Errorf("id length invalid: %d", len(id))
		}
		fmt.Println(id)
	}
}
