package bromide

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/cobbinma/bromide/internal"
	"github.com/davecgh/go-spew/spew"
)

const folder = "snapshots"

func Snapshot[K comparable](t *testing.T, item K) {
	t.Helper()

	currentDir, err := os.Getwd()
	if err != nil {
		t.Error("bromide: unable to get working directory")
		t.Log(err.Error())
		return
	}

	snapshotDir := fmt.Sprintf("%s/%s", currentDir, folder)
	acceptedPath := fmt.Sprintf("%s/%s%s", snapshotDir, t.Name(), internal.Accepted.Extension())
	pendingPath := fmt.Sprintf("%s/%s%s", snapshotDir, t.Name(), internal.Pending.Extension())

	incoming := serialize(item)

	if err := os.MkdirAll(snapshotDir, 0755); err != nil {
		t.Error("bromide: unable to create snapshot directory")
		t.Log(err.Error())
		return
	}

	file, err := os.Open(acceptedPath)
	if err != nil {
		if !os.IsNotExist(err) {
			t.Error("bromide: unable to open accepted snapshot")
			t.Log(err.Error())
			return
		}

		// test does not have accepted snapshot
		file, err := os.OpenFile(pendingPath, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			t.Error("bromide: unable to open pending snapshot")
			t.Log(err.Error())
			return
		}
		defer file.Close()

		if _, err := file.WriteString(incoming); err != nil {
			t.Error("bromide: unable to write pending snapshot")
			t.Log(err.Error())
			return
		}

		t.Errorf("new snapshot 📸")
		t.Log("to update snapshots run `bromide`")
		return
	}
	defer file.Close()

	// test has accepted snapshot
	existing := new(strings.Builder)
	if _, err := io.Copy(existing, file); err != nil {
		t.Error("bromide: unable to copy accepted file")
		t.Log(err.Error())
		return
	}

	if existing.String() != incoming {
		diff := internal.Diff(existing.String(), incoming)

		t.Log("snapshot mismatch")
		t.Log("\n" + diff)
		t.Log("to update snapshots run `bromide`")

		if err := os.WriteFile(pendingPath, []byte(incoming), 0644); err != nil {
			t.Error("bromide: unable to write pending snapshot")
			t.Log(err.Error())
			return
		}

		t.Fail()
		return
	}
}

func serialize[K any](item K) string {
	config := &spew.ConfigState{
		Indent:                  "  ",
		SortKeys:                true,
		DisablePointerAddresses: true,
		DisableCapacities:       true,
		SpewKeys:                true,
		DisableMethods:          false,
	}
	return config.Sdump(item)
}
