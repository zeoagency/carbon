package models

// KeywordSet is kind a Set Data Structure implementation.
// It has an Add method that only works if the Keyword is doesn't exist already.
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
	Keywords map[string]bool
}

// NewKeywordSet inits the KeywordSet to use.
func NewKeywordSet() *KeywordSet {
	var k KeywordSet
	k.Keywords = make(map[string]bool)
	return &k
}

// Add method adds new keywords if it is doesn't already exist.
func (k *KeywordSet) Add(keywords ...string) error {
	for _, keyword := range keywords {
		if _, ok := k.Keywords[keyword]; ok {
			// Continue if it is already exists.
			continue
		}

		k.Keywords[keyword] = true
	}

	return nil
}

// ToStringSlice create a string slice that includes all Keywords.
func (kw *KeywordSet) ToStringSlice() []string {
	r := []string{}

	for k, _ := range kw.Keywords {
		r = append(r, k)
	}

	return r
}
