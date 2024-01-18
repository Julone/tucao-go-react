package middleware

import "github.com/gofiber/fiber/v2"

var ServerStatus = 0

func Maintenance(ctx *fiber.Ctx) error {
	if ServerStatus == 0 {
		return ctx.Next()
	}
	return ctx.JSON(fiber.Map{"msg": "server is updating now"})
}
