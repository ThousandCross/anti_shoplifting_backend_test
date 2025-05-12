package models

type Prefecture struct {
	Id      uint      `json:"id" gorm:"primaryKey"`
	Name    string    `json:"name" gorm:"not null"`
	Store   []Store   `gorm:"foreignKey:PrefectureId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Company []Company `gorm:"foreignKey:PrefectureId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
