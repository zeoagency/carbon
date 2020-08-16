package models

import (
	"gitlab.com/seo.do/zeo-carbon/helpers"
)

// URLSet is kind a Set Data Structure implementation.
// It has an Add method that only works if the URL doesn't exist already.
// Also, It is able to split the url to FullURL, BaseURL and Keywords.
// If there is an error, add it to the fail list. (only adds when the fail info doesn't exist already.)
//
// When you use it, firstly create with NewURLSet method.
// Example:
// s := NewURLSet()
// s.Add(
// 	 "https://boratanrikulu.dev/postgresql-nedir-nasil-calisir/",
//   "https://boratanrikulu.dev/smtp-nasil-calisir-ve-postfix-kurulumu/",
// )
//
// You can access by using string output,
// like that: s.URLs["boratanrikulu.dev postgresql nedir nasil calisir"].BaseURL
type URLSet struct {
	URLs      map[string]url     // the key is url.String() (BaseURL + Keywords).
	Successes map[string]success // the key is the Original URL.
	Fails     map[string]fail    // the key is the Original URL.
}

// NewURLSet inits the URLSet to use.
func NewURLSet() *URLSet {
	var s URLSet
	s.URLs = make(map[string]url)
	s.Successes = make(map[string]success)
	s.Fails = make(map[string]fail)
	return &s
}

// url is used to keep the URL as parsed.
type url struct {
	FullURL  string
	BaseURL  string
	Keywords string
}

// urlSucces is used to keep the result.
type success struct {
	RelatedURLs []string
}

// urlFail is used to keep urls that could not be processed with a reason.
type fail struct {
	Reason string
}

// Add method adds new URLs if it is doesn't exist already.
// Also, It splits the url to FullURL, BaseURL and Keywords.
// It except only valid URLs.
// If there is an issue with the given URL, It adds it to the fail list with a reason.
func (s *URLSet) Add(urls ...string) {
	for _, url := range urls {
		u, err := convertToURL(url)
		if err != nil {
			s.AddFail(url, "That's not a URL.")
			continue
		}

		if _, ok := s.URLs[u.String()]; ok {
			// Continue if it is already exists.
			continue
		}

		s.URLs[u.String()] = u
	}
}

// AddSuccess adds the result to the success list, if it doesn't exist already.
func (s *URLSet) AddSuccess(originalURL string, result []string) {
	if _, ok := s.Successes[originalURL]; ok {
		return // Return if it is already exists.
	}

	s.Successes[originalURL] = success{
		RelatedURLs: result,
	}
}

// AddFail adds the that's failed to the fail list with a reason, if it doesn't exist already.
func (s *URLSet) AddFail(originalURL string, reason string) {
	if _, ok := s.Fails[originalURL]; ok {
		return // Return if it is already exists.
	}

	s.Fails[originalURL] = fail{
		Reason: reason,
	}
}

// convertToURL converts the given string to a parsed URL.
//
// For example;
// FROM:     https://boratanrikulu.dev/postgresql-nedir-nasil-calisir.html
// TO:
// FullURL:  https://boratanrikulu.dev/postgresql-nedir-nasil-calisir.html
// BaseURL:  boratanrikulu.dev
// Keywords: postgresql nedir nasil calisir
func convertToURL(fullURL string) (url, error) {
	result := url{}

	baseURL, keywords, err := helpers.ExtractURL(fullURL)
	if err != nil {
		return result, err
	}

	result.FullURL = fullURL
	result.BaseURL = baseURL
	result.Keywords = keywords

	return result, nil
}

// ToStringSlice create a string slice that includes all URLs.
func (u *URLSet) ToStringSlice() []string {
	r := []string{}

	for k, _ := range u.URLs {
		r = append(r, k)
	}

	return r
}

// String method for URL model.
func (u *url) String() string {
	result := u.BaseURL + " " + u.Keywords
	return result
}
