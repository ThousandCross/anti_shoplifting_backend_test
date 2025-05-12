package models

type IncidentType struct {
	Id       uint       `json:"id" gorm:"primaryKey"`
	Name     string     `json:"name" gorm:"not null"`
	Incident []Incident `gorm:"foreignKey:IncidentTypeId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
