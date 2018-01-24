package bushwack

import (
	"bufio"
	"bytes"
	"errors"
	"strings"
)

var InvalidLogFormat = errors.New("Invalid log format, expecting 20 fields")

func ParseLog(src string) (LogEntries, error) {
	logs := NewLogEntries()

	for _, line := range strings.Split(src, "\n") {
		if line != "" {
			e, err := parseLine(line)
			if err != nil {
				return nil, err
			}

			logs.PushEntry(e)
		}
	}

	return logs, nil
}

func parseLine(line string) ([]string, error) {
	r := strings.NewReader(line)
	scanner := bufio.NewScanner(r)
	scanner.Split(splitOnSpaceOrQuotes)
	var args []string

	for scanner.Scan() {
		args = append(args, scanner.Text())
	}

	if len(args) != 20 {
		return nil, InvalidLogFormat
	}

	return args, nil
}

func splitOnSpaceOrQuotes(data []byte, atEOF bool) (int, []byte, error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	i, word, err := bufio.ScanWords(data, atEOF)

	if i != 0 && word[0] == '"' {
		// The whole word is quoted
		if l := len(word); word[l-1] == '"' {
			return i, word[1 : l-1], nil
		}

		// Iterate through the remaining data searching for the end quote
		for j := i; j < len(data); j++ {
			if data[j] == '"' {
				trim := bytes.TrimSpace(data[0:j])
				trim = bytes.Trim(trim, "\"")

				return j + 1, trim, nil
			}
		}
	}

	// Return bufio.ScanWords output if we haven't found a quote
	return i, word, err
}
