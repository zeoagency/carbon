package helpers

import (
	"fmt"
	"log"
	"testing"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalln(err)
	}
}

func TestRandomAPICred(t *testing.T) {
	address, key, err := RandomAPICred()
	if err != nil {
		t.Fatal(err)
	}

	if address == "" || key == "" {
		t.Fatal("Address or key is empty.")
	}

	fmt.Println(address, key)
}
