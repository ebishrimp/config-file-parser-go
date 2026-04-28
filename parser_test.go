package confparser

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	input := `
# This is a comment
key1 value1
key2 value2 # This is another comment

key3 value3-1 value3-2

key4 8080
`
	r := strings.NewReader(input)
	conf, err := Parse(r)

	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if get := conf.GetValue("key1"); get != "value1" {
		t.Errorf("Expected 'value1' for key1, got '%s'", get)
	}

	if get := conf.GetValue("key2"); get != "value2" {
		t.Errorf("Expected 'value2' for key2, got '%s'", get)
	}

	if get := conf.GetValue("key3"); get != "" {
		t.Errorf("Expected no value for key3, got '%s'", get)
	}

	if get, err := conf.IntGetValue("key4"); err != nil || get != 8080 {
		t.Errorf("Expected 8080 for key4, got '%d'", get)
	}

	if conf.ExistsValue("key3") {
		t.Errorf("Expected no value for key3, but it exists")
	}

	if !conf.ExistsValue("key4") {
		t.Errorf("Expected value for key4, but it does not exist")
	}
}
