package utils

import (
	"testing"
)

func TestRandString(t *testing.T) {
	for i := 1; i < 100; i++ {
		if RandString(i) == RandString(i) {
			t.Errorf("RandString should generate different result")
		}
	}
}
