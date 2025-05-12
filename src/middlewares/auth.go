package middlewares

import (
	"fmt"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func IsAuthenticated(c *fiber.Ctx) error {
	fmt.Printf("IsAuthenticated called")

	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	if data["from"] == "web" {
		// case1. accessed from web app
		cookie := c.Cookies("jwt")

		token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})

		if err != nil || !token.Valid {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{
				"message": "unauthenticated",
			})
		}
	} else if data["from"] == "mobile" {
		fmt.Println("IsAuthenticated from mobile in")

		// case2. accessed from mobile app
		jwtString := data["jwt"]
		fmt.Println("IsAuthenticated get jwt completed")
		token, err := jwt.Parse(jwtString, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte("secret"), nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			fmt.Println("IsAuthenticated token valid !!")
			// fmt.Printf("user_id: %v\n", int64(claims["user_id"].(int64)))
			// fmt.Printf("exp: %v\n", int64(claims["exp"].(int64)))
			fmt.Printf("user_id: %v\n", claims["user_id"].(string))
			fmt.Printf("exp: %v\n", int64(claims["exp"].(float64)))
		} else {
			fmt.Println(err)
		}

		fmt.Println("IsAuthenticated from mobile out")

	} else {
		// case3. accessed from others
		return c.JSON(fiber.Map{
			"message": "Invalid Credentials!",
		})
	}

	return c.Next()
}

func GetUserId(c *fiber.Ctx) (uint, error) {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return 0, err
	}

	var token *jwt.Token

	if data["from"] == "web" {
		// case1. accessed from web app
		cookie := c.Cookies("jwt")

		tokenString, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})
		token = tokenString
		if err != nil {
			return 0, err
		}

		payload := token.Claims.(*jwt.StandardClaims)

		id, _ := strconv.Atoi(payload.Subject)

		return uint(id), nil

	} else if data["from"] == "mobile" {
		// case2. accessed from mobile app
		jwtString := data["jwt"]
		tokenString, err := jwt.Parse(jwtString, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte("secret"), nil
		})
		token = tokenString
		if err != nil {
			return 0, err
		}

		claims, _ := token.Claims.(jwt.MapClaims)
		// fmt.Printf("user_id: %v\n", int64(claims["user_id"].(float64)))
		// fmt.Printf("exp: %v\n", int64(claims["exp"].(float64)))
		id, _ := strconv.ParseUint(claims["user_id"].(string), 10, 64)
		return uint(id), nil

	} else {
		// case3. accessed from others
		return 0, nil
	}
}
