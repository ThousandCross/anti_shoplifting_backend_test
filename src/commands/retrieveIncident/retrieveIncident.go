package main

import (
	"anti-shoplifting/src/database"
	"anti-shoplifting/src/models"
	"encoding/json"
	"fmt"
	"time"
)

func main() {
	database.Connect()

	//var incidents []models.Incident
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
		1,
	).Order(
		"incidents.start_datetime",
	).Find(&incidents)

	type Videos struct {
		StartDatetime time.Time `json:"start_datetime"`
		EndDatetime   time.Time `json:"end_datetime"`
		Url           string    `json:"url"`
		SerialNo      string    `json:"serial_no"`
	}

	//var result = map[string]map[string]string{}
	var result []map[string]interface{}
	//var result = map[int]map[string]string{}
	//var result []map[string]string
	//var result []string
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
			1,
		).Order(
			"incident_videos.start_datetime",
		).Find(&videos)

		var tempMap map[string]interface{}
		data, _ := json.Marshal(incident)
		json.Unmarshal(data, &tempMap)

		tempMap["videos"] = videos

		result = append(result, tempMap)
		//result
		//jsons, _ := json.Marshal(tempMap)

		//fmt.Println(string(jsons))
	}

	jsons, _ := json.Marshal(result)

	fmt.Println(string(jsons))

}
