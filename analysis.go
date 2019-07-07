package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/hhatto/gocloc"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

// Analysis
type Analysis struct {
	repos  []string
	clocs  []string
	github *GitHub
	result *gocloc.Result
	wg     sync.WaitGroup
	sync.RWMutex
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
			if err.Error() != "remote repository is empty" {
				fmt.Fprintf(os.Stderr, fmt.Sprintf("clone error: %s\n%s\n", URL, err))
			}
			return
		} else {
			fmt.Fprintf(os.Stderr, fmt.Sprintf("cloned %s\n", URL))
		}
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
