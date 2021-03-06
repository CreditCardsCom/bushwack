package bushwack

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

var InvalidLogFormat = errors.New("Invalid log format, expecting >=20 fields")
var ClosingQuoteNotFound = errors.New("Closing quote not found in line")

func ProcessLog(filename string) (int, string, error) {
	contents, err := decompress(filename)
	if err != nil {
		return 0, "", err
	}

	entries, err := ParseLog(string(contents))
	if err != nil {
		return 0, "", err
	}

	num := len(entries)
	if num == 0 {
		return 0, "", nil
	}

	body, err := entries.SerializeBulkBody()
	if err != nil {
		return 0, "", err
	}

	return num, body, nil
}

func ParseLog(src string) (LogEntries, error) {
	logs := NewLogEntries()

	for i, line := range strings.Split(src, "\n") {
		if line != "" {
			e, err := parseLine(line)
			if err != nil {
				// Discard the line, but continue with the log
				if err == InvalidLogFormat || err == ClosingQuoteNotFound {
					log.Printf("Error parsing line %d: %s\n", i+1, err)
					continue
				}

				return nil, err
			}

			logs.PushEntry(e)
		}
	}

	return logs, nil
}

func parseLine(line string) ([]string, error) {
	r := strings.NewReader(line)
	l := len(line)
	scanner := bufio.NewScanner(r)
	scanner.Split(splitOnSpaceOrQuotes)
	// Override the max token size to the string length
	scanner.Buffer(make([]byte, l), l)
	var args []string

	for scanner.Scan() {
		args = append(args, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(args) < 20 {
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
		if l := len(word); l > 1 && word[l-1] == '"' {
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

		return 0, nil, ClosingQuoteNotFound
	}

	// Return bufio.ScanWords output if we haven't found a quote
	return i, word, err
}

func decompress(f string) ([]byte, error) {
	fd, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	r, err := gzip.NewReader(fd)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	return ioutil.ReadAll(r)
}
