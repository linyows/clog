package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// JSONData
type JSONData struct {
	Langs Langs `json:"languages"`
	Total Lang  `json:"total"`
}

// Result
type Result struct {
	name string
	data JSONData
}

// Report
type Report struct {
	analysis *Analysis
	results  []Result
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
		var data JSONData
		err := json.Unmarshal(bytes, &data)
		if err != nil {
			fmt.Fprintf(os.Stderr, fmt.Sprintf("Error: %s\n", err))
			continue
		}
		r.results = append(r.results, Result{
			name: strings.Replace(name, ".json", "", 1),
			data: data,
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

	var orgLangs Langs
	orgTotal := &Lang{Name: "TOTAL", FilesCount: 0, Comments: 0, Code: 0, Blanks: 0}

	for _, rr := range r.results {
		tt := rr.data.Total
		ll := rr.data.Langs

		for _, l := range ll {
			nofound := true
			for _, ol := range orgLangs {
				if ol.Name == l.Name {
					ol.FilesCount += l.FilesCount
					ol.Blanks += l.Blanks
					ol.Comments += l.Comments
					ol.Code += l.Code
					nofound = false
					break
				}
			}
			if nofound {
				orgLangs = append(orgLangs, l)
			}
		}

		orgTotal.FilesCount += tt.FilesCount
		orgTotal.Blanks += tt.Blanks
		orgTotal.Comments += tt.Comments
		orgTotal.Code += tt.Code
	}

	var sortedOrgLangs Langs
	for _, l := range orgLangs {
		sortedOrgLangs = append(sortedOrgLangs, l)
	}
	sort.Sort(sortedOrgLangs)

	fmt.Printf("%.[2]*[1]s\n", separator, rowLen)
	fmt.Printf("%-[2]*[1]s %[3]s\n", header, headerLen, commonHeader)
	fmt.Printf("%.[2]*[1]s\n", separator, rowLen)
	for _, v := range sortedOrgLangs {
		fmt.Printf("%-27v %6v %14v %14v %14v\n", v.Name, v.FilesCount, v.Blanks, v.Comments, v.Code)
	}
	fmt.Printf("%.[2]*[1]s\n", separator, rowLen)
	fmt.Printf("%-27v %6v %14v %14v %14v\n", "TOTAL", orgTotal.FilesCount, orgTotal.Blanks, orgTotal.Comments, orgTotal.Code)
	fmt.Printf("%.[2]*[1]s\n", separator, rowLen)

	for _, rr := range r.results {
		if len(rr.data.Langs) == 0 {
			continue
		}
		ll := rr.data.Langs
		fmt.Printf("%s/%s\n--\n", r.analysis.github.org, rr.name)
		for _, l := range ll {
			fmt.Printf("%-27v %6v %14v %14v %14v\n", l.Name, l.FilesCount, l.Blanks, l.Comments, l.Code)
		}
		fmt.Printf("%.[2]*[1]s\n", separator, rowLen)
	}
}
