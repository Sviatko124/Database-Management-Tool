# EVA - Cheatsheet/Notes Database tool
EVA database tool gives you easy access to search through your hacking notes and cheatsheet, all from the comfort of your terminal. The tool is intuitive, and after adding all of your notes, allows you to quickly query and modify your notes, so that you can find exactly what you need easily. This tool is perfect for red teamers who have a lot of notes through which they have to manually find what they're looking for. 

## Features

- Store and manage exploit notes, commands, and attack step, so that you can categorize each entry to a specific penetration testing process step. 
- Search all text in an entry for your given search keywords
- Efficient search functionality, so if you enter a search with two keyword (keywords are separated by a comma), only the entries containing both keywords will be displayed
- Pretty color-coded terminal interface
- Automatic entry reindexing upon entry deletion
- SQLite3 backend for quick and reliable data storage

## Build from source

Requirements:
- Go 1.x or higher
- GCC
- musl-dev (for static compilation)

```bash
# Install dependencies (Ubuntu/Debian)
sudo apt-get install gcc musl-dev

# Clone the repository
git clone https://github.com/Sviatko124/eva.git
cd eva

# Build static binary
CGO_ENABLED=1 go build -ldflags="-s -w -linkmode external -extldflags '-static'" src/eva.go

# Optional but highly recommended: move the binary to /usr/local/bin so that you can run the program from anywhere in your system
mv eva /usr/local/bin/eva
```
## Usage
Just run:
eva

You will be greeted with a simple interface. Upon first launch, a database will be create in your home directory under .eva/eva.db. To start creating notes, just begin adding them by selecting the "Add entry" option. By the way, make sure to pick easy to remember keywords for your entries, so that you can find what you're looking for easily. 
