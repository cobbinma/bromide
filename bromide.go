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
		t.Errorf("error getting current directory: %v", err)
	}

	snapshotDir := fmt.Sprintf("%s/%s", currentDir, folder)
	acceptedPath := fmt.Sprintf("%s/%s%s", snapshotDir, t.Name(), internal.Accepted.Extension())
	pendingPath := fmt.Sprintf("%s/%s%s", snapshotDir, t.Name(), internal.Pending.Extension())

	incoming := serialize(item)

	if err := os.MkdirAll(snapshotDir, 0755); err != nil {
		panic("unable to create snapshot directory")
	}

	file, err := os.Open(acceptedPath)
	if err != nil {
		if !os.IsNotExist(err) {
			panic("unable to open existing snapshot")
		}

		// test does not have accepted snapshot
		file, err := os.OpenFile(pendingPath, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			panic("unable to create new snapshot" + err.Error())
		}
		defer file.Close()

		if _, err := file.WriteString(incoming); err != nil {
			panic("unable to create new snapshot")
		}

		t.Errorf("new snapshot 📸")
		t.Log("to update snapshots run `bromide review`")
		return
	}
	defer file.Close()

	// test has accepted snapshot
	existing := new(strings.Builder)
	io.Copy(existing, file)

	if existing.String() != incoming {
		diff := internal.Diff(existing.String(), incoming)

		t.Log("snapshot mismatch")
		t.Log(diff)
		t.Log("to update snapshots run `bromide review`")

		os.WriteFile(pendingPath, []byte(incoming), 0644)

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
