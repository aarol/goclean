# Goclean

### Clean your filesystem with ease

<img title="" src="resource/screenshot.png" alt="Screenshot" data-align="center">

Goclean is written in Go, using the [bubbletea](https://github.com/charmbracelet/bubbletea) framework.

--- 

1. Enter the directories you want to find

2. Press <kbd>space</kbd>/<kbd>delete</kbd> to delete the directory

3. Press <kbd>Q</kbd> to quit when you're done

### Installation

Go 1.20+ is required. [Install go here](https://go.dev/doc/install)

```bash
go install github.com/aarol/goclean@latest
```

This should build the program for your system. You should also add add "$(go env GOPATH)/bin" to your $PATH for this command to be accessible anywhere on your system.

## Usage

```bash
goclean <dir> (dir2 dir3...)
```

### Options

* `-e`,`--exclude "<dir> (dir2...)"`: Exclude directories from search (default: node_modules, .git)
* `-a`,`--all`: Include hidden directories (directories starting with ".")
