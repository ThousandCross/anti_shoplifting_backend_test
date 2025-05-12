package main

import (
	"anti-shoplifting/src/database"
	"anti-shoplifting/src/models"
	"fmt"
	"time"
)

func main() {
	database.Connect()

	now := time.Now()
	fmt.Println(now.Format(time.RFC3339))

	nowUTC := now.UTC()
	fmt.Println(nowUTC.Format(time.RFC3339))
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)

	nowJST := nowUTC.In(jst)

	for i := 0; i < 10; i++ {
		start_datetime := nowJST.Add(time.Duration(i) * time.Minute)
		end_datetime := start_datetime.Add(time.Duration(2) * time.Minute)

		incident := models.Incident{
			UserId:         1,
			IncidentTypeId: 1,
			StartDatetime:  start_datetime,
			EndDatetime:    end_datetime,
		}

		database.DB.Create(&incident)
	}
}
