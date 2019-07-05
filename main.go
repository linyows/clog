package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
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
	repos  []string
	clocs  []string
	github *GitHub
	result *gocloc.Result
	wg     sync.WaitGroup
	sync.RWMutex
}

type Result struct {
	name string
	data gocloc.JSONLanguagesResult
}

type Report struct {
	results []Result
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

	err = a.Do()
	if err != nil {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("Error: %s\n", err))
		return
	}

	fmt.Fprintf(os.Stderr, fmt.Sprintf("analyzed:  %d/%d\n", len(a.clocs), len(a.repos)))

	r := &Report{}
	err = r.loadData()
	if err != nil {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("Error: %s\n", err))
		return
	}

	r.show()
}

func (r *Report) loadData() error {
	dir := "analyzed"
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		name := f.Name()
		bytes, _ := ioutil.ReadFile(filepath.Join(dir, name))
		var result gocloc.JSONLanguagesResult
		err := json.Unmarshal(bytes, &result)
		if err != nil {
			fmt.Fprintf(os.Stderr, fmt.Sprintf("Error: %s\n", err))
			continue
		}
		r.results = append(r.results, Result{
			name: strings.Replace(name, ".json", "", 1),
			data: result,
		})
	}

	return nil
}

func (r *Report) show() {
	header := "Language"
	commonHeader := "files          blank        comment           code"
	separator := "-------------------------------------------------------------------------" +
		"-------------------------------------------------------------------------" +
		"-------------------------------------------------------------------------"
	rowLen := 79
	headerLen := 28

	result := make(map[string]*gocloc.ClocLanguage)
	total := &gocloc.ClocLanguage{Name: "TOTAL"}

	for _, rr := range r.results {
		tt := rr.data.Total
		langs := rr.data.Languages

		//fmt.Printf("%s\n", rr.name)
		//fmt.Printf("%.[2]*[1]s\n", separator, rowLen)
		//fmt.Printf("%-[2]*[1]s %[3]s\n", header, headerLen, commonHeader)
		//fmt.Printf("%.[2]*[1]s\n", separator, rowLen)
		for _, l := range langs {
			//fmt.Printf("%-27v %6v %14v %14v %14v\n", l.Name, l.FilesCount, l.Blanks, l.Comments, l.Code)
			if _, ok := result[l.Name]; !ok {
				//fmt.Printf("%s, %d\n", l.Name, l.Code)
				result[l.Name] = &l
			} else {
				result[l.Name].FilesCount += l.FilesCount
				result[l.Name].Blanks += l.Blanks
				result[l.Name].Comments += l.Comments
				result[l.Name].Code += l.Code
			}
		}
		//fmt.Printf("%.[2]*[1]s\n", separator, rowLen)
		//fmt.Printf("%-27v %6v %14v %14v %14v\n", "TOTAL", tt.FilesCount, tt.Blanks, tt.Comments, tt.Code)
		//fmt.Printf("%.[2]*[1]s\n\n", separator, rowLen)

		total.FilesCount += tt.FilesCount
		total.Blanks += tt.Blanks
		total.Comments += tt.Comments
		total.Code += tt.Code
	}

	fmt.Printf("all\n")
	fmt.Printf("%.[2]*[1]s\n", separator, rowLen)
	fmt.Printf("%-[2]*[1]s %[3]s\n", header, headerLen, commonHeader)
	fmt.Printf("%.[2]*[1]s\n", separator, rowLen)
	for k, v := range result {
		fmt.Printf("%-27v %6v %14v %14v %14v\n", k, v.FilesCount, v.Blanks, v.Comments, v.Code)
	}
	fmt.Printf("%.[2]*[1]s\n", separator, rowLen)
	fmt.Printf("%-27v %6v %14v %14v %14v\n", "TOTAL", total.FilesCount, total.Blanks, total.Comments, total.Code)
	fmt.Printf("%.[2]*[1]s\n\n", separator, rowLen)
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

func (a *Analysis) Do() error {
	repoDir := "repos"
	if _, err := os.Stat(repoDir); os.IsNotExist(err) {
		err = os.Mkdir(repoDir, 0755)
		if err != nil {
			return err
		}
	}

	analyzedDir := "analyzed"
	if _, err := os.Stat(analyzedDir); os.IsNotExist(err) {
		err = os.Mkdir(analyzedDir, 0755)
		if err != nil {
			return err
		}
	}

	//defer os.RemoveAll(repos)

	a.wg = sync.WaitGroup{}
	for _, repo := range a.repos {
		a.wg.Add(1)
		go a.doClocBeforeClone(repo)
	}

	a.wg.Wait()

	return nil
}

func (a *Analysis) clone(URL string, path string) error {
	_, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:   URL,
		Depth: 1,
		Auth: &http.BasicAuth{
			Username: a.github.user,
			Password: a.github.token,
		},
		//Progress: os.Stderr,
		Progress: bytes.NewBuffer(nil),
	})

	return err
}

func (a *Analysis) doClocBeforeClone(URL string) {
	defer a.wg.Done()

	name := strings.Replace(filepath.Base(URL), ".git", "", 1)
	path := fmt.Sprintf("repos/%s", name)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := a.clone(URL, path)
		if err != nil {
			fmt.Fprintf(os.Stderr, fmt.Sprintf("clone error: %s\n%s\n", URL, err))
			return
		} else {
			fmt.Fprintf(os.Stderr, fmt.Sprintf("cloned %s\n", URL))
		}
	} else {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("clone skipped %s\n", URL))
	}

	jsonPath := "analyzed/" + name + ".json"
	if _, err := os.Stat(jsonPath); !os.IsNotExist(err) {
		return
	}

	langs := gocloc.NewDefinedLanguages()
	opts := gocloc.NewClocOptions()
	p := gocloc.NewProcessor(langs, opts)
	result, err := p.Analyze([]string{path})
	if err != nil {
		fmt.Printf("fail gocloc analyze. error: %v\n", err)
		return
	}

	var sortedLanguages gocloc.Languages
	for _, language := range result.Languages {
		if len(language.Files) != 0 {
			sortedLanguages = append(sortedLanguages, *language)
		}
	}
	sort.Sort(sortedLanguages)

	jsonResult := gocloc.NewJSONLanguagesResultFromCloc(result.Total, sortedLanguages)
	buf, err := json.Marshal(jsonResult)
	if err != nil {
		fmt.Printf("json marshal error: %v\n", err)
		return
	}

	err = ioutil.WriteFile(jsonPath, buf, 0644)
	if err != nil {
		fmt.Printf("file write error: %v\n", err)
		return
	}

	a.Lock()
	a.clocs = append(a.clocs, name)
	a.Unlock()
}
