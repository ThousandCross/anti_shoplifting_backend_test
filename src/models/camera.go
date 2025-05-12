package models

type Camera struct {
	Id            uint            `json:"id" gorm:"primaryKey"` // as dynamo.primary_key
	SerialNo      string          `json:"serial_no" gorm:"index:uidx_camera_1,unique;not null"` // permit serial.no dupication between stores
	UserId        uint            `json:"user_id" gorm:"index:uidx_camera_1,unique;not null"`
	IpAddress     string          `json:"ip_address" gorm:"default null"` // IPアドレス
	Latitude      string          `json:"latitude" gorm:"default null"`   // 緯度
	Longitude     string          `json:"longitude" gorm:"default null"`  // 経度
	Elevation     string          `json:"elevation" gorm:"default null"`  // 標高
	IncidentVideo []IncidentVideo `gorm:"foreignKey:CameraId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
