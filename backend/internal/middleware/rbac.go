package middleware

import (
	"github.com/MrVirul/icecream-inventory/internal/auth"
	"github.com/MrVirul/icecream-inventory/internal/httpx"
	"github.com/gofiber/fiber/v2"
)

func RequirePermission(p auth.Permission) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleVal := c.Locals(CtxRole)
		if roleVal == nil {
			return httpx.Forbidden()
		}
		roleStr, ok := roleVal.(string)
		if !ok || !auth.Can(auth.Role(roleStr), p) {
			return httpx.Forbidden()
		}
		return c.Next()
	}
}
