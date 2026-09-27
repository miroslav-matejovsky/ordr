// Command ordr generates decision-support projections from an ORDR workspace.
package main

import (
	"os"

	"github.com/miroslav-matejovsky/ordr/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:], os.Stdout, os.Stderr))
}
