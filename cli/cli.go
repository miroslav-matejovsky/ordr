// Package cli is the command-line interface of ORDR.
//
// It parses arguments, calls the workspace and projections packages and
// presents results. It holds no parsing rules, graph algorithms or
// projection logic.
package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"

	"github.com/miroslav-matejovsky/ordr/projections"
	"github.com/miroslav-matejovsky/ordr/workspace"
)

// Exit codes returned by Run.
const (
	ExitOK    = 0 // success
	ExitError = 1 // invalid sources or I/O failure
	ExitUsage = 2 // invalid command line
)

const usage = `usage: ordr project [PROJECTION] [--workspace PATH] [--stdout]

PROJECTION  all (default), focus, value, uncertainty, complexity, relationships, readiness
--workspace folder inside the workspace; ordr.yaml is searched upward (default ".")
--stdout    print projections instead of writing them to the projections folder
`

var errUsage = errors.New("usage")

// Run executes the command line args (without the program name) and returns
// the process exit code.
func Run(args []string, stdout, stderr io.Writer) int {
	err := run(args, stdout)
	switch {
	case err == nil:
		return ExitOK
	case errors.Is(err, errUsage):
		// A failed write to stderr cannot be reported anywhere; the exit code remains.
		_, _ = fmt.Fprintf(stderr, "%v\n\n%s", err, usage)
		return ExitUsage
	default:
		_, _ = fmt.Fprintf(stderr, "error: %v\n", err)
		return ExitError
	}
}

func run(args []string, stdout io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("%w: missing command", errUsage)
	}
	if args[0] != "project" {
		return fmt.Errorf("%w: unknown command %q", errUsage, args[0])
	}

	fs := flag.NewFlagSet("project", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	dir := fs.String("workspace", ".", "")
	toStdout := fs.Bool("stdout", false, "")
	// Flags may appear before or after the projection name.
	if err := fs.Parse(args[1:]); err != nil {
		return fmt.Errorf("%w: %v", errUsage, err)
	}
	name := projections.All
	if fs.NArg() > 0 {
		name = fs.Arg(0)
		if err := fs.Parse(fs.Args()[1:]); err != nil {
			return fmt.Errorf("%w: %v", errUsage, err)
		}
	}
	if fs.NArg() > 0 {
		return fmt.Errorf("%w: unexpected arguments %v", errUsage, fs.Args())
	}
	kinds, err := projections.Select(name)
	if err != nil {
		return fmt.Errorf("%w: %v", errUsage, err)
	}

	ws, err := workspace.Discover(*dir)
	if err != nil {
		return err
	}
	register, g, err := ws.Load()
	if err != nil {
		return err
	}
	for _, k := range kinds {
		doc := projections.Render(k, register, g)
		if *toStdout {
			if _, err := fmt.Fprint(stdout, doc.Content); err != nil {
				return err
			}
			continue
		}
		written, err := ws.WriteProjection(doc.Name, doc.Content)
		if err != nil {
			return err
		}
		if _, err := fmt.Fprintf(stdout, "wrote %s\n", written); err != nil {
			return err
		}
	}
	return nil
}
