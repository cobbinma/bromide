package internal

import (
	"github.com/sergi/go-diff/diffmatchpatch"
)

func Diff(old, new string) string {
	dmp := diffmatchpatch.New()

	diffs := dmp.DiffMain(old, new, true)

	return dmp.DiffPrettyText(diffs)
}
