package controllers

import (
	"anti-shoplifting/src/database"
	"anti-shoplifting/src/middlewares"
	"anti-shoplifting/src/models"
	dynamomodels "anti-shoplifting/src/models/dynamo"
	"context"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slices"
	"gorm.io/gorm/clause"
)

func GetUserIdBySensorCd(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	// Sensor regist
	var sensor models.Sensor
	result := database.DB.Model(&models.Sensor{}).Where("sensors.sensor_cd = ?", data["sensor_cd"]).First(&sensor)
	if result.Error != nil {
		// error
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"result":  "ng",
			"message": "Sensor not found!",
			"user_id": "",
		})
	}

	return c.JSON(fiber.Map{
		"result":  "ok",
		"message": "",
		"user_id": sensor.UserId,
	})
}

func InitializeSensorPosition(c *fiber.Ctx) error {
	type Data struct {
		SensorCd  string  `json:"sensor_cd"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Elevation float64 `json:"elevation"`
	}
	var data Data

	if err := c.BodyParser(&data); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"result":  "ng",
			"message": "Can't parse request parameter",
		})
	}

	var sensor models.Sensor
	result := database.DB.Model(&models.Sensor{}).Where("sensors.sensor_cd = ?", data.SensorCd).First(&sensor)
	if result.Error != nil {
		// error
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"result":  "ng",
			"message": "Sensor not found!",
		})
	}

	sensor.Latitude = data.Latitude
	sensor.Longitude = data.Longitude
	sensor.Elevation = data.Elevation
	database.DB.Save(&sensor)

	return c.JSON(fiber.Map{
		"result":  "ok",
		"message": "",
	})
}

func Sensors(c *fiber.Ctx) error {
	user_id, _ := middlewares.GetUserId(c)

	type Moniterings struct {
		SensorCd string `json:"sensor_cd"`
		UserId   uint
		//CreatedAt time.Time `json:"created_at"`
		Url string `json:"url"`
	}
	var monitorings []Moniterings

	database.DB.Model(&models.Sensor{}).Select(
		"sensors.sensor_cd",
		"sensors.user_id",
		//"sensors.created_at",
		"sales_items.url",
	).Joins(
		"inner join sales_items on sales_items.id = sensors.sales_item_id",
	).Where("sensors.user_id = ?", user_id).Order("sensors.sensor_cd").Find(&monitorings)

	// combine all the sensors' status from dynamoDB!!!!!!!!!
	var ctx = context.Background()
	out, err := database.DynamoDB.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String("LogSensor"),
		IndexName:              aws.String("UserId-index"),
		KeyConditionExpression: aws.String("UserId = :userId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":userId": &types.AttributeValueMemberN{Value: strconv.FormatUint(uint64(user_id), 10)}, // uint to string
		},
	})
	if err != nil {
		panic(err)
	}

	var logs []dynamomodels.LogSensor
	attributevalue.UnmarshalListOfMaps(out.Items, &logs)

	sort.SliceStable(logs, func(i, j int) bool { return logs[i].SensorCd < logs[j].SensorCd })
	fmt.Printf("SensorCdでSort(Stable):%+v\n", logs)

	type Result struct {
		SensorCd  string    `json:"sensor_cd"`
		Status    bool      `json:"status"`
		Url       string    `json:"url"`
		CreatedAt time.Time `json:"created_at"`
	}

	var results []Result
	var result Result

	// dummy created_at
	now := time.Now()
	nowUTC := now.UTC()
	fmt.Println(nowUTC.Format(time.RFC3339))
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	nowJST := nowUTC.In(jst)

	for _, monitoring := range monitorings {
		// result.SensorCd
		result.SensorCd = monitoring.SensorCd
		// result.Status
		idx := slices.IndexFunc(logs, func(c dynamomodels.LogSensor) bool { return c.SensorCd == monitoring.SensorCd })
		if idx < 0 { // sensor.sensor_cd was not found in dynamoDB
			result.Status = false
		} else {
			diff := nowJST.Sub(logs[idx].LatestTimestamp)
			if diff.Hours() >= 2 {
				result.Status = false
			} else {
				result.Status = true
			}
		}
		// result.Url
		result.Url = monitoring.Url
		// result.CreatedAt = camera.CreatedAt
		result.CreatedAt = nowJST

		results = append(results, result)
	}

	//return c.JSON(cameras)
	return c.JSON(results)
}

func RegistSensors(c *fiber.Ctx) error {
	user_id, _ := middlewares.GetUserId(c)

	type Data struct {
		SensorCd  string  `json:"sensor_cd"`
		Jan       string  `json:"jan"`
		IpAddress string  `json:"ip_address"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Elevation float64 `json:"elevation"`
	}
	var data Data

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	// SalesItem regist
	var salesitem models.SalesItem
	result_sales_item := database.DB.Model(&models.SalesItem{}).Where("sales_items.user_id = ? AND sales_items.jan = ?", user_id, data.Jan).First(&salesitem)
	if result_sales_item.Error != nil {
		// create a new record
		salesitem.UserId = user_id
		salesitem.Jan = data.Jan
		database.DB.Create(&salesitem)
	}
	sales_item_Id := salesitem.Id

	// Sensor regist
	var sensor models.Sensor
	result_sensor := database.DB.Model(&models.Sensor{}).Where("sensors.sensor_cd = ?", data.SensorCd).First(&sensor)
	if result_sensor.Error != nil {
		// for Creating a new record
		sensor.SensorCd = data.SensorCd
		sensor.UserId = user_id
	}

	// required item
	sensor.SalesItemId = sales_item_Id

	// optional item
	if data.IpAddress != "" {
		sensor.IpAddress = data.IpAddress
	}
	if data.Latitude != 0 {
		sensor.Latitude = data.Latitude
	}
	if data.Longitude != 0 {
		sensor.Longitude = data.Longitude
	}
	if data.Elevation != 0 {
		sensor.Elevation = data.Elevation // data["elevation"]
	}

	//  create a record if it does not exists, and update if it exists
	database.DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "sensor_cd"}}, // key colume
		DoUpdates: clause.AssignmentColumns([]string{
			"user_id",
			"ip_address",
			"latitude",
			"longitude",
			"elevation",
			"sales_item_id",
		}),
	}).Create(&sensor)

	//go database.ClearCache("products_frontend", "products_backend")

	return c.JSON(sensor)
}

func GetSensor(c *fiber.Ctx) error {
	user_id, _ := middlewares.GetUserId(c)

	sensor_cd, _ := strconv.Atoi(c.Params("sensor_cd"))

	var sensor models.Sensor
	result := database.DB.Model(&models.Sensor{}).Where("sensors.sensor_cd = ? AND sensors.user_id = ?", sensor_cd, user_id).First(&sensor)
	if result.Error != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid Operations!",
		})
	}

	return c.JSON(sensor)
}

func DeleteSensor(c *fiber.Ctx) error {
	user_id, _ := middlewares.GetUserId(c)

	sensor_cd, _ := strconv.Atoi(c.Params("sensor_cd"))

	var sensor models.Sensor
	result := database.DB.Model(&models.Camera{}).Where("sensors.sensor_cd = ? AND sensors.user_id = ?", sensor_cd, user_id).First(&sensor)
	if result.Error != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid Operations!",
		})
	}

	database.DB.Delete(&sensor)

	//go database.ClearCache("cameras_frontend", "cameras_backend")

	return nil
}
