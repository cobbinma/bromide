package internal

type ReviewState int

const (
	Pending ReviewState = iota
	Accepted
)

func (s ReviewState) Extension() (extension string) {
	switch s {
	case Pending:
		extension = ".new"
	case Accepted:
		extension = ".snap"
	}

	return
}
