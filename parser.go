package confparser

import (
	"bufio"
	"fmt"
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

// Recommended to use ConfigurationMap instead of MultiConfig or Config. The ConfigurationMap struct has unified functions to get value, get multiple values, check if key exists, and it can handle both single and multiple values for the same key.
type ConfigurationMap struct {
	data map[string][]string
}

type ParseError struct {
	line int
	msg  string
}

func (err *ParseError) Error() string {
	return "line " + strconv.Itoa(err.line) + ": " + err.msg
}

// Read and parse the configuration through the io.Reader, return a ConfigurationMap struct. If there is an error during reading, return a ParseError struct which contains the line number and the error message. The ConfigurationMap struct should support multiple values for the same key.
func ParseConfig(r io.Reader) (*ConfigurationMap, *ParseError) {
	conf := &ConfigurationMap{
		data: make(map[string][]string),
	}

	scanner := bufio.NewScanner(r)

	var l int = 0 //line number
	for scanner.Scan() {
		l++
		line := scanner.Text()

		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		if parts := strings.TrimSpace(line); len(parts) < 2 {
			return nil, &ParseError{line: l, msg: "no value for key found"}
		}

		p := strings.Fields(line)
		if len(p) >= 3 {
			if strings.HasPrefix(p[1], "\"") && strings.HasSuffix(p[len(p)-1], "\"") {
				key := p[0]

				str := strings.Join(p[1:], " ")
				str = strings.Trim(str, "\"")

				conf.data[key] = append(conf.data[key], str)
			} else {
				return nil, &ParseError{line: l, msg: "missing quotes for value with spaces, or too many values for a key"}
			}
		} else {
			parts := strings.Fields(line)
			key := parts[0]
			value := parts[1]
			conf.data[key] = append(conf.data[key], value)
		}
	}
	return conf, nil
}

// Return a slice of strings which are associated with the key, if the key does not exist, return an empty slice
func (c *ConfigurationMap) StringSlice(key string) []string {
	return c.data[key]
}

// Return the value which is associated with the key, if the key does not exist, return empty string. If there are multiple values for the same key, return the first value.
func (c *ConfigurationMap) String(key string) string {
	values := c.data[key]
	if len(values) > 0 {
		return values[0]
	}
	return ""
}

// Return a slice of integers which are associated with the key, if the key does not exist, return an empty slice.
func (c *ConfigurationMap) IntSlice(key string) ([]int, error) {
	var result []int
	for i := 0; i < len(c.data[key]); i++ {
		if intdata, err := strconv.Atoi(c.data[key][i]); err != nil {
			return nil, err
		} else {
			result = append(result, intdata)
		}
	}
	return result, nil
}

// Return the value which is associated with the key and is converted to int, if the key does not exist or the value cannot be converted to int, return error. If there are multiple values for the same key, return the first value.
func (c *ConfigurationMap) Int(key string) (int, error) {
	values := c.data[key]
	if len(values) > 0 {
		return strconv.Atoi(values[0])
	}
	return 0, fmt.Errorf("key not found or no valid integer value")
}

// Return a slice of floats which are associated with the key, if the key does not exist, return an empty slice. If there is an error, return an error.
func (c *ConfigurationMap) FloatSlice(key string) ([]float64, error) {
	var result []float64
	for i := 0; i < len(c.data[key]); i++ {
		if floatdata, err := strconv.ParseFloat(c.data[key][i], 64); err != nil {
			return nil, err
		} else {
			result = append(result, floatdata)
		}
	}
	return result, nil
}

// Return the value which is associated with the key and is converted to float, if the key does not exist or the value cannot be converted to float, return error. If there are multiple values for the same key, return the first value.
func (c *ConfigurationMap) Float(key string) (float64, error) {
	values := c.data[key]
	if len(values) > 0 {
		return strconv.ParseFloat(values[0], 64)
	}
	return 0, fmt.Errorf("key not found or no valid float value")
}

// Return the value which is associated with the key and is converted to bool, if the key does not exist or the value cannot be converted to bool, return error. Multiple values for the same key are not allowed for boolean values, if there are multiple values, return error.
func (c ConfigurationMap) Bool(key string) (bool, error) {
	value := c.data[key][0]
	if len(value) == 1 {
		if booldata, err := strconv.ParseBool(value); err != nil {
			if strings.EqualFold(value, "yes") || strings.EqualFold(value, "y") {
				return true, nil
			} else if strings.EqualFold(value, "no") || strings.EqualFold(value, "n") {
				return false, nil
			}
			return false, fmt.Errorf("key not found or no valid boolean value")
		} else {
			return booldata, nil
		}
	}
	return false, fmt.Errorf("key not found or no valid boolean value")
}

// Return true if the key exists, otherwise return false
func (c *ConfigurationMap) Exists(key string) bool {
	_, ok := c.data[key]
	return ok
}

// Return the number of values which are associated with the key, if the key does not exist, return 0. If there is an error, return an error.
func (c *ConfigurationMap) Length(key string) (int, error) {
	_, ok := c.data[key]
	if ok {
		return len(c.data[key]), nil
	} else {
		return 0, fmt.Errorf("key not found")
	}
}
