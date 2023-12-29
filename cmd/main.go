package main

import (
	"os"

	"github.com/taylormonacelli/jitteryrambo"
)

func main() {
	code := jitteryrambo.Execute()
	os.Exit(code)
}
