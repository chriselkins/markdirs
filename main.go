package main

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

const (
	// Version is the current version of the tool.
	Version = "1.0.0"

	defaultPermission = 0644 // Default file permission for created files
)

var (
	overwrite = flag.Bool("o", false, "Overwrite existing file if present")
	quiet     = flag.Bool("q", false, "Quiet mode (suppress output except errors)")
	help      = flag.Bool("h", false, "Show help and exit")
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

	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error accessing %s: %v\n", path, err)
			return nil
		}

		if d.IsDir() {
			targetPath := filepath.Join(path, fileName)

			if _, err := os.Stat(targetPath); err == nil && !*overwrite {
				if !*quiet {
					fmt.Printf("Skipped existing %s\n", targetPath)
				}

				return nil
			}

			fileFlag := func() int {
				if *overwrite {
					// If overwriting, use O_TRUNC to clear the file
					return os.O_WRONLY | os.O_CREATE | os.O_TRUNC
				}

				// If not overwriting, use O_EXCL to ensure we don't overwrite existing files
				return os.O_WRONLY | os.O_CREATE | os.O_EXCL
			}()

			file, err := os.OpenFile(targetPath, fileFlag, defaultPermission)

			if err != nil {
				// Only print error if not just skipping for not overwriting
				if !os.IsExist(err) || *overwrite {
					fmt.Fprintf(os.Stderr, "Failed to create %s: %v\n", targetPath, err)
				}
			}

			_, werr := file.Write(content)
			file.Close()

			if werr != nil {
				fmt.Fprintf(os.Stderr, "Failed to write to %s: %v\n", targetPath, werr)
			} else if !*quiet {
				fmt.Printf("Created %s\n", targetPath)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "WalkDir error: %v\n", err)
		return 1
	}

	return 0
}
