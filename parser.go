package confparser

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

type Config struct {
	data map[string]string
}

// read and parse the configuration from the reader, return a Config struct, if there is an error during reading, return an error
func Parse(r io.Reader) (*Config, error) {
	conf := &Config{
		data: make(map[string]string),
	}

	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()

		// serch "#"
		if idx := strings.Index(line, "#"); idx != -1 {
			// if "#" found, then eliminate strings after "#"
			line = line[:idx]
		}

		// 2. eliminate spases
		line = strings.TrimSpace(line)

		// if no strings, skip the line
		if line == "" {
			continue
		}

		// devide into key and value
		parts := strings.Fields(line)
		if len(parts) != 2 {
			continue
		}

		key := parts[0]
		value := parts[1]
		conf.data[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return conf, nil
}

// return value which is associated with the key, if the key does not exist, return empty string
func (c *Config) GetValue(key string) string {
	return c.data[key]
}

// return value which is assocated with key and is converted to int, if the key does not exist or the value cannot be converted to int, return error
func (c *Config) IntGetValue(key string) (int, error) {
	return strconv.Atoi(c.data[key])
}

// return true if the key exists, otherwise return false
func (c *Config) ExistsValue(key string) bool {
	_, ok := c.data[key]
	return ok
}

type MultiConfig struct {
	data map[string][]string
}

// read and parse the configuration from the reader, return a MultiConfig struct, if there is an error during reading, return an error. The MultiConfig struct should support multiple values for the same key, and the values should be stored in a slice of strings.
func ParseMultipleValues(r io.Reader) (*MultiConfig, error) {
	conf := &MultiConfig{
		data: make(map[string][]string),
	}

	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()

		if idx := strings.Index(line, "#"); idx != -1 {
			line = line[:idx]
		}

		line = strings.TrimSpace(line)

		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		key := parts[0]
		value := parts[1]

		conf.data[key] = append(conf.data[key], value)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return conf, nil
}

// return a slice of strings which are associated with the key, if the key does not exist, return an empty slice
func (c *MultiConfig) GetMultipleValues(key string) []string {
	return c.data[key]
}

// return the first value which is associated with the key, if the key does not exist, return empty string
func (c *MultiConfig) GetFirstValue(key string) string {
	values := c.data[key]
	if len(values) > 0 {
		return values[0]
	}
	return ""
}

// return true if the key exists, otherwise return false
func (c *MultiConfig) ExistsMultipleValues(key string) bool {
	_, ok := c.data[key]
	return ok
}
