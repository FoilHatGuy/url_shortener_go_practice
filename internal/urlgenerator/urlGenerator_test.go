package urlgenerator

import (
	"testing"
)

func TestRandSeq(t *testing.T) {
	if seq := RandSeq(10); len(seq) != 10 {
		t.Error("length of random sequence is wrong")
	}
}
