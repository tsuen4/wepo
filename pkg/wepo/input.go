package wepo

import (
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"
)

// Input from arguments or pipeline
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
