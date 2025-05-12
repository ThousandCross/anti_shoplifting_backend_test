package controllers

import (
	"anti-shoplifting/src/database"
	"anti-shoplifting/src/middlewares"
	"anti-shoplifting/src/models"
	"anti-shoplifting/src/utils"
	"fmt"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func RegistCompany(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		fmt.Printf("%+v\n", data)
		fmt.Printf("%+v\n", err)

		return c.JSON(fiber.Map{
			"result":  "ng",
			"message": "parse ng",
			//"err":         err,
			"company_id":  "",
			"company_cd":  "",
			"company_key": "",
			"manager_key": "",
		})
	}

	fmt.Printf("%+v\n", data)
	// validation

	// DB registration
	prefecture_id, _ := strconv.Atoi(data["prefecture_id"])
	random_company_cd, _ := utils.MakeRandomStr(7)
	company_cd := "crp" + random_company_cd
	company_key, _ := utils.MakeRandomStr(8)
	manager_key, _ := utils.MakeRandomStr(8)

	company := models.Company{
		CompanyCd:                    company_cd,
		Name:                         data["company_name"],
		RepresentativeFamilyName:     data["representative_family_name"],
		RepresentativeFirstName:      data["representative_first_name"],
		RepresentativeFamilyNameKana: data["representative_family_name_kana"],
		RepresentativeFirstNameKana:  data["representative_first_name_kana"],
		Zipcode:                      data["zipcode"],
		PrefectureId:                 uint(prefecture_id),
		City:                         data["city"],
		Street:                       data["street"],
		Building:                     data["building"],
		Tel:                          data["tel"],
		Mail:                         data["mail"],
		ManagerFamilyName:            data["manager_family_name"],
		ManagerFirstName:             data["manager_first_name"],
		ManagerFamilyNameKana:        data["manager_family_name_kana"],
		ManagerFirstNameKana:         data["manager_first_name_kana"],
		ManagerTel:                   data["manager_tel"],
		ManagerMail:                  data["manager_mail"],
		ManagerMailVerified:          false,
	}
	company.SetCompanyKey(company_key)
	company.SetManagerKey(manager_key)
	database.DB.Create(&company)

	company_id := company.Id

	fmt.Printf("company_id: %d\n", company_id)
	fmt.Printf("company_cd: %s\n", company_cd)
	fmt.Printf("company_key: %s\n", company_key)
	fmt.Printf("manager_key: %s\n", manager_key)

	// mail to manager of company for notifying a url including company_id(as key1) and hashed company key(as key2)
	// but url is not this route but web app url
	// !pending!

	return c.JSON(fiber.Map{
		"result":              "ok",
		"message":             "company registeration completed!!",
		"company_id":          company_id,
		"company_cd":          company_cd,
		"company_key":         company_key,
		"company_manager_key": manager_key,
	})
}

func VerifyCompany(c *fiber.Ctx) error {
	// If there is no query string, it returns an empty string.
	id := c.Query("key1")
	company_key := c.Query("key2")

	// validation
	var company models.Company
	result := database.DB.Model(&models.Company{}).Where("companies.id = ? and companies.company_key = ?", id, company_key).First(&company)
	if result.Error != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"result":     "ng",
			"message":    "Invalid Operations!",
			"company_cd": "",
		})
	}
	if company.ManagerMailVerified == true {
		// error due to already verified
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"result":     "ng",
			"message":    "Invalid Operations!",
			"company_cd": "",
		})
	}

	// Update record
	database.DB.Model(&company).Update("manager_mail_verified", true)

	// display company_cd on screen od web app
	return c.JSON(fiber.Map{
		"result":     "ok",
		"message":    "company registeration completed!!",
		"company_cd": company.CompanyCd,
	})
}

