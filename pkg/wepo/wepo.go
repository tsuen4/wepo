package wepo

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/tsuen4/wepo/pkg/wepo/config"
	"golang.org/x/term"
)

// Wepo structure provide the client. Wepo holds thw config.
type Wepo struct {
	cfg *config.WepoConfig
}

// New returns a Wepo client. Requires '{cfgDirPath}/config.ini'.
func New(cfgDirPath string) (*Wepo, error) {
	cfg, err := config.New(cfgDirPath)
	if err != nil {
		return nil, err
	}

	return &Wepo{
		cfg: cfg,
	}, nil
}

// Input returns a string. The string is entered from an argument or pipeline.
func Input(args []string, fd int) (string, error) {
	var c string

	// fd: 0 -> default
	if term.IsTerminal(fd) {
		c = strings.Join(args, " ")
	} else {
		cBytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", err
		}
		c = string(cBytes)
	}

	if len(c) == 0 {
		return "", fmt.Errorf("empty value")
	}
	return c, nil
}

// PostDiscord sends content to Wepoconfig.URL
func (w Wepo) PostDiscord(content string) error {
	body := strings.ReplaceAll(w.cfg.Payload, "{input}", content)
	resp, err := http.Post(w.cfg.URL, "application/json", bytes.NewBuffer([]byte(body)))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("status code err: got: %d, want: %d", resp.StatusCode, http.StatusNoContent)
	}
	return nil
}

var errNoNeedSplit = fmt.Errorf("no need to split")

// NewContent is create post contents from string
func (w Wepo) NewContents(input string) ([]string, error) {
	lines := strings.Split(input, "\n")
	limit := w.cfg.CharLimit

	contents := []string{}
	for i, line := range lines {
		// when last line has an empty value
		if i == len(lines)-1 && len(line) == 0 {
			break
		}

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
		body += lines[idx] + "\n"
	}
	body += str

	if len([]rune(body)) > limit {
		lines = append(lines, str)
	} else {
		lines[idx] = body
	}
	return lines, nil
}