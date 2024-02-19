package bromide

import "regexp"

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
// add an optional title to a snapshot
// useful if tests have multiple snapshots
func WithSnapshotTitle(title string) Option {
	return func(c *config) {
		// Define a regular expression pattern to match invalid characters
		pattern := "[\\/\\\\\\x00:\\*\\?\"<>\\|&\\#]"

		// Compile the regular expression pattern
		regex := regexp.MustCompile(pattern)

		// Replace invalid characters with an empty string
		stripped := regex.ReplaceAllString(title, "")

		c.title = stripped
	}
}

// WithPassingNewSnapshots
// configure passing tests if new snapshots are created
func WithPassingNewSnapshots(pass bool) Option {
	return func(c *config) {
		c.passNew = pass
	}
}