func RegistStore(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	// validation
	// company_cd
	var company models.Company
	result := database.DB.Model(&models.Company{}).Where("companies.company_cd = ?", data["company_cd"]).First(&company)
	if result.Error != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"result":    "ng",
			"message":   "Invalid Parameters!",
			"store_id":  "",
			"store_cd":  "",
			"store_key": "",
		})
	}
	// company_key
	if _, result := data["company_key"]; !result {
		c.Status(400)
		return c.JSON(fiber.Map{
			"result":    "ng",
			"message":   "Invalid Parameters!",
			"store_id":  "",
			"store_cd":  "",
			"store_key": "",
		})
	}
	if err := company.CompareCompanyKey(company.CompanyKey, []byte(data["company_key"])); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"result":    "ng",
			"message":   "Invalid Credentials!",
			"store_id":  "",
			"store_cd":  "",
			"store_key": "",
		})
	}
	// password, password_confirm
	if data["password"] != data["password_confirm"] {
		c.Status(400)
		return c.JSON(fiber.Map{
			"result":    "ng",
			"message":   "passwords do not match",
			"store_id":  "",
			"store_cd":  "",
			"store_key": "",
		})
	}

	// DB registration
	prefecture_id, _ := strconv.Atoi(data["prefecture_id"])
	random_store_cd, _ := utils.MakeRandomStr(7)
	store_cd := "str" + random_store_cd
	store_key, _ := utils.MakeRandomStr(8)

	store := models.Store{
		StoreCd:               store_cd,
		CompanyId:             company.Id,
		IsApproved:            false,
		IsAdmin:               false,
		Name:                  data["name"],
		Zipcode:               data["zipcode"],
		PrefectureId:          uint(prefecture_id),
		City:                  data["city"],
		Street:                data["street"],
		Building:              data["building"],
		ManagerFamilyName:     data["manager_family_name"],
		ManagerFirstName:      data["manager_first_name"],
		ManagerFamilyNameKana: data["manager_family_name_kana"],
		ManagerFirstNameKana:  data["manager_first_name_kana"],
		ManagerTel:            data["manager_tel"],
		ManagerMail:           data["manager_mail"],
		ManagerMailVerified:   false,
	}
	store.SetStoreKey(store_key)
	database.DB.Create(&store)

	user_id := store.Id
	user := models.User{
		Id: user_id,
	}
	user.SetPassword(data["password"])
	database.DB.Create(&user)

	// mail to manager of store for notifying a url including store_id(as key1) and hashed store key(as key2)
	// but url is not this route but web app url
	// !pending!

	// mail to manager to company for notifying a url including
	// store_id(as key1) and
	// store_key(as key2) and
	// hashed company_key(as key3)
	// but url is not this route but web app url
	// !pending!

	fmt.Printf("store_id: %d\n", user_id)
	fmt.Printf("store_cd: %s\n", store_cd)
	fmt.Printf("store_key: %s\n", store_key)

	// If manager of company approved store regist request,
	// system will notify store_cd and company_cd to manager of store
	return c.JSON(fiber.Map{
		"result":    "ok",
		"message":   "company registeration completed!!",
		"store_id":  user_id,
		"store_cd":  store_cd,
		"store_key": store_key,
	})
}

func ApproveStore(c *fiber.Ctx) error {
	// If there is no query string, it returns an empty string.
	store_id := c.Query("key1")
	store_key := c.Query("key2")
	company_key := c.Query("key3")

	//var result Result
	var store models.Store
	database.DB.Model(&models.Store{}).Select(
		"stores.*",
	).Joins(
		"inner join companies on stores.company_id = companies.id",
	).Where(
		"stores.id = ? AND stores.store_key = ? AND companies.company_key = ?",
		store_id,
		store_key,
		company_key,
	).First(&store)
	if store.Id == 0 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"result":   "ng",
			"message":  "Invalid Credentials!",
			"store_cd": "",
		})
	}
	if store.IsApproved == true {
		// error due to already verified
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"result":   "ng",
			"message":  "Invalid Operations!",
			"store_cd": "",
		})
	}

	// Update record
	database.DB.Model(&store).Update("is_approved", true)

	// mail to manager of store for notifying company_cd and store_cd
	// !pending!

	// display store information fo newly added store to company manager.
	return c.JSON(fiber.Map{
		"result":   "ok",
		"message":  "approval of new store registration was completed!!",
		"store_cd": store.StoreCd,
	})
}

func VerifyStore(c *fiber.Ctx) error {
	// If there is no query string, it returns an empty string.
	store_id := c.Query("key1")
	store_key := c.Query("key2")

	// validation
	var store models.Store
	result := database.DB.Model(&models.Store{}).Where("stores.id = ? and stores.store_key = ?", store_id, store_key).First(&store)
	if result.Error != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"result":  "ng",
			"message": "Invalid Operations!",
		})
	}
	if store.ManagerMailVerified == true {
		// error due to already verified
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"result":  "ng",
			"message": "Invalid Operations!",
		})
	}

	// Update record
	database.DB.Model(&store).Update("manager_mail_verified", true)

	// display notification of mail verifying completed.
	// Do not notify store_cd until store approval completed.
	return c.JSON(fiber.Map{
		"result":  "ok",
		"message": "store email verification completed!",
	})
}

