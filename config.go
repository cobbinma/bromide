package bromide

type Option func(*config)

type config struct {
	snapshotDir string
	passNew     bool
}

// WithSnapshotDirectory
// configure a non default snapshot directory
func WithSnapshotDirectory(directory string) Option {
	return func(c *config) {
		c.snapshotDir = directory
	}
}

// WithPassingNewSnapshots
// configure passing tests if new snapshots are created
func WithPassingNewSnapshots(pass bool) Option {
	return func(c *config) {
		c.passNew = pass
	}
}
