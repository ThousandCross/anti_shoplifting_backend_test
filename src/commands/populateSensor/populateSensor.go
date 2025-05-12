package main

import (
	"anti-shoplifting/src/database"
	"anti-shoplifting/src/models"
	"anti-shoplifting/src/utils"
	"fmt"
	"math/rand"

	"github.com/bxcodec/faker/v3"
)

func main() {
	database.Connect()

	company_cd := "crpMreQRxT"
	store_cd := "strTVrFFY9"
	pre_sensor_cd := company_cd + store_cd

	for i := 1; i < 11; i++ {
		camera := models.Sensor{
			SensorCd:    pre_sensor_cd + fmt.Sprintf("%010d", i),
			UserId:      1,
			IpAddress:   faker.IPv4(),
			Latitude:    utils.MakeRandomFloats(-90.0, 90.0),
			Longitude:   utils.MakeRandomFloats(-180.0, 180.0),
			Elevation:   rand.Float64() * (400 - 300),
			SalesItemId: uint(i),
		}

		database.DB.Create(&camera)
	}
}
