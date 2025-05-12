package controllers

import (
	"anti-shoplifting/src/database"
	"anti-shoplifting/src/middlewares"
	"anti-shoplifting/src/models"
	"anti-shoplifting/src/notifier"
	"anti-shoplifting/src/utils"

	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"firebase.google.com/go/messaging"
	"github.com/gofiber/fiber/v2"
)

func Incidents(c *fiber.Ctx) error {
	user_id, _ := middlewares.GetUserId(c)

	type Incidents struct {
		Id             uint      `json:"incident_id"`
		IncidentTypeId uint      `json:"incident_type_id"`
		UserId         uint      `json:"user_id"`
		StartDatetime  time.Time `json:"start_datetime"`
		EndDatetime    time.Time `json:"end_datetime"`
		Name           string    `json:"incident_type_name"`
	}
	var incidents []Incidents

	database.DB.Model(&models.Incident{}).Select(
		"incidents.id",
		"incidents.incident_type_id",
		"incidents.user_id",
		"incidents.start_datetime",
		"incidents.end_datetime",
		"incident_types.name",
	).Joins(
		"inner join incident_types on incident_types.id = incidents.incident_type_id",
	).Where(
		"incidents.user_id = ?",
		user_id,
	).Order(
		"incidents.start_datetime",
	).Find(&incidents)

	type Videos struct {
		StartDatetime time.Time `json:"start_datetime"`
		EndDatetime   time.Time `json:"end_datetime"`
		Url           string    `json:"url"`
		SerialNo      string    `json:"serial_no"`
	}

	var result []map[string]interface{}
	var videos []Videos
	for _, incident := range incidents {
		database.DB.Model(&models.IncidentVideo{}).Select(
			"incident_videos.start_datetime",
			"incident_videos.end_datetime",
			"incident_videos.url",
			"cameras.serial_no",
		).Joins(
			"inner join cameras on cameras.id = incident_videos.camera_id",
		).Where(
			"incident_videos.incident_id = ? and cameras.user_id = ?",
			incident.Id,
			user_id,
		).Order(
			"incident_videos.start_datetime",
		).Find(&videos)

		var tempMap map[string]interface{}
		data, _ := json.Marshal(incident)
		json.Unmarshal(data, &tempMap)

		tempMap["videos"] = videos

		result = append(result, tempMap)
	}

	return c.JSON(result)
}

func RegistIncidents(c *fiber.Ctx) error {
	user_id, _ := middlewares.GetUserId(c)

	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	start_datetime := utils.StringToTime(data["start_datetime"])
	end_datetime := utils.StringToTime(data["end_datetime"])

	incident := models.Incident{
		UserId:         user_id,
		IncidentTypeId: 1,
		StartDatetime:  start_datetime,
		EndDatetime:    end_datetime,
	}

	database.DB.Create(&incident)

	// regist incidentvideo within loop
	// pending!!!

	//go database.ClearCache("products_frontend", "products_backend")

	return c.JSON(incident)

}

// func GetIncident(c *fiber.Ctx) error {
// 	user_id, _ := middlewares.GetUserId(c)

// 	id, _ := strconv.Atoi(c.Params("id"))

// 	var blacklist models.BlackList

// 	// validation
// 	result := database.DB.Model(&models.Camera{}).Where("blacklists.id = ? AND blacklists.user_id = ?", id, user_id).First(&blacklist)
// 	if result.Error != nil {
// 		c.Status(fiber.StatusBadRequest)
// 		return c.JSON(fiber.Map{
// 			"message": "Invalid Operations!",
// 		})
// 	}

// 	return c.JSON(&blacklist)
// }

func NotifyIncidents(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.JSON(fiber.Map{
			"result":  "ng",
			"message": "Invalid Operations!",
		})
	}

	// validation
	if _, ok := data["user_id"]; !ok {
		return c.JSON(fiber.Map{
			"result":  "ng",
			"message": "Invalid Operations!",
		})
	}

	registrationTokens := []string{}

	database.DB.Model(&models.FcmToken{}).Select(
		"fcm_tokens.token",
	).Where(
		"fcm_tokens.user_id = ? and fcm_tokens.is_enabled = ?",
		data["user_id"],
		true,
	).Find(&registrationTokens)

	// Create a list containing up to 500 registration tokens.
	// This registration tokens come from the client FCM SDKs.
	// registrationTokens := []string{
	// 	"YOUR_REGISTRATION_TOKEN_1",
	// 	// ...
	// 	"YOUR_REGISTRATION_TOKEN_n",
	// }
	fmt.Printf("%v\n", registrationTokens)

	message := &messaging.MulticastMessage{
		Tokens: registrationTokens,
		Data: map[string]string{
			"score": "850",
			"time":  "2:45",
		},
		Notification: &messaging.Notification{
			Title:    "Title of Your Notification A",
			Body:     "Body of Your Notification",
			ImageURL: "",
		},
		Android: &messaging.AndroidConfig{
			Priority: "high",
		},
	}

	br, err := notifier.FcmClient.SendMulticast(context.Background(), message)
	if err != nil {
		log.Fatalln(err)
	}

	// See the BatchResponse reference documentation
	// for the contents of response.
	fmt.Printf("%d messages were sent successfully\n", br.SuccessCount)

	return c.JSON(fiber.Map{
		"result": "ok",
	})
}
