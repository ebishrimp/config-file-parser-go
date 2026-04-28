# config-file-parser-go

`config-file-parser-go` is a simple and lightweight Key-Value configuration file parser for Go. 
It easily parses simple, space-separated formats commonly found in Linux config files, making them easy to use within your Go applications.

## Features

- **Simple:** Zero external dependencies (uses only the Go standard library).
- **Lightweight:** Accepts an `io.Reader`, ensuring a low memory footprint even when processing large files.
- **Intuitive Rules:**
  - Parses space- or tab-separated entries (e.g., `key value`).
  - Ignores anything after a `#` as a comment.
  - Safely skips empty lines and lines with three or more elements (e.g., `key value1 value2`) as invalid formats.

## Installation

```bash
go get github.com/ebishrimp/config-file-parser-go
```

## Sample

```bash
Usage
1. Prepare a Configuration File
Create a config file (e.g., config.conf) that you want to parse.

Plaintext
# Example of config.conf
port 8080
host localhost

# Lines with 3 or more elements will be ignored
invalid_key too many values
2. Use in Your Go Code
You can parse directly from a file or a string, as long as it satisfies the io.Reader interface.


package main

import (
	"fmt"
	"log"
	"os"

	confparser "github.com/ebishrimp/config-file-parser-go"
)

func main() {
	// 1. Open the file (you can also pass a string using strings.NewReader)
	f, err := os.Open("config.conf")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// 2. Parse the configuration
	conf, err := confparser.Parse(f)
	if err != nil {
		log.Fatalf("Parse error: %v", err)
	}

	// 3. Retrieve values

	// Get as a string (Return)
	host := conf.GetValue("host")
	fmt.Printf("Host: %s\n", host)

	// Get as an integer (IntReturn)
	port, err := conf.IntGetValue("port")
	if err != nil {
		log.Printf("Failed to get port: %v", err)
	} else {
		fmt.Printf("Port: %d\n", port)
	}

	// Check if a key exists (Exists)
	if conf.ExistsValue("invalid_key") {
		fmt.Println("invalid_key exists")
	} else {
		fmt.Println("invalid_key does not exist because it was skipped")
	}
}
