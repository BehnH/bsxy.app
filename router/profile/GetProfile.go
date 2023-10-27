package profile

import (
	"bsxy_app/atproto"
	"bsxy_app/database"
	"bsxy_app/logger"
	"bytes"
	"context"
	"github.com/gofiber/fiber/v2"
	"html/template"
	"strings"
)

func GetProfile(c *fiber.Ctx) error {
	profileRep, err := database.GetCachedUserProfile(c.Params("id"))

	if err != nil && strings.HasPrefix(err.Error(), "cache miss") {
		profileRep, err = atproto.GetActorProfile(context.TODO(), c.Params("id"))
		if err != nil {
			logger.SLog().Panicf("panicking because %w", err)
			panic(err)
		}
		err = nil
	}
	if err != nil {
		logger.SLog().Panicf("panicking because %w", err)
		panic(err)
	}

	var tmplFile = "./router/profile/profile.tmpl"
	tmpl, err := template.New("profile.tmpl").ParseFiles(tmplFile)
	if err != nil {
		panic(err)
	}

	var output bytes.Buffer
	err = tmpl.Execute(&output, map[string]interface{}{
		"displayName": profileRep.DisplayName,
		"handle":      profileRep.Handle,
		"banner":      profileRep.Banner,
		"description": profileRep.Description,
		"followers":   profileRep.FollowersCount,
		"following":   profileRep.FollowsCount,
	})
	if err != nil {
		panic(err)
	}

	err = database.AddProfileToCache(*profileRep)
	if err != nil {
		logger.SLog().Errorf("failed to cache profile: %w", err)
	}

	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
	return c.Status(200).SendString(output.String())
}
