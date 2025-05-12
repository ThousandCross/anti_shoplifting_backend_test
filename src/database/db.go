package database

import (
	"anti-shoplifting/src/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	var err error

	//jst, _ := time.LoadLocation("Asia/Tokyo")
	// change loc=Local to loc=UTC
	dsn := "docker:password@tcp(db:3306)/anti-shoplifting?charset=utf8mb4&parseTime=True&loc=Local"

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Could not connect with the database!")
	}
}

func AutoMigrate() {
	// Drop table if exists (will ignore or delete foreign key constraints when dropping)
	DB.Migrator().DropTable(
		models.Prefecture{},
		models.Company{},
		models.IncidentType{},
		models.Store{},
		models.User{},
		models.FcmToken{},
		models.SalesItem{},
		models.Camera{},
		models.Sensor{},
		models.Incident{},
		models.IncidentVideo{},
		models.BlackList{},
	)

	// Create New Tables
	DB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").AutoMigrate(
		models.Prefecture{},
		models.Company{},
		models.IncidentType{},
		models.Store{},
		models.User{},
		models.FcmToken{},
		models.SalesItem{},
		models.Camera{},
		models.Sensor{},
		models.Incident{},
		models.IncidentVideo{},
		models.BlackList{},
	)

	// Create Default Records
	// 1. prefectures
	prefectures := []models.Prefecture{
		{Id: 1, Name: "北海道"},
		{Id: 2, Name: "青森県"},
		{Id: 3, Name: "岩手県"},
		{Id: 4, Name: "宮城県"},
		{Id: 5, Name: "秋田県"},
		{Id: 6, Name: "山形県"},
		{Id: 7, Name: "福島県"},
		{Id: 8, Name: "茨城県"},
		{Id: 9, Name: "栃木県"},
		{Id: 10, Name: "群馬県"},
		{Id: 11, Name: "埼玉県"},
		{Id: 12, Name: "千葉県"},
		{Id: 13, Name: "東京都"},
		{Id: 14, Name: "神奈川県"},
		{Id: 15, Name: "新潟県"},
		{Id: 16, Name: "富山県"},
		{Id: 17, Name: "石川県"},
		{Id: 18, Name: "福井県"},
		{Id: 19, Name: "山梨県"},
		{Id: 20, Name: "長野県"},
		{Id: 21, Name: "岐阜県"},
		{Id: 22, Name: "静岡県"},
		{Id: 23, Name: "愛知県"},
		{Id: 24, Name: "三重県"},
		{Id: 25, Name: "滋賀県"},
		{Id: 26, Name: "京都府"},
		{Id: 27, Name: "大阪府"},
		{Id: 28, Name: "兵庫県"},
		{Id: 29, Name: "奈良県"},
		{Id: 30, Name: "和歌山県"},
		{Id: 31, Name: "鳥取県"},
		{Id: 32, Name: "島根県"},
		{Id: 33, Name: "岡山県"},
		{Id: 34, Name: "広島県"},
		{Id: 35, Name: "山口県"},
		{Id: 36, Name: "徳島県"},
		{Id: 37, Name: "香川県"},
		{Id: 38, Name: "愛媛県"},
		{Id: 39, Name: "高知県"},
		{Id: 40, Name: "福岡県"},
		{Id: 41, Name: "佐賀県"},
		{Id: 42, Name: "長崎県"},
		{Id: 43, Name: "熊本県"},
		{Id: 44, Name: "大分県"},
		{Id: 45, Name: "宮崎県"},
		{Id: 46, Name: "鹿児島県"},
		{Id: 47, Name: "沖縄県"},
	}
	DB.Create(&prefectures)

	// 2. incident type
	incident_types := []models.IncidentType{
		{Id: 1, Name: "万引き"},
	}
	DB.Create(&incident_types)
}
