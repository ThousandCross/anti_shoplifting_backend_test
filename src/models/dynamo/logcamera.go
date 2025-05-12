package dynamomodels

import (
	"time"
)

type LogCamera struct {
	CameraId        int       `json:"camera_id" dynamodbav:"CameraId"`
	LatestTimestamp time.Time `json:"latest_timestamp" dynamodbav:"LatestTimestamp"`
	UserId          int       `json:"user_id" dynamodbav:"UserId"`
}
