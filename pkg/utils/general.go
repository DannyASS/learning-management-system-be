package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/DannyAss/users/pkg/i18n"
	"github.com/DannyAss/users/pkg/presentation"

	"github.com/gofiber/fiber/v2"
)

func CleanText(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func RandomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"

	result := make([]byte, length)
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}

	return string(result), nil
}

func TryCatch(c *fiber.Ctx, fn func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = presentation.Response[any]().
				SetStatus(false).
				SetStatusCode(fiber.StatusInternalServerError).
				SetStatusText(i18n.T(c, "global.http.500")).
				SetMessage(i18n.T(c, "global.message.1")). // server error msg
				SetErrorDetail(fmt.Sprint(r)).
				Json(c)
		}
	}()
	return fn()
}

func FilterSlice[T any](data []T, cond func(T) bool) []T {
	out := make([]T, 0)
	for _, v := range data {
		if cond(v) {
			out = append(out, v)
		}
	}
	return out
}

func UIntPtr(v int) *uint {
	u := uint(v)
	return &u
}

func CookieConfig(env, domain string, value string) fiber.Cookie {
	isProd := env == "production"

	if value == "" {
		return fiber.Cookie{
			Name:     "refresh_token",
			Value:    value,
			HTTPOnly: true,
			Secure:   isProd, // local=false, prod=true
			SameSite: func() string {
				return "Lax" // localhost aman pake Lax
			}(),
			Path:    "/",
			Domain:  domain,
			Expires: time.Now().Add(-time.Hour),
			MaxAge:  -1,
		}
	}

	return fiber.Cookie{
		Name:     "refresh_token",
		Value:    value,
		HTTPOnly: true,
		Secure:   isProd, // local=false, prod=true
		SameSite: func() string {
			return "Lax" // localhost aman pake Lax
		}(),
		Path:   "/",
		Domain: domain,
	}
}
