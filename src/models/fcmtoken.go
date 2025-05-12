package models

type FcmToken struct {
	Id        uint   `json:"id" gorm:"primaryKey"`
	UserId    uint   `json:"user_id" gorm:"default:null"`
	Token     string `json:"token" gorm:"index:uidx_fcm_tokens_1,unique;type:varchar(256);not null"`
	IsEnabled bool   `json:"is_enabled" gorm:"not null;default:false"`
}
