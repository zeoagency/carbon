package helpers

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

// ExtractURL works like this:
//
// Given URL: "https://text.blog.boratanrikulu.dev.tr/archlinux-install.html"
// Result:
//   BaseURL: "boratanrikulu.dev"
//   Keywords: "archlinux install"
func ExtractURL(url string) (string, string, error) {
	u, err := neturl.Parse(url)
	if err != nil || u.Scheme == "" {
		return "", "", errors.New("That's not a valid URL.")
	}

	parts := strings.Split(u.Hostname(), ".")

	count := 2
	_, contains := StringSliceContains(countryTopLevelDomains, parts[len(parts)-1])
	if len(parts) > 2 && contains {
		// Set count to 3 if it is an url that contains country domain.
		count = 3
	}

	if len(parts) < count {
		return "", "", errors.New("That's not a valid URL.")
	}

	keywords := ReplaceAnyWithSpace(
		u.EscapedPath(),
		"/", "\\", "-", "_", ".", "html", "js", "php", "aspx",
	)

	return parts[len(parts)-count] + "." + parts[len(parts)-(count-1)], keywords, nil
}
