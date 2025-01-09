package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const repos = "Code/GitHub"

type input struct {
	base string
	repo string
	name string
	link string
	path string
}

func parseInput(args []string) (input, error) {
	var in input

	flags := flag.NewFlagSet(args[0], flag.ExitOnError)

	home, err := os.UserHomeDir()
	if err != nil {
		return in, err
	}

	defaultBase := filepath.Join(home, repos)

	flags.StringVar(&in.base, "base", defaultBase,
		"The base directory to clone into")

	flags.StringVar(&in.name, "name", "",
		"The name to use for the project, if blank use GitHub name")

	if err := flags.Parse(args[1:]); err != nil {
		return in, err
	}

	in.link = flags.Arg(0)

	in.repo, err = extractRepo(in.link, in.name)
	if err != nil {
		return in, err
	}

	in.path = filepath.Join(in.base, in.repo)

	return in, nil
}

func extractRepo(link, name string) (string, error) {
	const (
		gistPrefix = "git@gist.github.com:"
		ghPrefix   = "git@github.com:"
		suffix     = ".git"
	)

	repo := link

	repo = strings.Replace(repo, "https://github.com/", ghPrefix, 1)
	repo = strings.Replace(repo, "https://gist.github.com/", gistPrefix, 1)

	gist := strings.HasPrefix(repo, gistPrefix)

	if !strings.HasPrefix(repo, ghPrefix) && !gist {
		return "", fmt.Errorf("Missing %q prefix", ghPrefix)
	}

	repo = strings.TrimPrefix(repo, gistPrefix)
	repo = strings.TrimPrefix(repo, ghPrefix)

	if !strings.HasSuffix(link, suffix) {
		return "", fmt.Errorf("Missing %q suffix", suffix)
	}

	repo = strings.TrimSuffix(repo, suffix)

	var account, project string

	if gist {
		account = "Gists"
		project = repo
	} else {
		var found bool

		account, project, found = strings.Cut(repo, "/")
		if !found {
			return "", fmt.Errorf("Could not override repo name")
		}
	}

	if name != "" {
		project = name
	}

	return filepath.Join(account, project), nil
}

func main() {
	if err := run(os.Args, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(args []string, stderr io.Writer) error {
	in, err := parseInput(args)
	if err != nil {
		return err
	}

	cmd := exec.Command("git", "clone", in.link, in.path)

	cmd.Stderr = stderr

	return cmd.Run()
}
