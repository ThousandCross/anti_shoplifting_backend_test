package models

import "golang.org/x/crypto/bcrypt"

type User struct {
	Id        uint        `json:"id" gorm:"primaryKey;autoIncrement:false"`
	Password  []byte      `json:"-" gorm:"not null"`
	FcmToken  []FcmToken  `gorm:"foreignKey:UserId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	SalesItem []SalesItem `gorm:"foreignKey:UserId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Camera    []Camera    `gorm:"foreignKey:UserId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Sensor    []Sensor    `gorm:"foreignKey:UserId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Incident  []Incident  `gorm:"foreignKey:UserId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	BlackList []BlackList `gorm:"foreignKey:UserId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (user *User) SetPassword(password string) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), 12)
	user.Password = hashedPassword
}

func ComparePassword(registeredHashedPassword []byte, password []byte) error {
	return bcrypt.CompareHashAndPassword(registeredHashedPassword, password)
}
