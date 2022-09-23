package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/tsuen4/wepo/internal/tui"
	"github.com/tsuen4/wepo/pkg/wepo"
	"github.com/tsuen4/wepo/pkg/wepo/config"
)

const appName = "wepo"

var (
	isTUIMode bool
	section   string
)

func init() {
	flag.BoolVar(&isTUIMode, "t", false, "Enable tui mode")
	flag.StringVar(&section, "s", "", fmt.Sprintf(`Section name of "%s" where "%s" is set`, cfgFilePath(), config.WebhookURLKey))
}

func main() {
	flag.Parse()

	if err := run(flag.Args()); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	var runWepo func(string, string, []string) error
	if isTUIMode {
		runWepo = tuiMode
	} else {
		runWepo = shellMode
	}

	if err := runWepo(cfgFilePath(), section, flag.Args()); err != nil {
		return err
	}

	return nil
}

func shellMode(cfgDirPath, section string, args []string) error {
	client, err := wepo.New(cfgDirPath, section)
	if err != nil {
		return err
	}

	input, err := wepo.Input(args, int(os.Stdin.Fd()))
	if err != nil {
		return err
	}

	contents, err := client.NewContents(input)
	if err != nil {
		return err
	}

	if err := client.PostContents(contents); err != nil {
		return err
	}

	return nil
}

func tuiMode(cfgDirPath, section string, args []string) error {
	if err := tui.Run(cfgDirPath, section, args); err != nil {
		return err
	}

	return nil
}

// cfgFilePath returns config file path
func cfgFilePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "$HOME"
	}
	return filepath.Join(home, ".config", appName, config.CfgFileName)
}
