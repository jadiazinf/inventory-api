package health

import (
"github.com/gofiber/fiber/v2"
"github.com/jadiazinf/inventory/internal/platform/i18n"
)

type Handler struct{}

func NewHandler() *Handler {
return &Handler{}
}

func (h *Handler) Check(c *fiber.Ctx) error {
message := i18n.Translate(c, "health_check")
return c.Status(fiber.StatusOK).JSON(fiber.Map{
"status":  "ok",
"message": message,
})
}

func (h *Handler) Greet(c *fiber.Ctx) error {
name := c.Params("name")
message := i18n.TranslateWithData(c, "greeting", map[string]interface{}{
"Name": name,
})
return c.JSON(fiber.Map{
"message": message,
})
}
