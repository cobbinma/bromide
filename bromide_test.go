package bromide_test

import (
	"testing"

	"github.com/cobbinma/bromide"
)

type testStruct struct {
	example string
	again   int
}

func Test_Snapshot(t *testing.T) {
	hello := testStruct{
		example: "hello",
		again:   5,
	}
	bromide.Snapshot(t, hello)
}
