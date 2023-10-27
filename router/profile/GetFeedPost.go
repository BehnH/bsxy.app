package profile

import (
	"bsxy_app/atproto"
	"bsxy_app/database"
	"bsxy_app/logger"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/gofiber/fiber/v2"
	"html/template"
	"strings"
)

func GetFeedPost(c *fiber.Ctx) error {
	user := c.Params("id")
	postId := c.Params("post")
	profileRep, err := database.GetCachedUserProfile(user)

	if err != nil && strings.HasPrefix(err.Error(), "cache miss") {
		profileRep, err = atproto.GetActorProfile(context.TODO(), user)
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

	uri := fmt.Sprintf("at://%s/app.bsky.feed.post/%s", profileRep.Did, postId)
	postRep, err := atproto.GetFeedPost(context.TODO(), uri)

	if err != nil {
		logger.SLog().Panic("panicking because ", err)
		panic(err)
	}

	var tmplFile = "./router/profile/post.tmpl"
	tmpl, err := template.New("post.tmpl").ParseFiles(tmplFile)
	if err != nil {
		panic(err)
	}

	var output bytes.Buffer
	rec := postRep.Posts[0].Record.Val.(*bsky.FeedPost)
	tmplData := map[string]interface{}{
		"displayName": postRep.Posts[0].Author.DisplayName,
		"handle":      postRep.Posts[0].Author.Handle,
		"description": rec.Text,
		"likeCount":   postRep.Posts[0].LikeCount,
		"replyCount":  postRep.Posts[0].ReplyCount,
		"repostCount": postRep.Posts[0].RepostCount,
		"post":        postId,
	}

	data, _ := json.Marshal(&postRep.Posts[0].Embed)
	if string(data) != "null" {
		tmplData["imgKind"] = "post_image"
		tmplData["image"] = postRep.Posts[0].Embed.EmbedImages_View.Images[0].Thumb
	} else {
		tmplData["imgKind"] = "user_image"
		tmplData["image"] = postRep.Posts[0].Author.Avatar
	}

	err = tmpl.Execute(&output, tmplData)
	if err != nil {
		panic(err)
	}

	c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
	return c.Status(200).SendString(output.String())
}
