package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/go-github/github"
	"github.com/hhatto/gocloc"
	"golang.org/x/oauth2"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

// GitHub
type GitHub struct {
	token    string
	user     string
	endpoint string
	org      string
}

// Analysis
type Analysis struct {
	repos []string
	paths []string
	dir   string
	sync.RWMutex
	github *GitHub
	wg     sync.WaitGroup
}

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

	a.Do(g)
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

func (a *Analysis) Do(g *GitHub) error {
	var err error
	a.dir, err = ioutil.TempDir("", "clog-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(a.dir)

	a.wg = sync.WaitGroup{}
	for _, repo := range a.repos {
		a.wg.Add(1)
		go a.clone(repo)
	}

	a.wg.Wait()

	langs := gocloc.NewDefinedLanguages()
	opts := gocloc.NewClocOptions()
	p := gocloc.NewProcessor(langs, opts)
	res, err := p.Analyze(a.paths)
	if err != nil {
		return err
	}

	fmt.Printf("%#v\n", res)
	return nil
}

func (a *Analysis) clone(URL string) {
	defer a.wg.Done()

	name := filepath.Base(URL)
	path := fmt.Sprintf("%s/%s", a.dir, name)

	_, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:   URL,
		Depth: 1,
		Auth: &http.BasicAuth{
			Username: a.github.user,
			Password: a.github.token,
		},
		Progress: os.Stdout,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("git-clone error(%s): %s\n", URL, err))
		return
	}

	a.Lock()
	a.paths = append(a.paths, path)
	a.Unlock()
}
