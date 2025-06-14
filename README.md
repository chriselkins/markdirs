# markdirs

**markdirs** is a command-line tool to recursively write a file with specified content to every directory under a given root directory.
This is useful for marking directories (e.g., with a `.backup` or `.nomedia` file), automation, or batch configuration.

## Features

* Recursively create a file in every directory under a root path
* Optionally overwrite existing files
* Content can be provided as a string or from standard input
* Quiet mode to suppress informational output
* Simple, single binary – no dependencies

## Usage

```shell
markdirs [flags] <directory> <file> <content | ->

Recursively write a file with the specified name and content to every folder under <directory>.
```

### Flags

| Flag | Description                              |
| ---- | ---------------------------------------- |
| `-o` | Overwrite existing files (default: skip) |
| `-q` | Quiet mode (suppress output)             |
| `-h` | Show help and exit                       |
| `-v` | Print version and exit                   |

## Examples

**Write a file named **`** with the content “This folder is backed up” to every directory under **`**:**

```shell
markdirs /data .backup "This folder is backed up"
```

**Write the contents of a local file (**`**) to every directory under **`**, overwriting existing files:**

```shell
cat notice.txt | markdirs -o /var/www NOTICE.txt -
```

**Quietly mark all directories with a **\`\`** file:**

```shell
markdirs -q /photos .nomedia ""
```

## Build

Requires Go 1.18 or newer.

```shell
go build -o markdirs
```

## Exit Codes

* `0`: Success
* `1`: Error (usage, I/O, or filesystem error)

## License

[MIT License](LICENSE)