package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"

	"gitlab.com/seo.do/zeo-carbon/controllers"
)

func init() {
	_ = godotenv.Load(".env")
}

func main() {
	lambda.Start(controllers.Result)
}
