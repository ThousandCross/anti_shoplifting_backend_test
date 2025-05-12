package main

import (
	"anti-shoplifting/src/database"
	dynamomodels "anti-shoplifting/src/models/dynamo"
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func main() {

	var ctx = context.Background()

	database.ConnectDynamoDB()

	now := time.Now()
	fmt.Println(now.Format(time.RFC3339))

	nowUTC := now.UTC()
	fmt.Println(nowUTC.Format(time.RFC3339))
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)

	nowJST := nowUTC.In(jst)
	fmt.Println(nowJST.Format(time.RFC3339))

	for i := 1; i < 21; i++ {
		putInput := dynamomodels.LogCamera{
			CameraId:        i,
			LatestTimestamp: nowJST,
			UserId:          1,
		}
		av, err := attributevalue.MarshalMap(putInput)
		if err != nil {
			fmt.Printf("dynamodb marshal: %s\n", err.Error())
			return
		}
		_, err = database.DynamoDB.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String("LogCamera"),
			Item:      av,
		})
		if err != nil {
			fmt.Printf("put item: %s\n", err.Error())
			return
		}
	}
}
