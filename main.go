package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/tsuen4/wepo/internal/tui"
	"github.com/tsuen4/wepo/pkg/wepo"
)

var isTUIMode bool

func init() {
	flag.BoolVar(&isTUIMode, "t", false, "Enable tui mode")
}

func main() {
	flag.Parse()

	if err := run(flag.Args()); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	cfgDirPath := filepath.Join(filepath.Dir(exe))

	var runWepo func(string, []string) error
	if isTUIMode {
		runWepo = tuiMode
	} else {
		runWepo = shellMode
	}

	if err := runWepo(cfgDirPath, flag.Args()); err != nil {
		return err
	}

	return nil
}

func shellMode(cfgDirPath string, args []string) error {
	client, err := wepo.New(cfgDirPath)
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

func tuiMode(cfgDirPath string, args []string) error {
	if err := tui.Run(cfgDirPath, args); err != nil {
		return err
	}

	return nil
}
