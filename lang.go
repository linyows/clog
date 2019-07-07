package main

// Lang
type Lang struct {
	Name       string `json:"name,omitempty"`
	FilesCount int32  `json:"files"`
	Code       int32  `json:"code"`
	Comments   int32  `json:"comment"`
	Blanks     int32  `json:"blank"`
}

// Langs
type Langs []*Lang

func (ll Langs) Len() int {
	return len(ll)
}

func (ll Langs) Swap(i, j int) {
	ll[i], ll[j] = ll[j], ll[i]
}

func (ll Langs) Less(i, j int) bool {
	if ll[i].Code == ll[j].Code {
		return ll[i].Name < ll[j].Name
	}
	return ll[i].Code > ll[j].Code
}
