package helpers

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type cred struct {
	Keys []struct {
		Address string `json:"address"`
		Key     string `json:"key"`
	} `json:"keys"`
}

// RandomAPICred returns randomly selected api-key values.
func RandomAPICred() (string, string, error) {
	credJSON := os.Getenv("SERP_API_CREDENTIALS_JSON")

	apiCred := cred{}
	err := json.Unmarshal([]byte(credJSON), &apiCred)
	if err != nil {
		fmt.Println(err)
		return "", "", err
	}

	length := len(apiCred.Keys)
	if length == 0 {
		return "", "", errors.New("No api-key.")
	}

	rand.Seed(time.Now().UnixNano())
	selected := rand.Intn(length)
	address, key := apiCred.Keys[selected].Address, apiCred.Keys[selected].Key

	// If the keys are empty, return error.
	if address == "" && key == "" {
		return "", "", errors.New("No api-key.")
	}

	return address, key, nil
}
