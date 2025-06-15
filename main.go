package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

const (
	defaultPermission = 0644 // Default file permission for created files
)

var (
	// Version is the version of the tool, set at build time
	Version   = "dev"
	Commit    = ""
	BuildDate = ""

	// Command-line flags
	overwrite   = flag.Bool("o", false, "Overwrite existing file if present")
	quiet       = flag.Bool("q", false, "Quiet mode (suppress output except errors)")
	help        = flag.Bool("h", false, "Show help and exit")
	showVersion = flag.Bool("v", false, "Print version and exit")
)

func usage() {
	fmt.Fprintf(os.Stderr, `Usage: %s [flags] <directory> <file> <content|- for stdin>

Recursively write a file with the specified name and content to every folder under <directory>.

Flags:
  -o    Overwrite existing files (default: do not overwrite)
  -q    Quiet mode (no informational output)
  -h    Show this help message
`, os.Args[0])
}

func main() {
	os.Exit(run())
}

func run() int {
	flag.Usage = usage
	flag.Parse()

	if *help {
		usage()
		return 0
	}

	if *showVersion {
		fmt.Printf("Version: %s\nCommit: %s\nBuild Date: %s\n", Version, Commit, BuildDate)
		os.Exit(0)
	}

	args := flag.Args()

	if len(args) != 3 {
		usage()
		return 1
	}

	root, fileName, contentArg := args[0], args[1], args[2]

	// content for the file can be provided as a string or read from stdin
	content, err := func() ([]byte, error) {
		if contentArg == "-" {
			return io.ReadAll(os.Stdin)
		}

		return []byte(contentArg), nil
	}()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read from stdin: %v\n", err)
		return 1
	}

	err = MarkDirs(root, fileName, content, *overwrite, *quiet)

	if err != nil {
		fmt.Fprintf(os.Stderr, "WalkDir error: %v\n", err)
		return 1
	}

	return 0
}
