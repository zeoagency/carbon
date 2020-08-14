package models

import (
	"errors"
	neturl "net/url"
	"strings"
)

// URLSet is kind a Set Data Structure implementation.
// It has an Add method that only works if the URL doesn't exist already.
// Also, It is able to split the url to FullURL, BaseURL and Keywords.
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
// like that: s.URLs["boratanrikulu.dev postgresql nedir nasil calisir"]
type URLSet struct {
	URLs map[string]URL
}

// URL is used to keep the URLs as parsed.
type URL struct {
	FullURL  string
	BaseURL  string
	Keywords string
}

// NewURLSet inits the URLSet to use.
func NewURLSet() *URLSet {
	var s URLSet
	s.URLs = make(map[string]URL)
	return &s
}

// Add method adds new URLs if it is doesn't already exist.
// Also, It splits the url to FullURL, BaseURL and Keywords.
// It except only valid URLs.
func (s *URLSet) Add(urls ...string) error {
	for _, url := range urls {
		u, err := convertToURL(url)
		if err != nil {
			return err
		}

		if _, ok := s.URLs[u.String()]; ok {
			// Continue if it is already exists.
			continue
		}

		s.URLs[u.String()] = u
	}

	return nil
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
func (u *URL) String() string {
	result := u.BaseURL + " " + u.Keywords
	return result
}

// convertToURL converts the given string to a parsed URL.
//
// For example;
// FROM:     https://boratanrikulu.dev/postgresql-nedir-nasil-calisir.html
// TO:
// FullURL:  https://boratanrikulu.dev/postgresql-nedir-nasil-calisir.html
// BaseURL:  boratanrikulu.dev
// Keywords: postgresql nedir nasil calisir
func convertToURL(fullURL string) (URL, error) {
	result := URL{}

	u, err := neturl.Parse(fullURL)
	parts := strings.Split(u.Hostname(), ".")
	if err != nil || u.Scheme == "" || len(parts) < 2 {
		return result, errors.New("The input includes non-url(s). Please check your input.")
	}

	result.FullURL = u.String()
	result.BaseURL = parts[len(parts)-2] + "." + parts[len(parts)-1]
	result.Keywords = replaceAnyWithSpace(
		u.EscapedPath(),
		"/", "\\", "-", "_", ".", "html", "js", "php", "aspx",
	)

	return result, nil
}

// replaceAnyWithSpace replace all given keys with space.
func replaceAnyWithSpace(s string, keys ...string) string {
	for _, k := range keys {
		s = strings.ReplaceAll(s, k, " ")
	}
	return strings.TrimSpace(s)
}
