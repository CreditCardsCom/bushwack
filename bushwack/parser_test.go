package bushwack

import (
	"strings"
	"testing"
)

const fixture = `
http 2018-01-22T23:55:03.306727Z app/loadbalancer-production/1f323822c01b6275 1.1.1.1:33312 10.20.21.202:32771 0.001 0.189 0.000 302 302 446 338 "GET http://example.com:80/index.html HTTP/1.1" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.10; rv:41.0) Gecko/20100101 Firefox/55.0" - - arn:aws:elasticloadbalancing:us-west-2:123456789012:targetgroup/targetgroup-production/1de8c5df30d77023 "Root=1-5a6679d7-460d74875250ee5a6778f85c" "-" "-"
h2 2018-01-22T23:55:03.688517Z app/loadbalancer-production/1f323822c01b6275 1.1.1.1:58705 10.20.20.112:32768 0.000 0.197 0.000 302 302 2196 288 "GET https://blog.example.com:443/?testing=asdf HTTP/2.0" "Mozilla/5.0 (iPad; CPU OS 11_2_1 like Mac OS X) AppleWebKit/604.4.7 (KHTML, like Gecko) Version/11.0 Mobile/15C153 Safari/604.1" ECDHE-RSA-AES128-GCM-SHA256 TLSv1.2 arn:aws:elasticloadbalancing:us-west-2:123456789012:targetgroup/targetgroup-production/1de8c5df30d77023 "Root=1-5a6679d7-718ffabe05a3580f2696ac33" "example.com" "session-renegotiated-or-reused"
`

var fixtures []string

func init() {
	t := strings.Trim(fixture, "\n")
	temp := strings.Split(t, "\n")

	for i := 0; i < 500; i++ {
		fixtures = append(fixtures, temp...)
	}
}

func TestParseLog(t *testing.T) {
	log := strings.Join(fixtures, "\n")

	entries, err := ParseLog(log)
	if err != nil {
		t.Fatalf("Error '%s' not expected.", err)
	}

	if len(entries) != 1000 {
		t.Fatalf("Expected 1000 entries got, %d.", len(entries))
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
	log := strings.Join(fixtures, "\n")

	for i := 0; i < b.N; i++ {
		ParseLog(log)
	}
}

func BenchmarkParseLine(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, line := range fixtures {
			parseLine(line)
		}
	}
}
