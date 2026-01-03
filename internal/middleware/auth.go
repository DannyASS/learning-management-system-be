package middleware

import (
	"strconv"
	"strings"

	"github.com/DannyAss/users/config"
	auth_usercase "github.com/DannyAss/users/internal/usecase/auth"
	"github.com/DannyAss/users/pkg/i18n"
	"github.com/DannyAss/users/pkg/presentation"
	"github.com/DannyAss/users/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(cfg *config.ConfigEnv) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return presentation.Response[any]().
				SetStatus(true).
				SetStatusCode(fiber.StatusUnauthorized).
				SetStatusText(i18n.T(c, "global.http."+strconv.Itoa(fiber.StatusUnauthorized))).
				SetMessage(i18n.T(c, "global.message.9")).
				Json(c)
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return presentation.Response[any]().
				SetStatus(true).
				SetStatusCode(fiber.StatusUnauthorized).
				SetStatusText(i18n.T(c, "global.http."+strconv.Itoa(fiber.StatusUnauthorized))).
				SetMessage(i18n.T(c, "global.message.10")).
				Json(c)
		}

		encToken := strings.TrimSpace(parts[1])

		// Decrypt first
		cs, _ := utils.NewCryptoService([]byte(cfg.AppKey))
		raw, err := cs.Decrypt(encToken)
		if err != nil {
			return presentation.Response[any]().
				SetStatus(true).
				SetStatusCode(fiber.StatusUnauthorized).
				SetStatusText(i18n.T(c, "global.http."+strconv.Itoa(fiber.StatusUnauthorized))).
				SetMessage(i18n.T(c, "global.message.11")).
				Json(c)
		}

		tokenString := string(raw)

		claims, err_ := auth_usercase.VerifyAccessToken([]byte(cfg.JWTSecretKey), tokenString)
		if err_ != nil {
			return presentation.Response[any]().
				SetStatus(true).
				SetStatusCode(fiber.StatusUnauthorized).
				SetStatusText(i18n.T(c, "global.http."+strconv.Itoa(fiber.StatusUnauthorized))).
				SetMessage(i18n.T(c, "global.message.12")).
				Json(c)
		}

		c.Locals("userID", claims.UserID)
		c.Locals("name", claims.Name)
		c.Locals("email", claims.Email)
		c.Locals("roleIDs", claims.RoleIds)

		return c.Next()
	}
}
