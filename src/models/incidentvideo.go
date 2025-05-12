package models

import "time"

type IncidentVideo struct {
	Id            uint      `json:"id" gorm:"primaryKey"`
	IncidentId    uint      `json:"incident_id" gorm:"index:uidx_incident_video_1,unique;not null"`
	CameraId      uint      `json:"camera_id" gorm:"index:uidx_incident_video_1,unique;not null"`
	StartDatetime time.Time `json:"start_datetime" gorm:"index:uidx_incident_video_1,unique;not null"`
	EndDatetime   time.Time `json:"end_datetime" gorm:"not null"`
	Url           string    `json:"url" gorm:"not null"`
}
