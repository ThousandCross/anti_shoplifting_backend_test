package models

type BlackList struct {
	Id     uint   `json:"id" gorm:"primaryKey"`
	Name   string `json:"name" gorm:"not null"`
	UserId uint   `json:"user_id" gorm:"not null"`
}
