package markdirs

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// MarkDirs writes a file with the specified name and content to every
// directory under root.
//
// If overwrite is true, existing files will be replaced. If quiet is true,
// informational output will be suppressed except for errors. If failFast is
// true, the operation stops on the first error encountered; otherwise, errors
// are logged and processing continues. Perm specifies the file permissions
// for the created files.
func MarkDirs(root, fileName string, content []byte, overwrite, quiet, failFast bool, perm os.FileMode) error {
	bytesReader := bytes.NewReader(content)

	return MarkDirsFromReaderAt(
		root, fileName, bytesReader, int64(len(content)),
		overwrite, quiet, failFast, perm,
	)
}

// MarkDirsFromReaderAt writes the content from src (an io.ReaderAt of length
// bytes) to a file with the specified name in every directory under root.
//
// The src must support random access (i.e., be an io.ReaderAt such as
// *os.File or bytes.Reader), as it will be read from the beginning for each
// file created.
//
// If overwrite is true, existing files will be replaced. If quiet is true,
// informational output will be suppressed except for errors. If failFast is
// true, the operation stops on the first error encountered; otherwise, errors
// are logged and processing continues. Perm specifies the file permissions
// for the created files.
func MarkDirsFromReaderAt(
	root, fileName string, src io.ReaderAt, length int64,
	overwrite, quiet, failFast bool, perm os.FileMode,
) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error accessing %s: %v\n", path, err)

			if failFast {
				return err
			}

			return nil
		}

		if d.IsDir() {
			targetPath := filepath.Join(path, fileName)

			if _, err := os.Stat(targetPath); err == nil && !overwrite {
				if !quiet {
					fmt.Printf("Skipped existing %s\n", targetPath)
				}

				return nil
			}

			fileFlag := func() int {
				if overwrite {
					// If overwriting, use O_TRUNC to clear the file
					return os.O_WRONLY | os.O_CREATE | os.O_TRUNC
				}

				// If not overwriting, use O_EXCL to ensure we don't overwrite existing files
				return os.O_WRONLY | os.O_CREATE | os.O_EXCL
			}()

			file, err := os.OpenFile(targetPath, fileFlag, perm)

			if err != nil {
				if os.IsExist(err) && !overwrite {
					if !quiet {
						fmt.Printf("Skipped existing %s\n", targetPath)
					}

					return nil // Skip existing file if not overwriting
				}

				fmt.Fprintf(os.Stderr, "Failed to create %s: %v\n", targetPath, err)

				if failFast {
					return err // Stop immediately on error
				}

				return nil
			}

			defer file.Close()

			r := io.NewSectionReader(src, 0, length)

			_, werr := io.Copy(file, r)

			if werr != nil {
				fmt.Fprintf(os.Stderr, "Failed to write to %s: %v\n", targetPath, werr)

				if failFast {
					return werr // Stop immediately on write error
				}
			}

			if !quiet {
				fmt.Printf("Created %s\n", targetPath)
			}
		}

		return nil
	})
}
