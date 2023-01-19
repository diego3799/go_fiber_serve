package routes

import (
	"fun_server/db"
	"fun_server/models"
	"fun_server/utils"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type userRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func UserRoutes(router *fiber.Router) {

	userRoutes := (*router).Group("/users")
	userRoutes.Post("/signup", func(c *fiber.Ctx) error {
		body := userRequest{}
		if err := c.BodyParser(&body); err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}

		password, err := utils.CreateHashPassword(body.Password)
		if err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}

		userDb := models.User{
			Email:    body.Email,
			Password: password,
		}

		createdUser := db.Connection.Create(&userDb)

		if createdUser.Error != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}

		userJwt, err := utils.CreateUserJwt(userDb.Email)

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(map[string]string{
				"message": "There was an error creating the jwt",
			})
		}
		return c.Status(http.StatusCreated).JSON(
			map[string]string{
				"jwt": userJwt,
			},
		)
	})

	userRoutes.Post("/signin", func(c *fiber.Ctx) error {
		body := userRequest{}
		if err := c.BodyParser(&body); err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}
		// See if the user exist
		user := models.User{}
		result := db.Connection.First(&user, models.User{
			Email: body.Email,
		})
		if result.Error != nil {
			if result.Error.Error() == utils.NotFound {
				return c.SendStatus(http.StatusNotFound)
			}
			return c.SendStatus(http.StatusInternalServerError)
		}
		// compare password
		err := utils.ComparePassword(user.Password, body.Password)

		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(map[string]string{
				"message": "Error in email or password",
			})
		}

		userJwt, err := utils.CreateUserJwt(user.Email)

		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(map[string]string{
				"message": "There was an error creating the jwt",
			})
		}
		return c.Status(http.StatusCreated).JSON(
			map[string]string{
				"jwt": userJwt,
			},
		)
	})

}
