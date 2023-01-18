package routes

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type userRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func UserRoutes(router *fiber.Router) {

	userRoutes := (*router).Group("/users")
	userRoutes.Post("/signin", func(c *fiber.Ctx) error {
		body := userRequest{}
		if err := c.BodyParser(&body); err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}

		fmt.Printf("%v", body)

		return c.SendStatus(http.StatusAccepted)
	})

	userRoutes.Post("/singup", func(c *fiber.Ctx) error {
		body := userRequest{}
		if err := c.BodyParser(&body); err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}

		fmt.Printf("%v", body)

		return c.SendStatus(http.StatusAccepted)
	})

}
