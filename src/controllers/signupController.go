package controllers

import (
	"anti-shoplifting/src/database"
	"anti-shoplifting/src/models"

	"github.com/gofiber/fiber/v2"
)

func GetPrefectures(c *fiber.Ctx) error {

	var prefectures []models.Prefecture

	// result := database.DB.Model(&models.Prefecture{}).Order("prefectures.id").Find(&prefectures)
	// if result.Error != nil {
	// 	c.Status(fiber.StatusBadRequest)
	// 	return c.JSON(fiber.Map{
	// 		"message": "Invalid Operations!",
	// 	})
	// }

	// database.DB.Model(&prefectures).Association("Store").Clear()
	// database.DB.Model(&prefectures).Association("Company").Clear()
	// database.DB.Model(&prefectures).Association("Store").Delete([]models.Store{})
	// database.DB.Model(&prefectures).Association("Company").Delete([]models.Company{})

	//database.DB.Select("id", "name").Order("prefectures.id").Find(&prefectures).Association().Replace()
	database.DB.Distinct("id", "name").Order("id").Find(&prefectures)

	return c.JSON(prefectures)
}
