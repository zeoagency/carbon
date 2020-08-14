package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	"gitlab.com/seo.do/zeo-carbon/controllers"
)

func main() {
	lambda.Start(controllers.Result)
}
