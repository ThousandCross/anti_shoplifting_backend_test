package dynamomodels

import (
	"time"
)

type LogSensor struct {
	SensorCd        string    `json:"sensor_cd" dynamo:"SensorCd,hash"`
	LatestTimestamp time.Time `json:"latest_timestamp" dynamodbav:"LatestTimestamp"`
	UserId          int       `json:"user_id" dynamodbav:"UserId"`
	WeightValue     string    `json:"weight_value" dynamodbav:"WeightValue"`
}
