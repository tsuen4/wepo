package wepo

import (
	"fmt"
	"strings"
)

var errNoNeedSplit = fmt.Errorf("no need to split")

// Content is create post contents from string array
func (cfg wepoConfig) Contents(input string) ([]string, error) {
	lines := strings.Split(input, "\n")
	divNum := cfg.CharLimit

	contents := []string{}
	for i, line := range lines {
		// when last line has an empty value
		if i == len(lines)-1 && len(line) == 0 {
			break
		}

		var err error
		contents, err = appendLine(contents, line, divNum)
		if err != nil {
			return nil, err
		}
	}
	return contents, nil
}

func splitRow(str string, divNum int) ([]string, error) {
	runes := []rune(str)
	lines := []string{}

	if len(runes) > divNum {
		for i := 0; i < len(runes); i += divNum {
			nextSep := i + divNum

			if nextSep < len(runes) {
				lines = append(lines, string(runes[i:nextSep]))
			} else {
				lines = append(lines, string(runes[i:]))
			}
		}
		return lines, nil
	} else {
		return nil, errNoNeedSplit
	}
}

func appendLine(lines []string, str string, divNum int) ([]string, error) {
	if len([]rune(str)) > divNum {
		splitted, err := splitRow(str, divNum)
		if err != nil {
			return nil, err
		}
		return append(lines, splitted...), nil
	}

	// initialize
	idx := 0
	if len(lines) == 0 {
		lines = append(lines, "")
	} else {
		idx = len(lines) - 1
	}

	body := ""
	if len(lines[idx]) != 0 {
		body += lines[idx] + "\n"
	}
	body += str

	if len([]rune(body)) > divNum {
		lines = append(lines, str)
	} else {
		lines[idx] = body
	}
	return lines, nil
}
