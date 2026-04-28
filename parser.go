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

// [return value which is associated with the key, if the key does not exist, return empty string]
func (c *Config) GetValue(key string) string {
	return c.data[key]
}

// [return value which is assocated with key and is converted to int, if the key does not exist or the value cannot be converted to int, return error]
func (c *Config) IntGetValue(key string) (int, error) {
	return strconv.Atoi(c.data[key])
}

// [return true if the key exists, otherwise return false]
func (c *Config) ExistsValue(key string) bool {
	_, ok := c.data[key]
	return ok
}
