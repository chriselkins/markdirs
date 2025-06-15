package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/chriselkins/markdirs/markdirs"
)

var (
	// version information is set at build time in the release.sh script
	Version   = "dev"
	Commit    = ""
	BuildDate = ""
)

func usage() {
	fmt.Fprintf(os.Stderr, `Usage: %s [flags] <directory>... <file> <content|- for stdin>

Recursively write a file with the specified name and content to every folder under each <directory>.

Flags:
  -o    Overwrite existing files (default: do not overwrite)
  -q    Quiet mode (suppress output except errors)
  -m    File permission/mode (default: 0644, e.g., 0600, 0640)
  -f    Fail fast on error (default: continue processing)
  -v    Show version information and exit
  -h    Show this help message
`, os.Args[0])
}

func main() {
	os.Exit(run())
}

func run() int {
	// Command-line flags
	overwrite := flag.Bool("o", false, "Overwrite existing file if present")
	quiet := flag.Bool("q", false, "Quiet mode (suppress output except errors)")
	mode := flag.String("m", "0644", "File permission/mode (default: 0644, e.g., 0600, 0640)")
	failFast := flag.Bool("f", false, "Fail immediately on error instead of continuing")
	help := flag.Bool("h", false, "Show help and exit")
	showVersion := flag.Bool("v", false, "Print version and exit")

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

	if len(args) < 3 {
		usage()
		return 1
	}

	// Everything but last two are directories to process
	dirs := args[:len(args)-2]
	fileName := args[len(args)-2]
	contentArg := args[len(args)-1]

	modeInt, err := strconv.ParseUint(*mode, 8, 32)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid file mode: %v\n", err)
		return 1
	}

	fileMode := os.FileMode(modeInt)

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

	for _, dir := range dirs {
		err = markdirs.MarkDirs(dir, fileName, content, *overwrite, *quiet, *failFast, fileMode)

		if err != nil {
			fmt.Fprintf(os.Stderr, "WalkDir error: %v\n", err)

			if *failFast {
				return 1
			}
		}
	}

	return 0
}
