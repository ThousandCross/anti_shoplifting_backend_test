package main

import (
	"anti-shoplifting/src/database"
	"anti-shoplifting/src/models"
	"anti-shoplifting/src/utils"
	"fmt"
	"math/rand"
	"strconv"

	"github.com/bxcodec/faker/v3"
)

func main() {
	database.Connect()

	for i := 1; i < 21; i++ {
		camera := models.Camera{
			SerialNo:  "CAM" + fmt.Sprintf("%07d", i),
			UserId:    1,
			IpAddress: faker.IPv4(),
			Latitude:  fmt.Sprintf("%g", utils.MakeRandomFloats(-90.0, 90.0)),
			Longitude: fmt.Sprintf("%g", utils.MakeRandomFloats(-180.0, 180.0)),
			Elevation: strconv.Itoa(300 + rand.Intn(100)),
		}

		database.DB.Create(&camera)
	}
}
