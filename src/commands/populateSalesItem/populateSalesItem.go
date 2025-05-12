package main

import (
	"anti-shoplifting/src/database"
	"anti-shoplifting/src/models"
)

func main() {
	database.Connect()

	salesitems := []models.SalesItem{
		{UserId: 1, Jan: "1001", Sku: "", Url: ""},
		{UserId: 1, Jan: "1002", Sku: "", Url: ""},
		{UserId: 1, Jan: "1003", Sku: "", Url: "https://icooon-mono.com/i/icon_13583/icon_135830.svg"},
		{UserId: 1, Jan: "1004", Sku: "", Url: ""},
		{UserId: 1, Jan: "1005", Sku: "", Url: "https://icooon-mono.com/i/icon_13585/icon_135850.svg"},
		{UserId: 1, Jan: "1006", Sku: "", Url: "https://icooon-mono.com/i/icon_13586/icon_135860.svg"},
		{UserId: 1, Jan: "1007", Sku: "", Url: "https://icooon-mono.com/i/icon_13587/icon_135870.svg"},
		{UserId: 1, Jan: "1008", Sku: "", Url: "https://icooon-mono.com/i/icon_13588/icon_135880.svg"},
		{UserId: 1, Jan: "1009", Sku: "", Url: "https://icooon-mono.com/i/icon_13589/icon_135890.svg"},
		{UserId: 1, Jan: "1010", Sku: "", Url: "https://icooon-mono.com/i/icon_13590/icon_135900.svg"},
	}

	database.DB.Create(&salesitems)

}
