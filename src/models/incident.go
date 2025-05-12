package models

import "time"

type Incident struct {
	Id             uint            `json:"id" gorm:"primaryKey"`
	UserId         uint            `json:"user_id" gorm:"not null"`
	IncidentTypeId uint            `json:"incident_type_id" gorm:"not null"`
	StartDatetime  time.Time       `json:"start_datetime" gorm:"not null"`
	EndDatetime    time.Time       `json:"end_datetime" gorm:"not null"`
	IncidentVideo  []IncidentVideo `gorm:"foreignKey:IncidentId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	// incidentとBlacklistの関連を保存するテーブルの外部キーとなる
	// テーブル名をどうするべきか？
}
