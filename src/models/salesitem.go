package models

type SalesItem struct {
	Id     uint     `json:"id" gorm:"primaryKey"`
	UserId uint     `json:"user_id" gorm:"index:uidx_salesitem_1,unique;not null"`
	Jan    string   `json:"jan" gorm:"index:uidx_salesitem_1,unique;type:varchar(10);default null"`
	Sku    string   `json:"sku" gorm:"default null"`
	Url    string   `json:"url" gorm:"default null"`
	Sensor []Sensor `gorm:"foreignKey:SalesItemId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
