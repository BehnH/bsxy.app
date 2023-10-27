package atproto

import (
	"bsxy_app/logger"
	"context"
	"encoding/json"
	"fmt"
	"github.com/bluesky-social/indigo/api/atproto"
	"github.com/bluesky-social/indigo/api/bsky"
	"github.com/bluesky-social/indigo/xrpc"
	"os"
)

func MakeXRPCClient(ctx context.Context) (*xrpc.Client, error) {
	username := os.Getenv("BSKY_USERNAME")
	pass := os.Getenv("BSKY_PASSWORD")

	xrpcc := &xrpc.Client{
		Host: "https://bsky.social",
		Auth: &xrpc.AuthInfo{Handle: username},
	}

	auth, err := atproto.ServerCreateSession(ctx, xrpcc, &atproto.ServerCreateSession_Input{
		Identifier: xrpcc.Auth.Handle,
		Password:   pass,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	xrpcc.Auth.Did = auth.Did
	xrpcc.Auth.AccessJwt = auth.AccessJwt
	xrpcc.Auth.RefreshJwt = auth.RefreshJwt

	return xrpcc, nil
}

func GetActorProfile(ctx context.Context, handle string) (*bsky.ActorDefs_ProfileViewDetailed, error) {
	xrpcc, err := MakeXRPCClient(ctx)
	if err != nil {
		logger.SLog().Errorf("cannot create client: %w", err)
		return nil, err
	}

	profile, err := bsky.ActorGetProfile(ctx, xrpcc, handle)
	if err != nil {
		logger.SLog().Errorf("cannot get profile: %w", err)
		return nil, err
	}

	json.NewEncoder(os.Stdout).Encode(profile)

	return profile, nil
}

func GetFeedPostAndParent(ctx context.Context, id string) (*bsky.FeedGetPostThread_Output, error) {
	xrpcc, err := MakeXRPCClient(ctx)
	if err != nil {
		logger.SLog().Errorf("cannot create client: %w", err)
		return nil, err
	}

	post, err := bsky.FeedGetPostThread(ctx, xrpcc, 0, 1, id)
	if err != nil {
		logger.SLog().Errorf("cannot get profile: %w", err)
		return nil, err
	}

	json.NewEncoder(os.Stdout).Encode(post)

	return post, nil
}

func GetFeedPost(ctx context.Context, uri string) (*bsky.FeedGetPosts_Output, error) {
	xrpcc, err := MakeXRPCClient(ctx)
	if err != nil {
		logger.SLog().Errorf("cannot create client: %w", err)
		return nil, err
	}

	post, err := bsky.FeedGetPosts(ctx, xrpcc, []string{uri})
	if err != nil {
		logger.SLog().Error("cannot get post: ", err)
		return nil, err
	}

	json.NewEncoder(os.Stdout).Encode(post)
	return post, nil
}
