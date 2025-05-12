package routes

import (
	"anti-shoplifting/src/controllers"
	"anti-shoplifting/src/middlewares"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	api := app.Group("api")

	// initialize
	initialize := api.Group("initialize")
	initialize.Post("sensor/position", controllers.InitializeSensorPosition)

	// search
	search := api.Group("search")
	search.Get("prefectures", controllers.GetPrefectures)
	search.Post("user/id", controllers.GetUserIdBySensorCd)

	// send message for push notification
	message := api.Group("message")
	message.Post("incidents", controllers.NotifyIncidents)

	user := api.Group("user")
	user.Post("register/company", controllers.RegistCompany)
	user.Get("register/company/verify", controllers.VerifyCompany)
	user.Post("register/store", controllers.RegistStore)
	user.Get("register/store/approve", controllers.ApproveStore)
	user.Get("register/store/verify", controllers.VerifyStore)
	user.Post("register/store/reset-password", controllers.ResetStorePassword)
	user.Post("register/fcm", controllers.RegistFcmToken)
	user.Post("register/fcm/refresh", controllers.RefreshFcmToken)
	user.Post("login", controllers.Login)

	userAuthenticated := user.Use(middlewares.IsAuthenticated)
	userAuthenticated.Post("/", controllers.User) // Get to Post
	userAuthenticated.Put("info", controllers.UpdateInfo)
	userAuthenticated.Put("password", controllers.UpdatePassword)
	userAuthenticated.Post("logout", controllers.Logout)

	// cameras
	userAuthenticated.Post("cameras", controllers.Cameras) // Get to Post
	userAuthenticated.Post("register/cameras", controllers.RegistCameras)
	userAuthenticated.Post("cameras/:serial_no", controllers.GetCamera)      // Get to Post
	userAuthenticated.Delete("cameras/:serial_no", controllers.DeleteCamera) // Get to Post

	// sensors
	userAuthenticated.Post("sensors", controllers.Sensors) // Get to Post
	userAuthenticated.Post("register/sensors", controllers.RegistSensors)
	userAuthenticated.Post("sensors/:sensor_cd", controllers.GetSensor)      // Get to Post
	userAuthenticated.Delete("sensors/:sensor_cd", controllers.DeleteSensor) // Get to Post

	// blacklists
	userAuthenticated.Post("blacklists", controllers.BlackLists) // Get to Post
	userAuthenticated.Post("register/blacklists", controllers.RegistBlackLists)
	userAuthenticated.Post("blacklists/:id", controllers.GetBlackList)      // Get to Post
	userAuthenticated.Delete("blacklists/:id", controllers.DeleteBlackList) // Get to Post

	// incidents
	userAuthenticated.Post("incidents", controllers.Incidents) // Get to Post
	userAuthenticated.Post("register/incidents", controllers.RegistIncidents)
	// userAuthenticated.Get("incidents/:id", controllers.GetIncident) .. replace with grpc
}
