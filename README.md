# bromide

A snapshot testing library for go 📸

## Introduction

Bromide is a snapshot library for Go, designed to simplify the process of taking and managing snapshots for testing purposes. 
With Bromide, you can easily capture test output and then compare to expected values.

## Usage

### Write a test

```sh
go get github.com/cobbinma/bromide
```

```go
import github.com/cobbinma/bromide

func TestSomething(t *testing.T) {
    something := "something"
    
    bromide.SnapshotT(t, something)
}
```

### Review snapshots

```sh
go install github.com/cobbinma/bromide/cmd/bromide@master
```

```sh
bromide review
```
