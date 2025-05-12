package models

import "golang.org/x/crypto/bcrypt"

type Company struct {
	Id                           uint    `json:"id" gorm:"primaryKey"`
	CompanyCd                    string  `json:"company_cd" gorm:"type:varchar(10);not null"`
	CompanyKey                   []byte  `json:"-" gorm:"not null"`
	Name                         string  `json:"name" gorm:"type:varchar(255);not null"`
	RepresentativeFamilyName     string  `json:"representative_family_name" gorm:"type:varchar(100);not null"`
	RepresentativeFirstName      string  `json:"representative_first_name" gorm:"type:varchar(100);not null"`
	RepresentativeFamilyNameKana string  `json:"representative_family_name_kana" gorm:"type:varchar(100);not null"`
	RepresentativeFirstNameKana  string  `json:"representative_first_name_kana" gorm:"type:varchar(100);not null"`
	Zipcode                      string  `json:"zipcode" gorm:"not null;size:7"`
	PrefectureId                 uint    `json:"prefecture_id" gorm:"not null"`
	City                         string  `json:"city" gorm:"not null"`
	Street                       string  `json:"street" gorm:"not null"`
	Building                     string  `json:"building" gorm:"default null"`
	Tel                          string  `json:"tel" gorm:"not null"`
	Mail                         string  `json:"mail" gorm:"not null"`
	ManagerKey                   []byte  `json:"-" gorm:"not null"`
	ManagerFamilyName            string  `json:"manager_family_name" gorm:"type:varchar(100);not null"`
	ManagerFirstName             string  `json:"manager_first_name" gorm:"type:varchar(100);not null"`
	ManagerFamilyNameKana        string  `json:"manager_family_name_kana" gorm:"type:varchar(100);not null"`
	ManagerFirstNameKana         string  `json:"manager_first_name_kana" gorm:"type:varchar(100);not null"`
	ManagerTel                   string  `json:"manager_tel" gorm:"not null"`
	ManagerMail                  string  `json:"manager_mail" gorm:"not null"`
	ManagerMailVerified          bool    `json:"manager_mail_verified" gorm:"not null;default:false"`
	Store                        []Store `gorm:"foreignKey:CompanyId;references:Id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (company *Company) SetCompanyKey(company_key string) {
	hashedCompanyKey, _ := bcrypt.GenerateFromPassword([]byte(company_key), 12)
	company.CompanyKey = hashedCompanyKey
}

func (company *Company) SetManagerKey(manager_key string) {
	hashedManagerKey, _ := bcrypt.GenerateFromPassword([]byte(manager_key), 12)
	company.ManagerKey = hashedManagerKey
}

func (company *Company) CompareCompanyKey(registeredHashedCompanyKey []byte, CompanyKey []byte) error {
	return bcrypt.CompareHashAndPassword(registeredHashedCompanyKey, CompanyKey)
}
