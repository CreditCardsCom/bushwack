package bushwack

import (
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"
)

var fixture string

func init() {
	p := filepath.Join("testdata", "sample-log.log")
	bytes, err := ioutil.ReadFile(p)

	if err != nil {
		panic(err)
	}

	var file []byte
	for i := 0; i < 100; i++ {
		file = append(file, bytes...)
	}

	fixture = string(file)
}

func TestProcessLog(t *testing.T) {
	p := filepath.Join("testdata", "sample-log.log.gz")
	num, body, err := ProcessLog(p)
	if err != nil {
		t.Fatalf("Error '%s' not expected.", err)
	}

	if num != 4 {
		t.Fatalf("Expected 4 entries got, %d.", num)
	}

	if !strings.HasSuffix(body, "\n") {
		t.Fatalf("Expected bulk body to contain trailing newline character.")
	}
}

func TestParseLog(t *testing.T) {
	entries, err := ParseLog(fixture)
	if err != nil {
		t.Fatalf("Error '%s' not expected.", err)
	}

	if len(entries) != 400 {
		t.Fatalf("Expected 400 entries got, %d.", len(entries))
	}

	f := `failure to parse`
	_, err = ParseLog(f)
	if err != InvalidLogFormat {
		t.Fatalf("Expected error of type 'InvalidLogFormat`, but got '%s'.", err)
	}

	if err == nil {
		t.Fatalf("Expected error to parse invalid log format.")
	}
}

func BenchmarkParseLog(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ParseLog(fixture)
	}
}

func BenchmarkParseLine(b *testing.B) {
	fixtures := strings.Split(fixture, "\n")

	for i := 0; i < b.N; i++ {
		for _, line := range fixtures {
			parseLine(line)
		}
	}
}
