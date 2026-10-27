package middleware

import (
	"errors"
	"log/slog"

	"github.com/MrVirul/icecream-inventory/internal/httpx"
	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	var appErr *httpx.AppError
	if errors.As(err, &appErr) {
		resp := fiber.Map{
			"error": fiber.Map{
				"code":    appErr.Code,
				"message": appErr.Message,
			},
		}
		if len(appErr.Details) > 0 {
			resp["error"].(fiber.Map)["details"] = appErr.Details
		}
		return c.Status(appErr.Status).JSON(resp)
	}

	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		code := "BAD_REQUEST"
		if fiberErr.Code == fiber.StatusRequestEntityTooLarge {
			code = "BAD_REQUEST"
		}
		return c.Status(fiberErr.Code).JSON(fiber.Map{
			"error": fiber.Map{
				"code":    code,
				"message": fiberErr.Message,
			},
		})
	}

	reqID := c.Locals("requestid")
	slog.Error("internal error", "error", err, "request_id", reqID, "path", c.Path())
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": fiber.Map{
			"code":    "INTERNAL_ERROR",
			"message": "internal server error",
		},
	})
}
