package router

import (
	"bsxy_app/logger"
	profile2 "bsxy_app/router/profile"
	"github.com/gofiber/fiber/v2"
)

func Router(app *fiber.App) {
	root := app.Group("/")
	root.Get("/", func(ctx *fiber.Ctx) error {
		logger.SLog().Info("redirecting to bluesky")
		return ctx.Redirect("https://bsky.app", fiber.StatusPermanentRedirect)
	})

	profile := app.Group("/profile")
	profile.Get("/:id", profile2.GetProfile)
	profile.Get("/:id/post/:post", profile2.GetFeedPost)
}
