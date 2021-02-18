package services

import (
	"context"
	"fmt"
	"log"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/srcabl/gateway/graph/model"
	"github.com/srcabl/gateway/internal/util"
	postspb "github.com/srcabl/protos/posts"
	sharedpb "github.com/srcabl/protos/shared"
	"github.com/srcabl/services/pkg/config"
	"google.golang.org/grpc"
)

// PostsClient defeines the behavior of a posts client
type PostsClient interface {
	Run() (func() error, error)
	CreatePost(context.Context, model.CreatePostRequest) (*model.CommonPostResponse, error)
	Posts(context.Context, model.PostsRequest) (*model.CommonPostsResponse, error)
	CurrentUsersPosts(context.Context) (*model.CommonPostsResponse, error)
}

type postsClient struct {
	postsPort    int
	postsConn    *grpc.ClientConn
	postsService postspb.PostsServiceClient

	sourcesClient SourcesClient
}

// NewPostsClient news up the posts client
func NewPostsClient(config *config.Gateway, sourcesClient SourcesClient) (PostsClient, error) {
	return &postsClient{
		postsPort:     config.Services.PostsPort,
		sourcesClient: sourcesClient,
	}, nil
}

// Run starts up the clients
func (c *postsClient) Run() (func() error, error) {
	log.Printf("Starting Posts Client Connection on: %d\n", c.postsPort)
	postsConn, err := grpc.Dial(fmt.Sprintf("localhost:%d", c.postsPort), grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to dial to posts port: %d", c.postsPort)
	}
	c.postsConn = postsConn
	c.postsService = postspb.NewPostsServiceClient(postsConn)

	return c.close(), nil
}

// Close closes the grpc connection
func (c *postsClient) close() func() error {
	return func() error {
		if err := c.postsConn.Close(); err != nil {
			return errors.Wrap(err, "failed to close posts connection")
		}
		return nil
	}
}

// CreatePost handles creating a post
func (c *postsClient) CreatePost(ctx context.Context, input model.CreatePostRequest) (*model.CommonPostResponse, error) {
	fmt.Println("In create post")
	userUUID := util.GetUserUUIDFromContext(ctx)
	if userUUID == nil {
		fmt.Printf("No user\n")
		return nil, errors.New("user not determined")
	}
	// Check if link exists
	var link *sharedpb.Link
	getLinkByURLReq := model.GetLinkByURLRequest(input)
	linkRes, err := c.postsService.GetLink(ctx, getLinkByURLReq)
	fmt.Printf("GetLinkRes: %+v\n", linkRes)
	fmt.Printf("GetLinkError: %+v\n", err)
	if err == nil {
		link = linkRes.Link
	}
	if link == nil {
		fmt.Println("Creating link")
		// if not, determine
		determineSourceReq := model.CreatePostRequestToPBDetermineSourceRequest(input)
		source, err := c.sourcesClient.Service().DetermineLinkSource(ctx, determineSourceReq)
		if err != nil {
			fmt.Printf("Source Error: %+v\n", err)
			return nil, errors.Wrapf(err, "failed to determine the source of %s", input.URL)
		}
		fmt.Printf("\n\n\n\nDetermined Sources %+v\n", source)
		createLinkReq := model.CreatePostRequestToPBCreateLinkRequest(input, source.PrimarySourceNodes)
		fmt.Printf("CreateLinkReq %+v\n", createLinkReq)
		createLinkRes, err := c.postsService.CreateLink(ctx, createLinkReq)
		fmt.Printf("CreateLinkRes: %+v\n", createLinkRes)
		fmt.Printf("CreateLinkError: %+v\n", err)
		if err != nil {
			fmt.Printf("Create link Error: %+v\n", err)
			return nil, errors.Wrap(err, "failed to create link")
		}
		link = createLinkRes.Link
	}
	createPostReq := model.CreatePostRequestToPBCreatePostRequest(input, userUUID, link.Uuid)
	createPostRes, resErr := c.postsService.CreatePost(ctx, createPostReq)
	fmt.Printf("CreatePostRes: %+v\n", createPostRes)
	fmt.Printf("CreatePostError: %+v\n", resErr)
	postRes := model.PBCreatePostLinkResponseToCommonPostResponse(createPostRes, link, resErr)
	fmt.Printf("Create Post Res: %+v\n", postRes)
	return postRes, nil
}

func (c *postsClient) CurrentUsersPosts(ctx context.Context) (*model.CommonPostsResponse, error) {
	// get the current user uuid
	userUUID := util.GetUserUUIDFromContext(ctx)
	if userUUID == nil {
		return nil, errors.New("no user found")
	}
	fmt.Printf("userUUID: %+v", userUUID)
	res, err := c.getPostsFromUser(ctx, userUUID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get posts for current user %s")
	}
	return res, nil
}

func (c *postsClient) Posts(ctx context.Context, input model.PostsRequest) (*model.CommonPostsResponse, error) {
	userUUID, err := uuid.FromString(input.UserID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert session user id to uuid")
	}
	res, err := c.getPostsFromUser(ctx, userUUID.Bytes())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get posts for user %s", input.UserID)
	}
	return res, nil
}

func (c *postsClient) getPostsFromUser(ctx context.Context, userID []byte) (*model.CommonPostsResponse, error) {
	req := &postspb.ListUsersPostsRequest{UserUuid: userID}
	fmt.Printf("req: %+v", req)
	res, err := c.postsService.ListUsersPosts(ctx, req)
	fmt.Printf("res: %+v", res)
	if err != nil {
		fmt.Printf("res error: %+v", err)
		return nil, errors.Wrap(err, "failed to get the posts for user")
	}
	fmt.Printf("posts: %+v", res.Posts)
	var posts []*model.PartialPost
	for i, p := range res.Posts {
		postUUID, err := uuid.FromBytes(p.Uuid)
		if err != nil {
			return nil, errors.Wrap(err, "failed to convert post id to uuid")
		}
		userUUID, err := uuid.FromBytes(p.UserUuid)
		if err != nil {
			return nil, errors.Wrap(err, "failed to convert user id to uuid")
		}
		post := model.PartialPost{
			ID:      postUUID.String(),
			UserID:  userUUID.String(),
			Comment: p.Comment.PrimaryContent,
			LinkURL: res.Links[i].Url,
		}
		posts = append(posts, &post)
	}
	return &model.CommonPostsResponse{
		Posts: posts,
	}, nil
}
