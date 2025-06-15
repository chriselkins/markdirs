package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func MarkDirs(root, fileName string, content []byte, overwrite, quiet bool) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error accessing %s: %v\n", path, err)
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

			file, err := os.OpenFile(targetPath, fileFlag, defaultPermission)

			if err != nil {
				// Only print error if not just skipping for not overwriting
				if !os.IsExist(err) || overwrite {
					fmt.Fprintf(os.Stderr, "Failed to create %s: %v\n", targetPath, err)
				}
			}

			_, werr := file.Write(content)
			file.Close()

			if werr != nil {
				fmt.Fprintf(os.Stderr, "Failed to write to %s: %v\n", targetPath, werr)
			} else if !quiet {
				fmt.Printf("Created %s\n", targetPath)
			}
		}

		return nil
	})
}
