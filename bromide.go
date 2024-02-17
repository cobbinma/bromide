package bromide

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/hexops/valast"
	"github.com/sergi/go-diff/diffmatchpatch"
)

const folder = "snapshots"

func Snapshot[K comparable](t *testing.T, item K) {
	t.Helper()

	currentDir, err := os.Getwd()
	if err != nil {
		t.Errorf("Error getting current directory: %v", err)
	}

	snapshotDir := fmt.Sprintf("%s/%s", currentDir, folder)
	accepted := fmt.Sprintf("%s/%s.accepted", snapshotDir, t.Name())
	neww := fmt.Sprintf("%s/%s.new", snapshotDir, t.Name())

	incoming := serialize(item)

	if err := os.MkdirAll(snapshotDir, 0755); err != nil {
		panic("unable to create snapshot directory")
	}

	file, err := os.Open(accepted)
	if err != nil {
		if !os.IsNotExist(err) {
			panic("unable to open existing snapshot")
		}

		// test does not have accepted snapshot
		file, err := os.OpenFile(neww, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			panic("unable to create new snapshot" + err.Error())
		}
		defer file.Close()

		if _, err := file.WriteString(incoming); err != nil {
			panic("unable to create new snapshot")
		}

		t.Errorf("new snapshot 📸")
		return
	}
	defer file.Close()

	// test has accepted snapshot
	existing := new(strings.Builder)
	io.Copy(existing, file)

	if existing.String() != incoming {
		dmp := diffmatchpatch.New()

		diffs := dmp.DiffMain(existing.String(), incoming, true)

		t.Log("snapshot does not match")
		t.Log(dmp.DiffPrettyText(diffs))

		os.WriteFile(neww, []byte(incoming), 0644)

		t.Fail()
		return
	}
}

func serialize[K any](item K) string {
	return valast.StringWithOptions(item, &valast.Options{Unqualify: false})
}
