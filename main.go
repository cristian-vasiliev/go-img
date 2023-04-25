package main

import (
	"fmt"
	"go-img/config"
	"go-img/controllers"
	"go-img/services"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.NewConfig()
	imageEncodingService := services.NewImageEncodingService()
	imageService := services.NewImageService(*imageEncodingService)
	imageController := controllers.NewImageController(*imageService, *cfg)

	engine := gin.Default()

	engine.GET("/image/:pathToImage", imageController.HandleImageRequest)

	// Check if the application is running on AWS Lambda
	if isRunningOnLambda() {
		// Wrap the Gin engine with the AWS Lambda proxy adapter
		lambdaAdapter := ginadapter.New(engine)

		// Start the Lambda function
		lambda.Start(lambdaAdapter.Proxy)
	} else {
		// Run the application locally
		engine.Run(fmt.Sprintf(":%d", cfg.Port))
	}
}

func isRunningOnLambda() bool {
	return os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != ""
}
