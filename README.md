# Hex-editor style encoder for go

[![PkgGoDev](https://pkg.go.dev/badge/github.com/kitschysynq/hexed)](https://pkg.go.dev/github.com/kitschysynq/hexed)

## Installation

```shell
go get github.com/kitschysynq/hexed
```

## Quickstart

```go
package main

import (
	"fmt"
	"os"

	"github.com/kitschysynq/hexed"
)

func main() {
	w := hexed.NewEncoder(os.Stdout)
	defer w.Close()

	fmt.Fprintf(w, "this is a totally rad example")
}
```

And that will produce output like this:

```
00000000: 7468 6973 2069 7320 6120 746f 7461 6c6c  this is a totall
00000010: 7920 7261 6420 6578 616d 706c 65         y rad example
```
