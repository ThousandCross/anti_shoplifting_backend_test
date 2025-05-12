package main

import (
	"anti-shoplifting/src/database"
	"anti-shoplifting/src/models"
	"fmt"
	"math/rand"
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
	fmt.Println(nowJST.Format(time.RFC3339))

	urls := []string{
		"http://commondatastorage.googleapis.com/gtv-videos-bucket/sample/BigBuckBunny.mp4",
		"http://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ElephantsDream.mp4",
		"http://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerBlazes.mp4",
		"http://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerEscapes.mp4",
		"http://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerFun.mp4",
		"http://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerJoyrides.mp4",
		"http://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerMeltdowns.mp4",
		"http://commondatastorage.googleapis.com/gtv-videos-bucket/sample/Sintel.mp4",
		"http://commondatastorage.googleapis.com/gtv-videos-bucket/sample/SubaruOutbackOnStreetAndDirt.mp4",
		"http://commondatastorage.googleapis.com/gtv-videos-bucket/sample/TearsOfSteel.mp4",
	}

	for i := 1; i < 11; i++ {
		for j := 1; j < 3; j++ {
			// choose video randomly
			num := rand.Intn(len(urls))

			start_datetime := nowJST.Add(time.Duration(2*(i-1)+j-1) * time.Minute)
			end_datetime := start_datetime.Add(time.Duration(1) * time.Minute)

			incident_video := models.IncidentVideo{
				IncidentId:    uint(i),
				CameraId:      uint(2*(i-1) + j),
				StartDatetime: start_datetime,
				EndDatetime:   end_datetime,
				Url:           urls[num],
			}

			database.DB.Create(&incident_video)
		}
	}
}
