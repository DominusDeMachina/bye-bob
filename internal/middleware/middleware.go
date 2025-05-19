package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// RequestLogger middleware for logging requests
func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Continue to the next middleware/handler
		err := c.Next()

		// Return any errors that occurred during the request
		return err
	}
}

// Other middleware functions will be added here as needed 