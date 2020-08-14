package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	"gitlab.com/zeo.org/carbon/controllers"
)

func main() {
	lambda.Start(controllers.Result)
}
