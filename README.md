[![Build Status](https://github.com/chriselkins/markdirs/actions/workflows/go.yml/badge.svg)](https://github.com/chriselkins/markdirs/actions/workflows/go.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/chriselkins/markdirs.svg)](https://pkg.go.dev/github.com/chriselkins/markdirs)
[![Go Report Card](https://goreportcard.com/badge/github.com/chriselkins/markdirs)](https://goreportcard.com/report/github.com/chriselkins/markdirs)
[![GitHub release](https://img.shields.io/github/v/release/chriselkins/markdirs)](https://github.com/chriselkins/markdirs/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

# markdirs

**markdirs** is a command-line tool to recursively write a file with specified content to every directory under one or more root directories.
This is useful for marking directories (e.g., with a `.backup` or `.nomedia` file), automation, or batch configuration tasks.

## Quick Start

**Mark all directories under /data1 and /data2 with a .backup file containing "This folder is backed up":**

`markdirs /data1 /data2 .backup "This folder is backed up"`

**Or, pipe file content from stdin (the dash means “read from stdin”) and overwrite any existing NOTICE.txt files:**

``cat notice.txt | markdirs -o /var/www NOTICE.txt -``

## Features

* Recursively create a file in every directory under each root path
* Optionally overwrite existing files
* Content can be provided as a string or from standard input (when - is supplied)
* Quiet mode to suppress informational output
* "Best effort" mode: by default, errors are ignored and processing continues; use -f to fail fast on errors
* Simple, single binary – no dependencies

## Usage

```shell
markdirs [flags] <directory>... <file> <content | ->

Recursively write a file with the specified name and content to every folder under each <directory>.
```

### Flags

| Flag | Description                                     |
| ---- | ----------------------------------------------- |
| `-o` | Overwrite existing files (default: skip)        |
| `-q` | Quiet mode (suppress output)                    |
| `-m` | File permission/mode (e.g., 0644, 0600)         |
| `-f` | Fail immediately on error instead of continuing |
| `-h` | Show help and exit                              |
| `-v` | Print version and exit                          |

## Examples

**Write a file named .backup with the content “This folder is backed up” to every directory under /data1 and /data2:**

```shell
markdirs /data1 /data2 .backup "This folder is backed up"
```

**Write the contents of a local file to every directory under /var/www, overwriting existing files named NOTICE.txt:**

```shell
cat notice.txt | markdirs -o /var/www NOTICE.txt -
```

**Quietly mark all directories in /photos with empty file named .nomedia:**

```shell
markdirs -q /photos .nomedia ""
```

## Build

Requires Go 1.24 or newer.

```shell
make release
```

## Exit Codes

* `0`: Success
* `1`: Error (usage, I/O, or filesystem error)

## License

[MIT License](LICENSE)