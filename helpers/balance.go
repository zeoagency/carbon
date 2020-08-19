package helpers

import (
	"encoding/json"
	"math/rand"
	"os"
	"time"
)

type API struct {
	Keys []struct {
		Address string `json:"address"`
		Key     string `json:"key"`
	} `json:"keys"`
}

// RandomAPICred returns randomly selected api-key index and the key list.
func RandomAPICred() (API, int) {
	credJSON := os.Getenv("SERP_API_CREDENTIALS_JSON")

	api := API{}
	err := json.Unmarshal([]byte(credJSON), &api)
	if err != nil {
		return api, -1
	}

	length := len(api.Keys)
	if length == 0 {
		return api, -1
	}

	rand.Seed(time.Now().UnixNano())
	selected := rand.Intn(length)

	// If the keys are empty, return error.
	if api.Keys[selected].Address == "" || api.Keys[selected].Key == "" {
		return api, -1
	}

	return api, selected
}
