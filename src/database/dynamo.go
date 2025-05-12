package database

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var DynamoDB *dynamodb.Client

func ConnectDynamoDB() {
	var err error
	var ctx = context.Background()

	cred := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider("", "", ""))
	//value, err := cred.Retrieve(context.TODO())
	if cred == nil {
		panic("failed to fetch credentials provider")
	}

	// DynamoDB クライアントの生成
	cfg, err := config.LoadDefaultConfig(ctx, config.WithCredentialsProvider(cred), config.WithRegion("ap-northeast-1"))
	if err != nil {
		fmt.Printf("load aws config: %s\n", err.Error())
		return
	}
	DynamoDB = dynamodb.NewFromConfig(cfg)
}
