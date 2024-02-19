package bromide_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/cobbinma/bromide"
	"github.com/cobbinma/bromide/internal"
)

type testStruct struct {
	example *string
	again   int
}

func Test_Snapshot(t *testing.T) {
	text := "hello"
	hello := testStruct{
		example: &text,
		again:   5,
	}

	bromide.Snapshot(t, hello, bromide.WithSnapshotDirectory("./internal/snapshots"))
}

func Test_PendingSnapshot(t *testing.T) {
	dir := t.TempDir()

	bromide.Snapshot(t, "something",
		bromide.WithSnapshotDirectory(dir),
		bromide.WithSnapshotTitle("title*/<"),
		bromide.WithPassingNewSnapshots(true))

	if _, err := os.Stat(fmt.Sprintf("%s/%s_title%s", dir, t.Name(), internal.Pending.Extension())); err != nil {
		t.Error(err.Error())
	}
}
