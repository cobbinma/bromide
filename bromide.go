package bromide

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/hexops/valast"
	"github.com/sergi/go-diff/diffmatchpatch"
)

const folder = "snapshots"

func Snapshot[K comparable](t *testing.T, item K) {
	t.Helper()

	update := false
	if u := os.Getenv("UPDATE_SNAPSHOTS"); u != "" {
		v, err := strconv.ParseBool(u)
		if err != nil {
			t.Errorf("unable to parse boolean : %s", err.Error())
			return
		}

		update = v
	}

	currentDir, err := os.Getwd()
	if err != nil {
		t.Errorf("Error getting current directory: %v", err)
	}

	snapshotDir := fmt.Sprintf("%s/%s", currentDir, folder)
	name := fmt.Sprintf("%s/%s.snap", snapshotDir, t.Name())

	incoming := serialize(item)

	if err := os.MkdirAll(snapshotDir, 0755); err != nil {
		panic("unable to create snapshot directory")
	}

	file, err := os.Open(name)
	if err != nil {
		if !os.IsNotExist(err) {
			panic("unable to open existing snapshot")
		}

		file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0644)
		if err != nil {
			panic("unable to create new snapshot" + err.Error())
		}
		defer file.Close()

		if _, err := file.WriteString(incoming); err != nil {
			panic("unable to create new snapshot")
		}

		t.Errorf("created snapshot")
		return
	}
	defer file.Close()

	existing := new(strings.Builder)
	io.Copy(existing, file)

	if existing.String() != incoming {
		dmp := diffmatchpatch.New()

		diffs := dmp.DiffMain(existing.String(), incoming, true)

		t.Log("snapshot does not match")
		t.Log(dmp.DiffPrettyText(diffs))

		if update {
			os.WriteFile(name, []byte(incoming), 0644)
			return
		}

		t.Fail()
		return
	}
}

func serialize[K any](item K) string {
	return valast.StringWithOptions(item, &valast.Options{Unqualify: false})
}
