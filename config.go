package bromide

import (
	"github.com/cobbinma/bromide/internal"
)

type Option func(*config)

type config struct {
	snapshotDir string
	passNew     bool
	title       string
}

// WithSnapshotDirectory
// configure a non default snapshot directory
func WithSnapshotDirectory(directory string) Option {
	return func(c *config) {
		c.snapshotDir = directory
	}
}

// WithSnapshotTitle
// appends an optional title to a snapshot
// useful if tests have multiple snapshots
func WithSnapshotTitle(title string) Option {
	return func(c *config) {
		c.title = internal.Sanitize(title)
	}
}

// WithPassingNewSnapshots
// configure passing tests if new snapshots are created
func WithPassingNewSnapshots(pass bool) Option {
	return func(c *config) {
		c.passNew = pass
	}
}
