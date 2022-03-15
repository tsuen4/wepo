package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/tsuen4/wepo/content"
	"gopkg.in/ini.v1"
)

type exitCode int

// error code
const (
	exitCodeOK exitCode = iota
	exitCodeErr
)

// config
const (
	cfgFileName = "config.ini"
	cfgURLKey   = "webhook_url"
)

// config key
var (
	addr = ""
)

func init() {
	flag.StringVar(&addr, "a", "", fmt.Sprintf(`Key with %s set in "%s"`, cfgURLKey, cfgFileName))
}

func main() {
	os.Exit(int(Main(os.Args[1:])))
}

func Main(args []string) exitCode {
	if err := run(args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		return exitCodeErr
	}
	return exitCodeOK
}

const sendCharLimit = 1024

func run(args []string) error {
	flag.Parse()
	url, err := getURL()
	if err != nil {
		return err
	}

	input, err := content.Input(args, int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	contents, err := content.New(input, sendCharLimit)
	if err != nil {
		return err
	}

	for _, c := range contents {
		if err := postDiscord(url, c); err != nil {
			return err
		}
	}

	return nil
}

func getURL() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}

	cfg, err := ini.Load(filepath.Join(filepath.Dir(exe), cfgFileName))
	if err != nil {
		return "", err
	}

	url := cfg.Section(addr).Key(cfgURLKey).String()
	if len(url) == 0 {
		msg := fmt.Sprintf(`"%s" is not set in "%s"`, cfgURLKey, cfgFileName)
		if len(addr) != 0 {
			msg = fmt.Sprintf("[%s] %s", addr, msg)
		}
		return "", fmt.Errorf(msg)
	}

	return url, nil
}

func postDiscord(url, content string) error {
	body := struct {
		Content string `json:"content"`
	}{
		Content: content,
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("status code err: got: %d, want: %d", resp.StatusCode, http.StatusNoContent)
	}
	return nil
}
