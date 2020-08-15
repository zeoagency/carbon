package models

import (
	"errors"
	neturl "net/url"
	"strings"
)

// countryTopLevelDomains includes all country-code-top-level-domains.
// source: https://en.wikipedia.org/wiki/List_of_Internet_top-level_domains#Country_code_top-level_domains
var countryTopLevelDomains = []string{
	"ac", "ad", "ae", "af", "ag", "ai", "al", "am", "ao", "aq", "ar", "as", "at", "au", "aw", "ax", "az", "ba", "bb", "bd", "be", "bf", "bg", "bh", "bi", "bj", "bm", "bn", "bo", "bq", "br", "bs", "bt", "bw", "by", "bz", "ca", "cc", "cd", "cf", "cg", "ch", "ci", "ck", "cl", "cm", "cn", "co", "cr", "cu", "cv", "cw", "cx", "cy", "cz", "de", "dj", "dk", "dm", "do", "dz", "ec", "ee", "eg", "eh", "er", "es", "et", "eu", "fi", "fj", "fk", "fm", "fo", "fr", "ga", "gd", "ge", "gf", "gg", "gh", "gi", "gl", "gm", "gn", "gp", "gq", "gr", "gs", "gt", "gu", "gw", "gy", "hk", "hm", "hn", "hr", "ht", "hu", "id", "ie", "il", "im", "in", "io", "iq", "ir", "is", "it", "je", "jm", "jo", "jp", "ke", "kg", "kh", "ki", "km", "kn", "kp", "kr", "kw", "ky", "kz", "la", "lb", "lc", "li", "lk", "lr", "ls", "lt", "lu", "lv", "ly", "ma", "mc", "md", "me", "mg", "mh", "mk", "ml", "mm", "mn", "mo", "mp", "mq", "mr", "ms", "mt", "mu", "mv", "mw", "mx", "my", "mz", "na", "nc", "ne", "nf", "ng", "ni", "nl", "no", "np", "nr", "nu", "nz", "om", "pa", "pe", "pf", "pg", "ph", "pk", "pl", "pm", "pn", "pr", "ps", "pt", "pw", "py", "qa", "re", "ro", "rs", "ru", "rw", "sa", "sb", "sc", "sd", "se", "sg", "sh", "si", "sk", "sl", "sm", "sn", "so", "sr", "ss", "st", "su", "sv", "sx", "sy", "sz", "tc", "td", "tf", "tg", "th", "tj", "tk", "tl", "tm", "tn", "to", "tr", "tt", "tv", "tw", "tz", "ua", "ug", "uk", "us", "uy", "uz", "va", "vc", "ve", "vg", "vi", "vn", "vu", "wf", "ws", "ye", "yt", "za", "zm", "zw",
}

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
			return errors.New("The input includes non-url(s). Please check your input.")
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

	baseURL, keywords, err := ExtractURL(fullURL)
	if err != nil {
		return result, err
	}

	result.FullURL = fullURL
	result.BaseURL = baseURL
	result.Keywords = keywords

	return result, nil
}

// ExtractURL works like this:
//
// Given URL: "https://text.blog.boratanrikulu.dev.tr/archlinux-install.html"
// Result:
//   BaseURL: "boratanrikulu.dev"
//   Keywords: "archlinux install"
func ExtractURL(url string) (string, string, error) {
	u, err := neturl.Parse(url)
	if err != nil || u.Scheme == "" {
		return "", "", err
	}

	parts := strings.Split(u.Hostname(), ".")

	count := 2
	if len(parts) > 2 && stringContains(countryTopLevelDomains, parts[len(parts)-1]) {
		// Set count to 3 if it is an url that contins country domain.
		count = 3
	}

	if len(parts) < count {
		return "", "", errors.New("That's not a valid hostname-url.")
	}

	keywords := replaceAnyWithSpace(
		u.EscapedPath(),
		"/", "\\", "-", "_", ".", "html", "js", "php", "aspx",
	)

	// TODO: there is an issue.
	// What will be happend if the url includes country flags?
	// like: boratanrikulu.dev.tr
	// Need to fix it1
	return parts[len(parts)-count] + "." + parts[len(parts)-(count-1)], keywords, nil
}

// replaceAnyWithSpace replace all given keys with space.
func replaceAnyWithSpace(s string, keys ...string) string {
	for _, k := range keys {
		s = strings.ReplaceAll(s, k, " ")
	}
	return strings.TrimSpace(s)
}

// contains tells whether a contains x.
func stringContains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
