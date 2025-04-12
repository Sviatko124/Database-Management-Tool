# Cheatsheet/Notes Database tool
This database management program written in Go allows you to easily search through your hacking notes and cheatsheets, all from the comfort of your terminal. The tool is intuitive, and after adding all of your notes, it allows you to quickly query and modify your notes so that you can easily find exactly what you need. This tool is perfect for red teamers who have a lot of notes through which they have to manually find what they're looking for. 

## Features

- Store and manage exploit notes, commands, and attack steps so that you can categorize each entry as a specific penetration testing process step. 
- Search all text in an entry for your given search keywords
- Efficient search functionality, so if you enter a search with two keywords (keywords are separated by a comma), only the entries containing both keywords will be displayed
- Pretty color-coded terminal interface
- Automatic entry reindexing upon entry deletion
- SQLite3 backend for quick and reliable data storage

## Build from source

Requirements:
- Go 1.19 or higher
- GCC
- musl-dev (for static compilation)

```bash
# Install dependencies (Ubuntu/Debian)
sudo apt-get install gcc musl-dev golang-go

# Clone the repository
git clone https://github.com/Sviatko124/eva.git
cd eva

# Initialize project, install dependencies, and finalize environment
go mod init eva
go get github.com/mattn/go-sqlite3
go get golang.org/x/term
go mod tidy

# Build static binary
CGO_ENABLED=1 go build -ldflags="-s -w -linkmode external -extldflags '-static'" eva.go

# Optional but highly recommended: move the binary to /usr/local/bin so that you can run the program from anywhere in your system
sudo mv eva /usr/local/bin/eva
```
If you get permission errors when running any of the commands, just put `sudo` before the command, which should fix the error. 

## Usage
Just run:
eva

You will be greeted with a simple interface. Upon first launch, a database will be created in your home directory under .eva/eva.db. To start creating notes, just begin adding them by selecting the "Add entry" option. By the way, make sure to pick easy to remember keywords for your entries, so that you can find what you're looking for easily. 
