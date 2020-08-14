package main

import (
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"

	"gitlab.com/seo.do/zeo-carbon/controllers"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	lambda.Start(controllers.Result)
}
