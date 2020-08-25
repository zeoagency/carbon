package models

import "strings"

// KeywordSet is kind a Set Data Structure implementation.
// It has an Add method that only works if the Keyword doesn't exist already.
//
// When you use it, firstly create with NewKeywordSet method.
// Example:
// k := NewKeywordSet()
// k.Add(
// 	 "boratanrikulu blog yazıları",
// )
//
// You can access by using keyword string output,
// like that: k.Keywords["boratanrikulu blog yazıları"]
type KeywordSet struct {
	Keywords  map[string]bool           // the key is the keyword.
	Successes map[string]keywordSuccess // the key is the keyword.
	Fails     map[string]keywordFail    // the key is the keyword.
}

// NewKeywordSet inits the KeywordSet to use.
func NewKeywordSet() *KeywordSet {
	var k KeywordSet
	k.Keywords = make(map[string]bool)
	k.Successes = make(map[string]keywordSuccess)
	k.Fails = make(map[string]keywordFail)
	return &k
}

// keywordSuccess is used to keep the result.
type keywordSuccess struct {
	Results []KeywordSuccessResult
}

// KeywordSuccessResults keeps the result for success elements.
// It is exported to use in services package.
type KeywordSuccessResult struct {
	Title string
	Desc  string
	URL   string
}

// fail is used to keep keywords that could not be processed with a reason.
type keywordFail struct {
	Reason string
}

// Add method adds new keywords if it is doesn't already exist.
func (k *KeywordSet) Add(keywords ...string) error {
	for _, keyword := range keywords {
		keyword = strings.TrimSpace(keyword)
		if keyword == "" {
			continue // If the keyword is empty, do nothing.
		}
		if _, ok := k.Keywords[keyword]; ok {
			// Continue if it is already exists.
			continue
		}

		k.Keywords[keyword] = true
	}

	return nil
}

// AddSuccess adds the result to the success list, if it doesn't exist already.
func (k *KeywordSet) AddSuccess(keyword string, results []KeywordSuccessResult) {
	if _, ok := k.Successes[keyword]; ok {
		return // Return if it is already exists.
	}

	k.Successes[keyword] = keywordSuccess{
		Results: results,
	}
}

// AddFail adds the keyword to the fail list with a reason, if it doesn't exist already.
func (k *KeywordSet) AddFail(keyword string, reason string) {
	if _, ok := k.Fails[keyword]; ok {
		return // Return if it is already exists.
	}

	k.Fails[keyword] = keywordFail{
		Reason: reason,
	}
}

// ToStringSlice create a string slice that includes all Keywords.
func (kw *KeywordSet) ToStringSlice() []string {
	r := []string{}

	for k, _ := range kw.Keywords {
		r = append(r, k)
	}

	return r
}

// Returns the original value for the given string.
func (kw *KeywordSet) Original(k string) string {
	if _, ok := kw.Keywords[k]; ok {
		return k
	}
	return ""
}
