package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/tsuen4/wepo/pkg/wepo"
)

type exitCode int

// error code
const (
	exitCodeOK exitCode = iota
	exitCodeErr
)

func main() {
	flag.Parse()
	os.Exit(int(Main(flag.Args())))
}

func Main(args []string) exitCode {
	if err := run(args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		return exitCodeErr
	}
	return exitCodeOK
}

func run(args []string) error {
	exe, err := os.Executable()
	if err != nil {
		return err
	}

	cfg, err := wepo.Config(filepath.Join(filepath.Dir(exe), wepo.CfgFileName))
	if err != nil {
		return err
	}

	input, err := wepo.Input(args, int(os.Stdin.Fd()))
	if err != nil {
		return err
	}
	contents, err := cfg.Contents(input)
	if err != nil {
		return err
	}

	for _, c := range contents {
		if err := cfg.PostDiscord(c); err != nil {
			return err
		}
	}

	return nil
}
