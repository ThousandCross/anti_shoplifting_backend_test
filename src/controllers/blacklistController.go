package controllers

import (
	"anti-shoplifting/src/database"
	"anti-shoplifting/src/middlewares"
	"anti-shoplifting/src/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm/clause"
)

func BlackLists(c *fiber.Ctx) error {
	user_id, _ := middlewares.GetUserId(c)

	var blacklists []models.BlackList

	database.DB.Model(&models.BlackList{}).Where("blacklists.user_id = ?", user_id).Find(&blacklists)

	// combine all the cameras' latest visit datetime and nunmber of visits from dynamoDB!!!!!!!!!
	// pending

	return c.JSON(blacklists)
}

func RegistBlackLists(c *fiber.Ctx) error {
	user_id, _ := middlewares.GetUserId(c)

	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}
	// SalesItem regist
	var blacklist models.BlackList
	result := database.DB.Model(&models.BlackList{}).Where("blacklists.id = ? AND blacklists.user_id = ?", data["id"], user_id).First(&blacklist)
	if result.Error != nil {
		// for Creating a new record
		blacklist.UserId = user_id
		blacklist.Name = data["name"]
	}

	//  create a record if it does not exists, and update if it exists
	database.DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}}, // key colume
		DoUpdates: clause.AssignmentColumns([]string{
			"user_id",
			"name",
		}),
	}).Create(&blacklist)

	// how to associate blacklist with incidents???????????????????
	// pending

	//go database.ClearCache("products_frontend", "products_backend")

	return c.JSON(blacklist)
}

func GetBlackList(c *fiber.Ctx) error {
	user_id, _ := middlewares.GetUserId(c)

	id, _ := strconv.Atoi(c.Params("id"))

	var blacklist models.BlackList

	// validation
	result := database.DB.Model(&models.Camera{}).Where("blacklists.id = ? AND blacklists.user_id = ?", id, user_id).First(&blacklist)
	if result.Error != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid Operations!",
		})
	}

	return c.JSON(&blacklist)
}

func DeleteBlackList(c *fiber.Ctx) error {
	user_id, _ := middlewares.GetUserId(c)

	id, _ := strconv.Atoi(c.Params("id"))

	var blacklist models.BlackList

	// validation
	result := database.DB.Model(&models.Camera{}).Where("blacklists.id = ? AND blacklists.user_id = ?", id, user_id).First(&blacklist)
	if result.Error != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid Operations!",
		})
	}

	database.DB.Delete(&blacklist)

	//go database.ClearCache("cameras_frontend", "cameras_backend")

	return nil
}
