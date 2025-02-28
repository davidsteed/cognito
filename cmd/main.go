package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"
	"github.com/davidsteed/cognito/lib/httpserver"
	"github.com/davidsteed/cognito/lib/logs"

	"go.uber.org/zap"
)

var (
	Version    = "dev"
	echoLambda *echoadapter.EchoLambda
)

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	defer logs.Log.Sync() //nolint: errcheck // No need

	return echoLambda.ProxyWithContext(ctx, req)
}

// main will only be called on a cold start when in lambda
func main() {
	defer func() {
		if r := recover(); r != nil {
			logs.Log.Error("panic detected", zap.Any("recover", r))
		}
	}()

	logs.NewZapWithLevel(Version)
	defer logs.Log.Sync() //nolint: errcheck // No need

	if err := Load(); err != nil {
		logs.Log.Error("failed to load config", zap.Error(err))
		return
	}

	if err := logs.SetLevel(Settings.LogLevel); err != nil {
		logs.Log.Error("failed to set log level", zap.Error(err))
		return
	}

	service,err := boostrap()
	if err != nil {
		logs.Log.Error("bootstrap failure", zap.Error(err))
		return
	}

	// this is mapped via api gateway custom domain
	// Log.Fatal("update basepath for api gateway custom domain and remove this line")
	echoLambda.StripBasePath("demo-api")

	// start service for local development
	if Version == "dev" {
		if err := service.Start(); err != nil {
			logs.Log.Warn("service failure", zap.Error(err))
		}
		return
	}

	lambda.Start(Handler)
}

// bootstrap initialise the service dependencies
func boostrap() (*Server, error) {
	// SessionConfig, err := config.LoadDefaultConfig(context.Background())
	// if err != nil {
	// 	return nil,fmt.Errorf("failed to create new AWS session: %w", err)
	// }

	service := NewServer(httpserver.New(logs.Log))
	echoLambda = echoadapter.New(service.Server)
	return service, nil
}
