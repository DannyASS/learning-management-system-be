package middleware

import (
	"strconv"
	"time"

	"github.com/DannyAss/users/config"

	user_model "github.com/DannyAss/users/internal/models/database_model/user"
	users_repository "github.com/DannyAss/users/internal/repositories/users"
	"github.com/DannyAss/users/pkg/i18n"
	"github.com/DannyAss/users/pkg/presentation"

	"github.com/gofiber/fiber/v2"
	"github.com/patrickmn/go-cache"
)

type PermissionMiddleware struct {
	userRepo users_repository.IuserRepos
	cache    *cache.Cache
	cfg      *config.ConfigEnv
}

func NewPermissionMiddleware(roleRepo users_repository.IuserRepos, cfg *config.ConfigEnv) *PermissionMiddleware {
	return &PermissionMiddleware{
		userRepo: roleRepo,
		cfg:      cfg,
		cache:    cache.New(time.Duration(cfg.CacheTTLExpiry), time.Duration(cfg.CachePeriodExpiry)),
	}
}

func (m *PermissionMiddleware) RequirePermissions(permissions []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := c.Locals("userID").(uint)
		if !ok {
			return presentation.Response[any]().
				SetStatus(false).
				SetStatusCode(fiber.StatusUnauthorized).
				SetStatusText(i18n.T(c, "global.http."+strconv.Itoa(fiber.StatusUnauthorized))).
				SetMessage(i18n.T(c, "global.message.13")).
				Json(c)
		}

		roleIDs, ok := c.Locals("roleIDs").(uint64)
		if !ok || roleIDs == 0 {
			return presentation.Response[any]().
				SetStatus(false).
				SetStatusCode(fiber.StatusForbidden).
				SetStatusText(i18n.T(c, "global.http."+strconv.Itoa(fiber.StatusForbidden))).
				SetMessage(i18n.T(c, "global.message.14")).
				Json(c)
		}

		cacheKey := "user_permissions:" + strconv.FormatUint(uint64(userID), 10)

		var userPermissions []string

		if cached, found := m.cache.Get(cacheKey); found {
			if perms, ok := cached.([]string); ok {
				userPermissions = perms
			}
		}

		// If not in cache, fetch from repository
		if userPermissions == nil {
			var err error
			roles := user_model.RoleFilter{
				Id: roleIDs,
			}
			role, err := m.userRepo.GetRole(roles, "")

			roleString := role.Name
			permissions = append(permissions, roleString)
			if err != nil {
				return presentation.Response[any]().
					SetStatus(false).
					SetStatusCode(fiber.StatusForbidden).
					SetStatusText(i18n.T(c, "global.http."+strconv.Itoa(fiber.StatusForbidden))).
					SetMessage(i18n.T(c, "global.message.15")).
					Json(c)
			}

			// Store in cache
			m.cache.Set(cacheKey, userPermissions, time.Duration(m.cfg.CacheTTLExpiry))

		}

		return m.checkPermissions(c, userPermissions, permissions)
	}
}

func (m *PermissionMiddleware) checkPermissions(c *fiber.Ctx, userPermissions []string, requiredPermissions []string) error {
	userPermsMap := make(map[string]bool)
	for _, perm := range userPermissions {
		userPermsMap[perm] = true
	}

	for _, requiredPerm := range requiredPermissions {
		if !userPermsMap[requiredPerm] {
			return presentation.Response[any]().
				SetStatus(false).
				SetStatusCode(fiber.StatusForbidden).
				SetStatusText(i18n.T(c, "global.http."+strconv.Itoa(fiber.StatusForbidden))).
				SetMessage(i18n.T(c, "global.message.16")).
				Json(c)
		}
	}

	return c.Next()
}

func (m *PermissionMiddleware) ClearUserPermissionsCache(userID uint) {
	cacheKey := "user_permissions:" + strconv.FormatUint(uint64(userID), 10)
	m.cache.Delete(cacheKey)
}
