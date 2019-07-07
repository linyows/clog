package main

import (
	"context"
	"fmt"
	"os"
)

func main() {
	g := &GitHub{
		endpoint: os.Getenv("GITHUB_ENDPOINT"),
		user:     os.Getenv("GITHUB_USER"),
		token:    os.Getenv("GITHUB_TOKEN"),
		org:      os.Getenv("GITHUB_ORG"),
	}

	if g.token == "" {
		fmt.Fprintf(os.Stderr, "GITHUB_TOKEN is required\n")
		return
	}

	if g.org == "" {
		fmt.Fprintf(os.Stderr, "GITHUB_ORG is required\n")
		return
	}

	ctx := context.Background()
	a := &Analysis{github: g}
	var err error
	a.repos, err = g.fetch(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("Error: %s\n", err))
		return
	}

	err = a.Do()
	if err != nil {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("Error: %s\n", err))
		return
	}

	fmt.Fprintf(os.Stderr, fmt.Sprintf("analyzed:  %d/%d\n", len(a.clocs), len(a.repos)))

	r := &Report{analysis: a}
	err = r.loadData()
	if err != nil {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("Error: %s\n", err))
		return
	}

	r.show()
}
