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

multiplekey 1
multiplekey 2
multiplekey 3

key5 value5

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

	r2 := strings.NewReader(input)

	multiConf, err := ParseMultipleValues(r2)
	if err != nil {
		t.Fatalf("ParseMultipleValues failed: %v", err)
	}

	if get := multiConf.GetMultipleValues("multiplekey"); len(get) != 3 || get[0] != "1" || get[1] != "2" || get[2] != "3" {
		t.Errorf("Expected ['1', '2', '3'] for multiplekey, got %v", get)
	}

	if get := multiConf.GetMultipleValues("key5"); len(get) != 1 || get[0] != "value5" {
		t.Errorf("Expected ['value5'] for key5, got %v", get)
	}

	if get := multiConf.GetFirstValue("multiplekey"); get != "1" {
		t.Errorf("Expected '1' for first value of multiplekey, got '%s'", get)
	}

	if !multiConf.ExistsMultipleValues("multiplekey") {
		t.Errorf("Expected multiple values for multiplekey, but none found")
	}

	input2_forConfigurationMap := `
key1 value1
key1 value2
key2 v1 v2 v3

key3 "str with spaces"
key4 "str contains " quotes"
key5 "str with # in quotes"

key6 123
key7 123abc

key8 1.23
key9 1.23abc
# This is a comment

# Empty line
key10
# no value for key10

key11 true
key12 false
key13 yes
key14 no

key15 v1
key15 v2
key15 v3
key15 v4

key16 value # with comment
`

	r3 := strings.NewReader(input2_forConfigurationMap)
	c, _ := ParseConfig(r3)

	if k := c.StringSlice("key1"); k[0] != "value1" || k[1] != "value2" {
		t.Errorf("Expected 'value1' and 'value2' for key1, got %v", c.StringSlice("key1"))
	}

	if c.Exists("key2") {
		t.Errorf("Expected no value for key2, but it exists")
	}

	if c.String("key3") != "str with spaces" {
		t.Errorf("Expected 'str with spaces' for key3, got '%s'", c.String("key3"))
	}

	if c.String("key4") != "str contains \" quotes" {
		t.Errorf("Expected 'str contains \" quotes' for key4, got '%s'", c.String("key4"))
	}

	if c.String("key5") != "str with # in quotes" {
		t.Errorf("Expected 'str with # in quotes' for key5, got '%s'", c.String("key5"))
	}

	if i, err := c.Int("key6"); err != nil || i != 123 {
		t.Errorf("Expected 123 for key6, got '%d'", i)
	}

	if _, err := c.Int("key7"); err == nil {
		t.Errorf("Expected error for key7, but got none")
	}

	if f, err := c.Float("key8"); err != nil || f != 1.23 {
		t.Errorf("Expected 1.23 for key8, got '%f'", f)
	}

	if _, err := c.Float("key9"); err == nil {
		t.Errorf("Expected error for key9, but got none")
	}

	if c.Exists("key10") {
		t.Errorf("Expected no value for key10, but it exists")
	}

	if b, err := c.Bool("key11"); err != nil || b != true {
		t.Errorf("Expected true for key11, got '%t'", b)
	}

	if b, err := c.Bool("key12"); err != nil || b != false {
		t.Errorf("Expected false for key12, got '%t'", b)
	}

	if b, err := c.Bool("key13"); err != nil || b != true {
		t.Errorf("Expected true for key13, got '%t'", b)
	}

	if b, err := c.Bool("key14"); err != nil || b != false {
		t.Errorf("Expected false for key14, got '%t'", b)
	}

	if l, err := c.Length("key15"); err != nil || l != 4 {
		t.Errorf("Expected length 4 for key15, got '%d'", l)
	}

	if c.Exists("key16") {
		t.Errorf("Expected no value for key16, but it exists")
	}
}
