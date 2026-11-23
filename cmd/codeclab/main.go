// Command codeclab is the entry point: subcommands serve, decode, encode, inspect and verify.
package main

import (
	"fmt"
	"io"
	"os"
)

const version = "0.0.0-scaffold"

const usage = `usage: codeclab <command> [arguments]

commands:
  serve     start the loopback server
  decode    decode a file
  encode    encode a raster
  inspect   inspect a file
  verify    run the release checks

flags:
  --help      show this help
  --version   show the version
`

// parents maps each command that is not built yet to its parent ticket.
var parents = map[string]string{
	"serve":   "#2",
	"decode":  "#3",
	"encode":  "#4",
	"inspect": "#3",
	"verify":  "#7",
}

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprint(stderr, usage)
		return 2
	}
	switch args[0] {
	case "--help", "-h", "help":
		fmt.Fprint(stdout, usage)
		return 0
	case "--version", "version":
		fmt.Fprintln(stdout, "codeclab "+version)
		return 0
	}
	if ticket, ok := parents[args[0]]; ok {
		fmt.Fprintf(stderr, "not implemented yet: see %s\n", ticket)
		return 1
	}
	fmt.Fprintf(stderr, "unknown command %q\n\n%s", args[0], usage)
	return 2
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}
