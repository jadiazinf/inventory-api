package i18n

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var bundle *i18n.Bundle

func InitI18n() {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	// Load translation files
	// In a real deployment, ensure these files are available in the working directory
	bundle.LoadMessageFile("locales/active.en.json")
	bundle.LoadMessageFile("locales/active.es.json")
}

func Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get language from Accept-Language header or query param "lang"
		lang := c.Query("lang")
		acceptLang := c.Get("Accept-Language")

		localizer := i18n.NewLocalizer(bundle, lang, acceptLang)

		// Store localizer in context for handlers to use
		c.Locals("localizer", localizer)

		return c.Next()
	}
}

func Translate(c *fiber.Ctx, messageID string) string {
	localizer, ok := c.Locals("localizer").(*i18n.Localizer)
	if !ok {
		return messageID
	}

	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: messageID,
	})
	if err != nil {
		return messageID
	}
	return msg
}

// TranslateWithData allows passing template data
func TranslateWithData(c *fiber.Ctx, messageID string, data map[string]interface{}) string {
	localizer, ok := c.Locals("localizer").(*i18n.Localizer)
	if !ok {
		return messageID
	}

	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: data,
	})
	if err != nil {
		return messageID
	}
	return msg
}
