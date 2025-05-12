package models

type Sensor struct {
	SensorCd    string  `json:"sensor_cd" gorm:"primaryKey;type:varchar(50);not null"`
	UserId      uint    `json:"user_id" gorm:"not null"`
	IpAddress   string  `json:"ip_address" gorm:"default null"` // IPアドレス
	Latitude    float64 `json:"latitude" gorm:"default null"`   // 緯度
	Longitude   float64 `json:"longitude" gorm:"default null"`  // 経度
	Elevation   float64 `json:"elevation" gorm:"default null"`  // 標高
	SalesItemId uint    `json:"sales_item_id" gorm:"default null"`
}
