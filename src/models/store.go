package models

import "golang.org/x/crypto/bcrypt"

type Store struct {
	Id                    uint   `json:"id" gorm:"primaryKey"`
	StoreCd               string `json:"store_cd" gorm:"index:uidx_store_1,unique;type:varchar(10);not null"`
	StoreKey              []byte `json:"-" gorm:"not null"`
	CompanyId             uint   `json:"company_id" gorm:"index:uidx_store_1,unique;not null"`
	IsApproved            bool   `json:"is_approved" gorm:"not null;default:false"`
	IsAdmin               bool   `json:"is_admin" gorm:"not null;default:false"`
	Name                  string `json:"name" gorm:"not null"`
	Zipcode               string `json:"zipcode" gorm:"not null;size:7"`
	PrefectureId          uint   `json:"prefecture_id" gorm:"not null"`
	City                  string `json:"city" gorm:"not null"`
	Street                string `json:"street" gorm:"not null"`
	Building              string `json:"building" gorm:"default null"`
	ManagerFamilyName     string `json:"manager_family_name" gorm:"type:varchar(100);not null"`
	ManagerFirstName      string `json:"manager_first_name" gorm:"type:varchar(100);not null"`
	ManagerFamilyNameKana string `json:"manager_family_name_kana" gorm:"type:varchar(100);not null"`
	ManagerFirstNameKana  string `json:"manager_first_name_kana" gorm:"type:varchar(100);not null"`
	ManagerTel            string `json:"manager_tel" gorm:"not null"`
	ManagerMail           string `json:"manager_mail" gorm:"not null"`
	ManagerMailVerified   bool   `json:"manager_mail_verified" gorm:"not null;default:false"`
	User                  User   `gorm:"foreignKey:Id;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (store *Store) SetStoreKey(store_key string) {
	hashedStoreKey, _ := bcrypt.GenerateFromPassword([]byte(store_key), 12)
	store.StoreKey = hashedStoreKey
}

func (store *Store) CompareStoreKey(registeredHashedStoreKey []byte, StoreKey []byte) error {
	return bcrypt.CompareHashAndPassword(registeredHashedStoreKey, StoreKey)
}
