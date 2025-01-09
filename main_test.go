package main

import "testing"

func TestExtractRepo(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		for _, tt := range []struct {
			uri  string
			gist string
			repo string
		}{
			{
				"git@github.com:peterhellberg/w4-init.git",
				"",
				"peterhellberg/w4-init",
			},
			{
				"git@github.com:tsoding/zigout.git",
				"",
				"tsoding/zigout",
			},
			{
				"git@github.com:tsoding/zigout.git",
				"zigbreak",
				"tsoding/zigbreak",
			},
			{
				"git@gist.github.com:124fb025981b9b167f942c39b87e6624.git",
				"",
				"Gists/124fb025981b9b167f942c39b87e6624",
			},
			{
				"git@gist.github.com:124fb025981b9b167f942c39b87e6624.git",
				"aisnake",
				"Gists/aisnake",
			},
			{
				"https://gist.github.com/8ab356001fea59cdc26d01d130bff17b.git",
				"",
				"Gists/8ab356001fea59cdc26d01d130bff17b",
			},
			{
				"https://github.com/c7/recept.c7.se.git",
				"",
				"c7/recept.c7.se",
			},
		} {
			t.Run(tt.uri, func(t *testing.T) {
				repo, err := extractRepo(tt.uri, tt.gist)
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}

				if got, want := repo, tt.repo; got != want {
					t.Fatalf("repo = %q, want %q", got, want)
				}
			})
		}
	})
}
