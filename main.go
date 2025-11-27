package main

import (
	"flag"
	"fmt"
	"io"
	"maps"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
)

const repos = "Code/"

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
		"The name to use for the project, if blank use repo name")

	if parseErr := flags.Parse(args[1:]); parseErr != nil {
		return in, parseErr
	}

	in.link = flags.Arg(0)

	in.repo, err = extractRepo(in.link, in.name)
	if err != nil {
		return in, err
	}

	in.path = joinPath(in)

	return in, nil
}

func extractRepo(link, name string) (string, error) {
	var prefixReplacements = map[string]string{
		"ssh://git@codeberg.org/":  "git@codeberg.org:",
		"https://codeberg.org/":    "git@codeberg.org:",
		"https://github.com/":      "git@github.com:",
		"https://gist.github.com/": "git@gist.github.com:",
	}

	var knownPrefixes = slices.Collect(maps.Values(prefixReplacements))

	repo := link

	for old, prefix := range prefixReplacements {
		repo = strings.Replace(repo, old, prefix, 1)
	}

	if !hasPrefix(repo, knownPrefixes) {
		return "", fmt.Errorf("Missing known prefix")
	}

	gist := strings.Contains(repo, "gist.github.com")

	for _, prefix := range knownPrefixes {
		repo = strings.TrimPrefix(repo, prefix)
	}

	const suffix = ".git"

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

func joinPath(in input) string {
	return filepath.Join(in.base, serviceName(in.link), in.repo)
}

func serviceName(link string) string {
	s := strings.ToLower(link)

	switch {
	case strings.Contains(s, "github.com"):
		return "GitHub"
	case strings.Contains(s, "codeberg.org"):
		return "Codeberg"
	default:
		return ""
	}
}

func hasPrefix(s string, prefixes []string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(s, p) {
			return true
		}
	}

	return false
}
