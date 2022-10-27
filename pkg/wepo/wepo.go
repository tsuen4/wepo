package wepo

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/tsuen4/wepo/pkg/wepo/config"
)

// wepo structure provide the client. wepo holds the config.
type wepo struct {
	cfg *config.WepoConfig
}

// New returns a wepo client. Requires 'config.ini'.
func New(iniPath, section string) (*wepo, error) {
	cfg, err := config.New(iniPath, section)
	if err != nil {
		return nil, err
	}

	return &wepo{
		cfg: cfg,
	}, nil
}

// ErrEmptyValue : error with when empty arguments
var ErrEmptyValue = fmt.Errorf("empty value")

// Input returns a string. The string is entered from an argument or pipeline.
func Input(args []string) (string, error) {
	bytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", err
	}
	input := string(bytes)

	if len(input) == 0 {
		return "", ErrEmptyValue
	}

	return input, nil
}

// PostContents sends content to wepo.WepoConfig.URL
func (w wepo) PostContents(contents []string) error {
	for _, c := range contents {
		if err := w.post(c); err != nil {
			return err
		}
	}
	return nil
}

func (w wepo) post(content string) error {
	body := strings.ReplaceAll(w.cfg.Payload, "{input}", content)
	resp, err := http.Post(w.cfg.URL, "application/json", bytes.NewBuffer([]byte(body)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		respBody := string(b)
		statusCodeError := fmt.Sprintf("status code err: got: %d, want: %d", resp.StatusCode, http.StatusNoContent)
		return fmt.Errorf("%s\n%s", respBody, statusCodeError)
	}
	return nil
}

var errNoNeedSplit = fmt.Errorf("no need to split")

// NewContent is create post contents from string
func (w wepo) NewContents(input string) ([]string, error) {
	lines := strings.Split(input, "\n")
	limit := w.cfg.CharLimit

	contents := []string{}
	for i, line := range lines {
		// when last line has an empty value
		if i == len(lines)-1 && len(line) == 0 {
			break
		}

		// escape double quotes
		line = strings.ReplaceAll(line, "\"", `\"`)

		var err error
		contents, err = appendLine(contents, line, limit)
		if err != nil {
			return nil, err
		}
	}
	return contents, nil
}

func splitRow(str string, limit int) ([]string, error) {
	runes := []rune(str)
	lines := []string{}

	if len(runes) > limit {
		for i := 0; i < len(runes); i += limit {
			nextSep := i + limit

			isBackSlashEnd := false
			if nextSep < len(runes) {
				// Prevents trailing backslashes
				if runes[nextSep-1] == '\\' {
					isBackSlashEnd = true
					nextSep--
				}

				lines = append(lines, string(runes[i:nextSep]))
			} else {
				lines = append(lines, string(runes[i:]))
			}

			if isBackSlashEnd {
				i--
			}
		}
		return lines, nil
	} else {
		return nil, errNoNeedSplit
	}
}

func appendLine(lines []string, str string, limit int) ([]string, error) {
	if len([]rune(str)) > limit {
		splitted, err := splitRow(str, limit)
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
		body += lines[idx] + `\n`
	}
	body += str

	if len([]rune(body)) > limit {
		lines = append(lines, str)
	} else {
		lines[idx] = body
	}
	return lines, nil
}
