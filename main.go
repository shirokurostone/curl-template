package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	var flags string
	flag.StringVar(&flags, "flags", "", "curl additional flags")
	flag.StringVar(&flags, "f", "", "curl additional flags (short)")

	var evalShellCommand bool
	flag.BoolVar(&evalShellCommand, "expand-shell-command", false, "expand shell commands")
	flag.BoolVar(&evalShellCommand, "s", false, "expand shell commands (short)")

	var expandEnv bool
	flag.BoolVar(&expandEnv, "expand-env", false, "expand environment variables")
	flag.BoolVar(&expandEnv, "e", false, "expand environment variables (short)")

	var prettyPrint bool
	flag.BoolVar(&prettyPrint, "prettyprint", false, "pretty print")
	flag.BoolVar(&prettyPrint, "p", false, "prett yprint (short)")

	flag.Parse()

	if evalShellCommand && expandEnv {
		flag.PrintDefaults()
		os.Exit(0)
	}

	rand.Seed(time.Now().UnixNano())

	args := flag.Args()
	var filename string
	if len(args) == 0 {
		filename = "/dev/stdin"
	} else if len(args) == 1 {
		filename = args[0]
	} else {
		flag.PrintDefaults()
		os.Exit(0)
	}

	err := Run(filename, flags, evalShellCommand, expandEnv, prettyPrint)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func Run(filename string, flags string, evalShellCommand bool, expandEnv bool, prettyPrint bool) error {

	req, err := OpenHttpFile(filename)
	if err != nil {
		return err
	}

	if evalShellCommand {
		req.ExpandShell(context.Background())
	} else if expandEnv {
		req.ExpandEnv()
	}

	if err != nil {
		return err
	}

	fmt.Println(req.CurlCommand(prettyPrint, flags))
	return nil
}

func OpenHttpFile(filename string) (*HttpRequest, error) {

	req, err := func(filename string) (*HttpRequest, error) {
		fp, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer fp.Close()

		return NewHttpRequestJson(fp)
	}(filename)
	if req != nil {
		return req, nil
	}

	req, err = func(filename string) (*HttpRequest, error) {
		fp, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		defer fp.Close()

		return NewHttpRequest(fp)
	}(filename)

	return req, err

}
