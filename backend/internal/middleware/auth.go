package middleware

import (
	"strings"
	"time"

	"github.com/MrVirul/icecream-inventory/internal/auth"
	"github.com/MrVirul/icecream-inventory/internal/httpx"
	"github.com/MrVirul/icecream-inventory/internal/models"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

const (
	CtxUserID = "userID"
	CtxRole   = "role"
	CtxJTI    = "jti"
	CtxExp    = "exp"
)

type AuthMiddleware struct {
	jwt *auth.JWTService
	db  *gorm.DB
}

func NewAuth(jwtSvc *auth.JWTService, db *gorm.DB) *AuthMiddleware {
	return &AuthMiddleware{jwt: jwtSvc, db: db}
}

func (m *AuthMiddleware) Authenticate() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return httpx.Unauthorized("missing or invalid token")
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return httpx.Unauthorized("missing or invalid token")
		}
		tokenStr := strings.TrimSpace(parts[1])
		if tokenStr == "" {
			return httpx.Unauthorized("missing or invalid token")
		}

		claims, err := m.jwt.Parse(tokenStr)
		if err != nil {
			return httpx.Unauthorized("missing or invalid token")
		}

		// Check if revoked
		var revoked models.RevokedToken
		if err := m.db.Where("jti = ?", claims.ID).First(&revoked).Error; err == nil {
			return httpx.Unauthorized("missing or invalid token")
		} else if err != gorm.ErrRecordNotFound {
			// treat as unauthorized on error? be conservative
			return httpx.Unauthorized("missing or invalid token")
		}

		// Load user from DB to get current role
		var user models.User
		if err := m.db.First(&user, claims.UserID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return httpx.Unauthorized("missing or invalid token")
			}
			return httpx.Unauthorized("missing or invalid token")
		}

		c.Locals(CtxUserID, user.ID)
		c.Locals(CtxRole, user.Role)
		c.Locals(CtxJTI, claims.ID)
		if claims.ExpiresAt != nil {
			c.Locals(CtxExp, claims.ExpiresAt.Time.Unix())
		} else {
			c.Locals(CtxExp, time.Now().Add(time.Hour).Unix())
		}
		return c.Next()
	}
}
