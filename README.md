# config-file-parser-go

`config-file-parser-go` is a simple and lightweight Key-Value configuration file parser for Go. 
It easily parses simple, space-separated formats commonly found in Linux config files, making them easy to use within your Go applications.

## Features

- **Simple:** Zero external dependencies (uses only the Go standard library).
- **Lightweight:** Accepts an `io.Reader`, ensuring a low memory footprint even when processing large files.
- **Intuitive Rules:**
  - Use Parse() method.
  - Parses space- or tab-separated entries (e.g., `key value`).
  - Ignores anything after a `#` as a comment.
  - Safely skips empty lines and lines with three or more elements (e.g., `key value1 value2`) as invalid formats.
- **Multiple Values:**
  - Use ParseMultipleValues() method
  - All values are stored into slices

## Installation

```bash
go get github.com/ebishrimp/config-file-parser-go
```

## Sample


Usage
1. Prepare a Configuration File
Create a config file (e.g., config.conf) that you want to parse.

```bash
# Example of config.conf
port 8080
host localhost

# Lines with 3 or more elements will be ignored
invalid_key too many values
```
```bash
# Example of config.conf which can be parsed by ParseMultipleValues(r io.Reader)
# The same key can hold multiple values
multiplevalues value1
multiplevalues value2

# Simultaneous description of values will be ignored
invalid_form value3 value4 value5
```
2. Use in Your Go Code
You can parse directly from a file or a string, as long as it satisfies the io.Reader interface.

```go
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

	// Get as a string
	host := conf.GetValue("host")
	fmt.Printf("Host: %s\n", host)

	// Get as an integer
	port, err := conf.IntGetValue("port")
	if err != nil {
		log.Printf("Failed to get port: %v", err)
	} else {
		fmt.Printf("Port: %d\n", port)
	}

	// Check if a key exists
	if conf.ExistsValue("invalid_key") {
		fmt.Println("invalid_key exists")
	} else {
		fmt.Println("invalid_key does not exist because it was skipped")
	}

	// 4. Parse multiple values
	multiConf, err := confparser.ParseMultipleValues(f)
	if err != nil {
		log.Fatalf("Parse error: %v", err)
	}

	// 5. Retrieve values

	// Get string slice
	multival := multiConf.GetMultipleValues("multiplevalues")
	// multival[0] == "value1"
	// multival[1] == "value2"

	// Get a first contest of slice
	firstval := multiConf.GetFirstValue("multiplevalues")
	// firstval == "value1"

	// Check if a key exists
	if conf.ExistsMultipleValues("invalid_form") {
		fmt.Println("invalid_form exists")
	} else {
		fmt.Println("invalid_form does not exist because it was skipped")
	}
}