func ResetStorePassword(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		fmt.Printf("%+v\n", data)
		fmt.Printf("%+v\n", err)

		return c.JSON(fiber.Map{
			"result":  "ng",
			"message": "parse ng",
		})
	}

	// validation
	// company_cd
	var company models.Company
	resultCompany := database.DB.Model(&models.Company{}).Where("companies.company_cd = ?", data["company_cd"]).First(&company)
	if resultCompany.Error != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"result":  "ng",
			"message": "Cannot find Company!",
		})
	}
	// store_cd
	if _, result := data["store_cd"]; !result {
		c.Status(400)
		return c.JSON(fiber.Map{
			"result":  "ng",
			"message": "Store Code doesn't exists in post Parameters!",
		})
	}
	var store models.Store
	resultStore := database.DB.Model(&models.Store{}).
		Where(
			"stores.store_cd = ? AND stores.company_id = ?",
			data["store_cd"],
			company.Id,
		).First(&store)
	if resultStore.Error != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"result":  "ng",
			"message": "Cannot find Store!",
		})
	}
	// store_key
	if _, result := data["store_key"]; !result {
		c.Status(400)
		return c.JSON(fiber.Map{
			"result":  "ng",
			"message": "Store key doesn't exists in post Parameters!",
		})
	}
	if err := store.CompareStoreKey(store.StoreKey, []byte(data["store_key"])); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"result":  "ng",
			"message": "Store key doesn't match with stored information!",
		})
	}
	// old_password
	var user models.User
	resultUser := database.DB.Model(&models.User{}).Where("users.id = ?", store.Id).First(&user)
	if resultUser.Error != nil {
		c.Status(400)
		return c.JSON(fiber.Map{
			"result":  "ng",
			"message": "Cannot find User!",
		})
	}
	if err := models.ComparePassword(user.Password, []byte(data["old_password"])); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"result":  "ng",
			"message": "Old password doesn't match with stored information!",
		})
	}

	// new_password, new_password_confirm
	if data["new_password"] != data["new_password_confirm"] {
		c.Status(400)
		return c.JSON(fiber.Map{
			"result":  "ng",
			"message": "New password doesn't match New Password Confirm",
		})
	}

	// DB registration
	// update users.password
	user.SetPassword(data["new_password"])
	database.DB.Save(&user)

	// mail to manager of store for notifying resetting password was completed!
	//pending

	// If manager of company approved store regist request,
	// system will notify store_cd and company_cd to manager of store
	return c.JSON(fiber.Map{
		"result":  "ok",
		"message": "reseting password completed!!",
	})

}

func RegistFcmToken(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	fcm_token := data["fcm_token"]
	if len(fcm_token) == 0 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"result":  "ng",
			"message": "Invalid Credentials!",
		})
	}

	var token models.FcmToken
	result := database.DB.Model(&models.FcmToken{}).Where("fcm_tokens.token = ?", fcm_token).First(&token)
	//if result.RowsAffected > 0 {
	if result.RowsAffected == 0 {
		token = models.FcmToken{
			//UserId:    NULL,
			Token:     fcm_token,
			IsEnabled: true,
		}
		database.DB.Create(&token)
		fmt.Printf("fcm_token create complete!!")

		// c.Status(fiber.StatusBadRequest)
		// return c.JSON(fiber.Map{
		// 	"result":  "ng",
		// 	"message": "Invalid Credentials!",
		// })
	}

	// create a new record(IsEnabled = true)
	// token = models.FcmToken{
	// 	//UserId:    NULL,
	// 	Token:     fcm_token,
	// 	IsEnabled: true,
	// }
	// database.DB.Create(&token)
	// fmt.Printf("fcm_token create complete!!")

	return c.JSON(fiber.Map{
		"result":  "ok",
		"message": "fcm token registration succeeded!",
	})

}

func RefreshFcmToken(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	old_fcm_token := data["old_fcm_token"]
	new_fcm_token := data["new_fcm_token"]
	if len(old_fcm_token) == 0 || len(new_fcm_token) == 0 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"result":  "ng",
			"message": "Invalid Credentials!",
		})
	}

	// update an old record (IsEnabled = false)
	var token_old models.FcmToken
	result_old := database.DB.Model(&models.FcmToken{}).Where("fcm_tokens.token = ?", old_fcm_token).First(&token_old)
	if result_old.Error != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"result":  "ng",
			"message": "Invalid Credentials!",
		})
	}
	token_old.IsEnabled = false
	database.DB.Model(&token_old).Update("is_enabled", false)

	// create a new record(IsEnabled = true)
	var token_new models.FcmToken
	result_new := database.DB.Model(&models.FcmToken{}).Where("fcm_tokens.token = ?", new_fcm_token).First(&token_new)
	if result_new.RowsAffected > 0 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"result":  "ng",
			"message": "Invalid Credentials!",
		})
	}
	token_new = models.FcmToken{
		//UserId:    NULL,
		Token:     new_fcm_token,
		IsEnabled: true,
	}
	database.DB.Create(&token_new)

	return c.JSON(fiber.Map{
		"result": "ok",
	})

}

