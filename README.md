# bromide

A snapshot testing library for go 📸

## Introduction

Bromide is a snapshot library for Go, designed to simplify managing snapshot tests. 

With Bromide, you can easily capture test output and check against an expected value.

## Usage

### Write a test

```sh
go get github.com/cobbinma/bromide
```

```go
import github.com/cobbinma/bromide

func TestSomething(t *testing.T) {
    something := "something"
    
    bromide.Snapshot(t, something)
}
```

### Review snapshots

```sh
go install github.com/cobbinma/bromide/cmd/bromide@master
```

```sh
bromide review
```
