# Doublemap
a Go container combining a map with a reverse map

[![GoDoc](https://godoc.org/github.com/rasteric/doublemap/go?status.svg)](https://godoc.org/github.com/rasteric/doublemap)
[![Go Report Card](https://goreportcard.com/badge/github.com/rasteric/doublemap)](https://goreportcard.com/report/github.com/rasteric/doublemap)

This package provides a map with a reverse map mechanism, i.e., it allows you to get the value by a given key and a key from a given value.

## Example usage

```
package main

import (
  "fmt"

  "github.com/rasteric/doublemap"
)

func main() {
    m := doublemap.New[string, int]()
    m.Set("first", 1)
    m.Set("second", 2)
    m.Set("third", 3)
    v, _ := m.Get("first")
    fmt.Println(v)
    k, _ := m.ByValue(3)
    fmt.Println(k)
}
```

See the reference for more information.

## License

This package is provided under the permissive MIT License, please see the accompanying LICENSE agreement for more information.
