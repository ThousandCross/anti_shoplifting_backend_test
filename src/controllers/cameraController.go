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
	"golang.org/x/exp/slices"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm/clause"
)

func Cameras(c *fiber.Ctx) error {
	user_id, _ := middlewares.GetUserId(c)

	// retrive all cameras corresponding with user_id from RDS
	var cameras []models.Camera
	database.DB.Model(&models.Camera{}).Where("cameras.user_id = ?", user_id).Order("cameras.serial_no").Find(&cameras)

	// combine all the cameras' status from dynamoDB!!!!!!!!!
	var ctx = context.Background()
	out, err := database.DynamoDB.Query(ctx, &dynamodb.QueryInput{
		TableName:              aws.String("LogCamera"),
		IndexName:              aws.String("UserId-index"),
		KeyConditionExpression: aws.String("UserId = :userId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":userId": &types.AttributeValueMemberN{Value: strconv.FormatUint(uint64(user_id), 10)},
		},
	})
	if err != nil {
		panic(err)
	}

	var logs []dynamomodels.LogCamera
	attributevalue.UnmarshalListOfMaps(out.Items, &logs)

	sort.SliceStable(logs, func(i, j int) bool { return logs[i].CameraId < logs[j].CameraId })
	fmt.Printf("CameraIdでSort(Stable):%+v\n", logs)

	type Result struct {
		SerialNo  string    `json:"serial_no"`
		Status    bool      `json:"status"`
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

	//fmt.Println("result empty")
	for _, camera := range cameras {
		fmt.Println("result range in!!!")
		// result.SerialNo
		result.SerialNo = camera.SerialNo
		// result.Status
		idx := slices.IndexFunc(logs, func(c dynamomodels.LogCamera) bool { return uint(c.CameraId) == camera.Id })
		//fmt.Printf("idx: %d\n", idx)
		if idx < 0 { // camera.Id was not found in dynamoDB
			//fmt.Printf("camera.Id was not found in dynamoDB: %d\n", camera.Id)

			result.Status = false
		} else {
			// fmt.Printf("camera.Id was found in dynamoDB!!!!: %d\n", camera.Id)
			// fmt.Printf("DynamoDB Latest Timestamp:%+v\n", logs[idx].LatestTimestamp)
			// fmt.Printf("nowJST:%+v\n", nowJST)
			diff := nowJST.Sub(logs[idx].LatestTimestamp)
			if diff.Hours() >= 2 {
				result.Status = false
			} else {
				result.Status = true
			}
		}
		// result.CreatedAt = camera.CreatedAt
		result.CreatedAt = nowJST

		results = append(results, result)
	}

	//return c.JSON(cameras)
	return c.JSON(results)
}

func RegistCameras(c *fiber.Ctx) error {
	user_id, _ := middlewares.GetUserId(c)

	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return c.JSON(fiber.Map{
			"result":  "ng",
			"message": "parse ng",
			"err":     err,
		})
	}

	var camera models.Camera
	result := database.DB.Model(&models.Camera{}).Where("cameras.serial_no = ? AND cameras.user_id = ?", data["serial_no"], user_id).First(&camera)
	if result.Error != nil {
		// for Creating a new record
		camera.SerialNo = data["serial_no"]
		camera.UserId = user_id
	}

	// optional item
	if val, result := data["ip_address"]; result {
		camera.IpAddress = val // data["ip_address"]
	}
	if val, result := data["latitude"]; result {
		camera.Latitude = val // data["latitude"]
	}
	if val, result := data["longitude"]; result {
		camera.Longitude = val // data["longitude"]
	}
	if val, result := data["elevation"]; result {
		camera.Elevation = val // data["elevation"]
	}

	//  create a record if it does not exists, and update if it exists
	database.DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}}, // key colume
		DoUpdates: clause.AssignmentColumns([]string{
			"serial_no",
			"user_id",
			"ip_address",
			"latitude",
			"longitude",
			"elevation",
		}),
	}).Create(&camera)

	//go database.ClearCache("products_frontend", "products_backend")

	return c.JSON(camera)
}

func GetCamera(c *fiber.Ctx) error {
	user_id, _ := middlewares.GetUserId(c)

	serial_no, _ := strconv.Atoi(c.Params("serial_no"))

	var camera models.Camera
	result := database.DB.Model(&models.Camera{}).Where("cameras.serial_no = ? AND cameras.user_id = ?", serial_no, user_id).First(&camera)
	if result.Error != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid Operations!",
		})
	}

	return c.JSON(camera)
}

func DeleteCamera(c *fiber.Ctx) error {
	user_id, _ := middlewares.GetUserId(c)

	serial_no, _ := strconv.Atoi(c.Params("serial_no"))

	var camera models.Camera
	result := database.DB.Model(&models.Camera{}).Where("cameras.serial_no = ? AND cameras.user_id = ?", serial_no, user_id).First(&camera)
	if result.Error != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "Invalid Operations!",
		})
	}

	database.DB.Delete(&camera)

	//go database.ClearCache("cameras_frontend", "cameras_backend")

	return nil
}
