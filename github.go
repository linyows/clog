package main

import (
	"context"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// GitHub
type GitHub struct {
	token    string
	user     string
	endpoint string
	org      string
}

func (g *GitHub) client(ctx context.Context) (*github.Client, error) {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: g.token},
	)
	tc := oauth2.NewClient(ctx, ts)
	var client *github.Client

	if g.endpoint == "" {
		client = github.NewClient(tc)
	} else {
		var err error
		client, err = github.NewEnterpriseClient(g.endpoint, "", tc)
		if err != nil {
			return nil, err
		}
	}

	return client, nil
}

func (g *GitHub) fetch(ctx context.Context) ([]string, error) {
	client, err := g.client(ctx)
	if err != nil {
		return nil, err
	}

	opt := &github.RepositoryListByOrgOptions{"all", github.ListOptions{Page: 1, PerPage: 100}}
	got, _, err := client.Repositories.ListByOrg(ctx, g.org, opt)
	if err != nil {
		return nil, err
	}

	repos := []string{}

	for _, v := range got {
		repos = append(repos, v.GetCloneURL())
	}

	return repos, nil
}
