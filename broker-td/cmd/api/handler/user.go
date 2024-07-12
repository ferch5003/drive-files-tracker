package handler

import (
	"broker-td/config"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userServiceBaseURL string
	Client             *fiber.Agent
}

func NewUserHandler(configs *config.EnvVars) *UserHandler {
	return &UserHandler{
		userServiceBaseURL: configs.UserServiceBaseURL,
		Client:             fiber.AcquireAgent(),
	}
}

func (h *UserHandler) GetAll(c *fiber.Ctx) error {
	req := h.Client.Request()
	req.Header.SetMethod(fiber.MethodGet)
	req.SetRequestURI(h.userServiceBaseURL + "/users")
	if err := h.Client.Parse(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}

	statusCode, body, errs := h.Client.Bytes()
	if len(errs) > 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": errs,
		})
	}

	var users []fiber.Map
	err := json.Unmarshal(body, &users)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}

	return c.Status(statusCode).JSON(users)
}
