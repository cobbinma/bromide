// Bromide is a snapshot library, designed to simplify managing snapshot tests.
//
// Snapshot tests are useful if inputs are large or change often.
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

// Snapshot compares a given input against a reference value.
// The test will fail if the value does not match, or a new snapshot is created.
//
// If the test fails, use bromide to interactively review changes.
// ```
// $ bromide
// ````
func Snapshot[K any](t *testing.T, item K, options ...Option) {
	t.Helper()

	config := &config{}
	for _, option := range options {
		option(config)
	}

	dir := config.snapshotDir
	if dir == "" {
		wd, err := os.Getwd()
		if err != nil {
			t.Error("bromide: unable to get working directory")
			t.Log(err.Error())
			return

		}

		dir = fmt.Sprintf("%s/snapshots", wd)
	}

	title := ""
	if config.title != "" {
		title = "_" + config.title
	}

	acceptedPath := fmt.Sprintf("%s/%s%s%s", dir, t.Name(), title, internal.Accepted.Extension())
	pendingPath := fmt.Sprintf("%s/%s%s%s", dir, t.Name(), title, internal.Pending.Extension())

	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Error("bromide: unable to create snapshot directory")
		t.Log(err.Error())
		return
	}

	incoming := serialize(item)

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

		t.Log("new snapshot 📸")
		t.Log("to accept snapshot run `bromide`")
		if !config.passNew {
			t.Fail()
		}
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
		t.Log("to update snapshot run `bromide`")

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
