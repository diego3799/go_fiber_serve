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

	(*router).Post("/signup", func(c *fiber.Ctx) error {
		body := userRequest{}
		if err := c.BodyParser(&body); err != nil {
			return utils.ErrorResponse(c, http.StatusInternalServerError, "Unknown error")
		}

		password, err := utils.CreateHashPassword(body.Password)
		if err != nil {
			return utils.ErrorResponse(c, http.StatusInternalServerError, "Unknown error")
		}

		userDb := models.User{
			Email:    body.Email,
			Password: password,
		}

		createdUser := db.Connection.Create(&userDb)

		if createdUser.Error != nil {
			return utils.ErrorResponse(c, http.StatusInternalServerError, "Error creating user")
		}

		userJwt, err := utils.CreateUserJwt(userDb.Email)

		if err != nil {
			return utils.ErrorResponse(c, http.StatusInternalServerError, "Error creating jwt")
		}
		return c.Status(http.StatusCreated).JSON(
			map[string]string{
				"jwt": userJwt,
			},
		)
	})

	(*router).Post("/signin", func(c *fiber.Ctx) error {
		body := userRequest{}
		if err := c.BodyParser(&body); err != nil {
			return utils.ErrorResponse(c, http.StatusInternalServerError, "Error in email or password")
		}
		// See if the user exist
		user := models.User{}
		result := db.Connection.First(&user, models.User{
			Email: body.Email,
		})
		if result.Error != nil {
			if result.Error.Error() == utils.NotFound {
				return utils.ErrorResponse(c, http.StatusNotFound, "user not found")
			}
			return utils.ErrorResponse(c, http.StatusBadRequest, "Error in email or password")
		}
		// compare password
		err := utils.ComparePassword(user.Password, body.Password)

		if err != nil {
			return utils.ErrorResponse(c, http.StatusBadRequest, "Error in email or password")
		}

		userJwt, err := utils.CreateUserJwt(user.Email)

		if err != nil {
			return utils.ErrorResponse(c, http.StatusInternalServerError, "Error creating jwt")
		}
		return c.Status(http.StatusCreated).JSON(
			map[string]string{
				"jwt": userJwt,
			},
		)
	})

}
