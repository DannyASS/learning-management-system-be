package web_home_handler

import (
	"github.com/DannyAss/users/config"

	"github.com/gofiber/fiber/v2"
)

type HomeHandler struct {
	cfg *config.ConfigEnv
}

func NewHomeHandler(cfg *config.ConfigEnv) *HomeHandler {
	return &HomeHandler{
		cfg: cfg,
	}
}

func (h *HomeHandler) HomePage(c *fiber.Ctx) error {
	return c.Render("web/home", fiber.Map{
		"AppName": h.cfg.AppName,
	}, "layouts/base")
}
