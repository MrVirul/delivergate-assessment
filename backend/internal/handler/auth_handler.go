package handler

import (
	"time"

	"github.com/MrVirul/icecream-inventory/internal/middleware"
	"github.com/MrVirul/icecream-inventory/internal/service"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	auth service.AuthService
}

func NewAuthHandler(authSvc service.AuthService) *AuthHandler {
	return &AuthHandler{auth: authSvc}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req loginRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "bad request")
	}
	res, err := h.auth.Login(req.Email, req.Password)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": fiber.Map{
			"access_token": res.AccessToken,
			"token_type":   res.TokenType,
			"expires_in":   res.ExpiresIn,
			"user": fiber.Map{
				"id":    res.User.ID,
				"name":  res.User.Name,
				"email": res.User.Email,
				"role":  res.User.Role,
			},
		},
	})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	jtiVal := c.Locals(middleware.CtxJTI)
	expVal := c.Locals(middleware.CtxExp)
	jti, _ := jtiVal.(string)
	var expUnix int64
	switch v := expVal.(type) {
	case int64:
		expUnix = v
	case int:
		expUnix = int64(v)
	case float64:
		expUnix = int64(v)
	default:
		expUnix = time.Now().Add(time.Hour).Unix()
	}
	if err := h.auth.Logout(jti, expUnix); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}
