package main

import (
	"anti-shoplifting/src/database"
	dynamomodels "anti-shoplifting/src/models/dynamo"
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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

	company_cd := "crpMreQRxT"
	store_cd := "strTVrFFY9"
	pre_sensor_cd := company_cd + store_cd

	// newly add test records
	for i := 1; i < 11; i++ {
		putInput := dynamomodels.LogSensor{
			SensorCd:        pre_sensor_cd + fmt.Sprintf("%010d", i),
			LatestTimestamp: nowJST,
			UserId:          1,
			WeightValue:     "0.0",
		}
		av, err := attributevalue.MarshalMap(putInput)
		if err != nil {
			fmt.Printf("dynamodb marshal: %s\n", err.Error())
			return
		}
		_, err = database.DynamoDB.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String("LogSensor"),
			Item:      av,
		})
		if err != nil {
			fmt.Printf("put item: %s\n", err.Error())
			return
		}
	}

	// seach all records by hash key
	fmt.Println("seach all records by hash key")
	getInput := &dynamodb.GetItemInput{
		TableName: aws.String("LogSensor"),
		Key: map[string]types.AttributeValue{
			"SensorCd": &types.AttributeValueMemberS{
				Value: "crpMreQRxTstrTVrFFY90000000009",
			},
		},
	}
	output, err := database.DynamoDB.GetItem(ctx, getInput)
	if err != nil {
		fmt.Printf("get item: %s\n", err.Error())
		return
	}
	gotLogSensor := dynamomodels.LogSensor{}
	err = attributevalue.UnmarshalMap(output.Item, &gotLogSensor)
	if err != nil {
		fmt.Printf("dynamodb unmarshal: %s\n", err.Error())
		return
	}
	fmt.Println(gotLogSensor)

	// seach all records by global secondary index
	fmt.Println("seach all records by global secondary index")
	out, err := database.DynamoDB.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String("LogSensor"),
		IndexName:              aws.String("UserId-index"),
		KeyConditionExpression: aws.String("UserId = :userId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":userId": &types.AttributeValueMemberN{Value: "1"},
		},
		// ExpressionAttributeValues: map[string]types.AttributeValue{
		//     ":hashKey":  &types.AttributeValueMemberS{Value: "123"},
		//     ":rangeKey": &types.AttributeValueMemberN{Value: "20150101"},
		// },
		// ExpressionAttributeNames: map[string]string{
		//     "#date": "date",
		// },
		//ScanIndexForward: aws.Bool(true), // true or false to sort by "date" Sort/Range key ascending or descending
	})
	if err != nil {
		panic(err)
	}

	var logs []dynamomodels.LogSensor
	attributevalue.UnmarshalListOfMaps(out.Items, &logs)

	// var logs []models.LogSensor = []models.LogSensor{}
	// for _, log := range out.Items {
	// 	//var log models.LogSensor
	// 	gotLogSensor := models.LogSensor{}
	// 	_ = attributevalue.UnmarshalListOfMaps(log, &gotLogSensor)
	// 	logs = append(logs, gotLogSensor)
	// }
	sort.SliceStable(logs, func(i, j int) bool { return logs[i].SensorCd < logs[j].SensorCd })
	fmt.Printf("SensorCdでSort(Stable):%+v\n", logs)

	//fmt.Println(logs)
}
