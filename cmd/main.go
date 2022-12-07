package main

import (
	"io"
	"os"

	"github.com/kitschysynq/hexed"
)

func main() {
	w := hexed.NewEncoder(os.Stdout)
	defer w.Close()

	io.Copy(w, os.Stdin)
}