func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	// response structure definition
	type Result struct {
		Id        uint
		Password  []byte
		StoreCd   string
		CompanyId uint
	}

	var result Result
	database.DB.Model(&models.User{}).Select(
		"users.id",
		"users.password",
		"stores.store_cd",
		"stores.company_id",
	).Joins(
		"inner join stores on stores.id = users.id",
	).Joins(
		"inner join companies on companies.id = stores.company_id",
	).Where(
		"companies.company_cd = ? AND stores.store_cd = ?",
		data["company_cd"],
		data["store_cd"],
	).First(&result)
	if result.Id == 0 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"result":  "ng",
			"message": "Invalid Credentials!",
		})
	}

	if err := models.ComparePassword(result.Password, []byte(data["password"])); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"result":  "ng",
			"message": "Invalid Credentials!",
		})
	}

	// fcm_token_registration
	fcm_token := data["fcm_token"]
	if len(fcm_token) == 0 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"result":  "ng",
			"message": "Invalid Credentials!",
		})
	}

	// update use id of fcm token record
	var token models.FcmToken
	result_token := database.DB.Model(&models.FcmToken{}).Where("fcm_tokens.token = ?", fcm_token).First(&token)
	if result_token.Error != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"result":  "ng",
			"message": "Invalid Credentials!",
		})
	}
	database.DB.Model(&token).Update("user_id", result.Id)

	// return jwt token
	if data["from"] == "web" {

		payload := jwt.StandardClaims{
			Subject:   strconv.Itoa(int(result.Id)),
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		}

		token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, payload).SignedString([]byte("secret"))
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"result":  "ng",
				"message": "Invalid Credentials!",
			})
		}

		cookie := fiber.Cookie{
			Name:     "jwt",
			Value:    token,
			Expires:  time.Now().Add(time.Hour * 24),
			HTTPOnly: true,
		}

		c.Cookie(&cookie)

		return c.JSON(fiber.Map{
			"result":  "ok",
			"message": "success",
		})

	} else if data["from"] == "mobile" {
		payload := jwt.MapClaims{
			"user_id": strconv.Itoa(int(result.Id)),
			"exp":     time.Now().Add(time.Hour * 24).Unix(),
		}

		token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, payload).SignedString([]byte("secret"))
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"result":  "ng",
				"message": "Invalid Credentials!",
			})
		}

		return c.JSON(fiber.Map{
			"result":  "ok",
			"message": "login from web succeeded!",
			"jwt":     token,
		})

	} else {
		return c.JSON(fiber.Map{
			"result":  "ng",
			"message": "Invalid Credentials!",
		})

	}
}

func User(c *fiber.Ctx) error {
	user_id, _ := middlewares.GetUserId(c)

	// response structure definition
	type Result struct {
		Id        uint
		StoreCd   string
		CompanyCd string
	}

	var result Result
	database.DB.Model(&models.User{}).Select(
		"user.id",
		"stores.company_cd",
		"stores.store_cd",
	).Joins(
		"inner join stores on stores.id = user.id",
	).Where(
		"users.user_id = ?",
		user_id,
	).First(&result)
	if result.Id == 0 {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"result":  "ng",
			"message": "Invalid Credentials!",
		})
	}

	return c.JSON(result)
}

func Logout(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	if data["from"] == "web" {
		// case1. accessed from web app
		cookie := fiber.Cookie{
			Name:     "jwt",
			Value:    "",
			Expires:  time.Now().Add(-time.Hour),
			HTTPOnly: true,
		}

		c.Cookie(&cookie)
	} else if data["from"] == "mobile" {
		// case2. accessed from mobile app
		// mobile app destroys jwt token.

	} else {
		// case3. accessed from others
		return c.JSON(fiber.Map{
			"result":  "ng",
			"message": "Invalid Credentials!",
		})
	}

	return c.JSON(fiber.Map{
		"result":  "ok",
		"message": "success",
	})
}

func UpdateInfo(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	id, _ := middlewares.GetUserId(c)

	prefecture_id, _ := strconv.Atoi(data["prefecture_Id"])
	store := models.Store{
		Id:                id,
		Name:              data["name"],
		Zipcode:           data["zipcode"],
		PrefectureId:      uint(prefecture_id),
		City:              data["city"],
		Street:            data["street"],
		Building:          data["building"],
		ManagerFamilyName: data["manager_family_name"],
		ManagerFirstName:  data["manager_first_name"],
		ManagerTel:        data["manager_tel"],
		ManagerMail:       data["manager_mail"],
	}

	database.DB.Model(&models.Store{}).Updates(&store)

	return c.JSON(store)
}

func UpdatePassword(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	// validation
	if data["password"] != data["password_confirm"] {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "passwords do not match",
		})
	}

	id, _ := middlewares.GetUserId(c)

	user := models.User{
		Id: id,
	}
	user.SetPassword(data["password"])

	database.DB.Model(&models.User{}).Updates(&user)

	return c.JSON(user)
}
