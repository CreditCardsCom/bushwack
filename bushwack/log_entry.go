package bushwack

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type LogEntries []LogEntry

type LogEntry struct {
	Protocol         string `json:"protocol"`
	Timestamp        string `json:"@timestamp"`
	LoadBalancer     string `json:"load_balancer"`
	RemoteAddress    string `json:"remote_address"`
	TargetAddress    string `json:"target_address"`
	StatusCode       int    `json:"status_code"`
	TargetStatusCode int    `json:"target_status_code"`
	Method           string `json:"method"`
	Url              string `json:"url"`
	UserAgent        string `json:"user_agent"`
	SslCipher        string `json:"ssl_cipher"`
	SslProtocol      string `json:"ssl_protocol"`
	TargetGroup      string `json:"target_group"`
	TargetGroupArn   string `json:"target_group_arn"`
}

func NewLogEntries() LogEntries {
	return make(LogEntries, 0)
}

func (entries *LogEntries) PushEntry(args []string) {
	sc := parseInt(args[7])
	tsc := parseInt(args[8])
	p := normalizeProtocol(args[0])
	lb := normalizeTripleSlash(args[2])
	tg := normalizeTripleSlash(args[16])
	method, url := splitRequest(args[12])
	e := LogEntry{
		Protocol:         p,
		Timestamp:        args[1],
		LoadBalancer:     lb,
		RemoteAddress:    args[3],
		TargetAddress:    args[4],
		StatusCode:       sc,
		TargetStatusCode: tsc,
		Method:           method,
		Url:              url,
		UserAgent:        args[13],
		SslCipher:        args[14],
		SslProtocol:      args[15],
		TargetGroup:      tg,
		TargetGroupArn:   args[16],
	}

	*entries = append(*entries, e)
}

func (entries LogEntries) SerializeBulkBody() (string, error) {
	lines := make([]string, len(entries)*2)

	for _, e := range entries {
		j, err := json.Marshal(e)
		if err != nil {
			return "", err
		}

		index := parseIndexName(e.Timestamp)
		action := fmt.Sprintf("{\"index\": {\"_index\": \"%s\", \"_type\": \"alb-access-log\"}}", index)
		lines = append(lines, action)
		lines = append(lines, string(j))
	}

	return strings.Join(lines, "\n"), nil
}

// 2018-01-22T23:55:03.306727Z
func parseIndexName(ts string) string {
	var t time.Time

	err := t.UnmarshalText([]byte(ts))
	if err != nil {
		t = time.Now()
	}

	return fmt.Sprintf("logstash-%d.%02d.%02d", t.Year(), t.Month(), t.Day())
}

func parseInt(i string) int {
	ret, err := strconv.ParseInt(i, 10, 0)
	if err != nil {
		ret = -1
	}

	return int(ret)
}

func normalizeProtocol(p string) string {
	switch p {
	case "h2":
		return "http2"
	default:
		return p
	}
}

func splitRequest(r string) (string, string) {
	parts := strings.Split(r, " ")

	// Just return the original value if we don't have 3 parts
	if len(parts) < 2 {
		return "", r
	}

	return parts[0], parts[1]
}

func normalizeTripleSlash(r string) string {
	parts := strings.Split(r, "/")

	if len(parts) < 2 {
		return r
	}

	return parts[1]
}
