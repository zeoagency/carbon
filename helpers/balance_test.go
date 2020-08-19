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
	api, index := RandomAPICred()

	if (index != -1) && (api.Keys[index].Address == "" || api.Keys[index].Key == "") {
		t.Fatal("Address or key is empty.")
	}

	fmt.Println(api.Keys[index].Address, api.Keys[index].Key)
}
