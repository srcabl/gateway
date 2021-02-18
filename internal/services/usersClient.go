package services

import (
	"context"
	"fmt"
	"log"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/srcabl/gateway/graph/model"
	"github.com/srcabl/gateway/internal/util"
	userspb "github.com/srcabl/protos/users"
	"github.com/srcabl/services/pkg/config"
	"google.golang.org/grpc"
)

// UsersClient defines the behavior of a users client
type UsersClient interface {
	Run() (func() error, error)
	//grapql handlers
	CurrentUser(context.Context) (*model.CommonUserResponse, error)
	ChangePassword(context.Context, model.ChangePasswordRequest) (*model.CommonUserResponse, error)
	ForgotPassword(context.Context, string) (bool, error)
	Register(context.Context, model.RegisterUserRequest) (*model.CommonUserResponse, error)
	Login(context.Context, model.LoginUserRequest) (*model.CommonUserResponse, error)
	Logout(context.Context) (bool, error)
	FollowUser(context.Context, model.FollowRequest) (bool, error)
	UnfollowUser(context.Context, model.FollowRequest) (bool, error)
	FollowSource(context.Context, model.FollowRequest) (bool, error)
	UnfollowSource(context.Context, model.FollowRequest) (bool, error)
}

type usersClient struct {
	usersPort   int
	usersConn   *grpc.ClientConn
	usersClient userspb.UsersServiceClient
}

// NewUsersClient news up the users client
func NewUsersClient(config *config.Gateway) (UsersClient, error) {
	return &usersClient{
		usersPort: config.Services.UsersPort,
	}, nil
}

// Run starts up the clients
func (c *usersClient) Run() (func() error, error) {
	log.Printf("Starting Users Client Connection on: %d\n", c.usersPort)
	usersConn, err := grpc.Dial(fmt.Sprintf("localhost:%d", c.usersPort), grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to dial to users port: %d", c.usersPort)
	}
	c.usersConn = usersConn
	c.usersClient = userspb.NewUsersServiceClient(usersConn)

	return c.close(), nil
}

// Close closes the grpc connection
func (c *usersClient) close() func() error {
	return func() error {
		if err := c.usersConn.Close(); err != nil {
			return errors.Wrap(err, "failed to close users connection")
		}
		return nil
	}
}

// ChangePassword handles change password requests
func (c *usersClient) ChangePassword(ctx context.Context, input model.ChangePasswordRequest) (*model.CommonUserResponse, error) {
	//TODO
	field := ""
	message := "not implemented"
	return &model.CommonUserResponse{
		Errors: []*model.Error{
			{
				Field:   &field,
				Message: &message,
			},
		},
	}, nil
}

// ForgotPassword handles forgot passowrd requests
func (c *usersClient) ForgotPassword(ctx context.Context, email string) (bool, error) {
	//TODO
	return false, nil
}

// Register handles user register requests
func (c *usersClient) Register(ctx context.Context, input model.RegisterUserRequest) (*model.CommonUserResponse, error) {
	userReq, commonErr := model.ValidateRegisterAndRegisterUserRequestToPBCreateUserRequest(input)
	if commonErr != nil {
		return &model.CommonUserResponse{
			Errors: commonErr,
		}, nil
	}
	createRes, err := c.usersClient.CreateUser(ctx, userReq)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create user")
	}
	// set the user id for the session
	util.SetUserUUIDToContext(ctx, createRes.User.Uuid)
	//get the user
	userRes := model.PBCreateUserResponseToCommonUserResponse(createRes, err)
	return userRes, nil
}

//Login handles login requests
func (c *usersClient) Login(ctx context.Context, input model.LoginUserRequest) (*model.CommonUserResponse, error) {
	userReq, err := model.LoginUserRequestToPBValidateUserCredentials(input)
	if err != nil {
		return nil, fmt.Errorf("failed to transform graph request to grpc request: %+v", err)
	}
	fmt.Printf("UserReq: %+v", userReq)
	res, resErr := c.usersClient.ValidateUserCredentials(ctx, userReq)
	if res != nil && res.User != nil {
		util.SetUserUUIDToContext(ctx, res.User.Uuid)
	}
	userRes := model.PBValidateUserResponseToCommonUserResponse(res, resErr)
	fmt.Printf("UserRes: %+v", userRes)
	return userRes, nil
}

//Logout handles logout requests
func (c *usersClient) Logout(ctx context.Context) (bool, error) {
	util.SetUserUUIDToContext(ctx, nil)
	return true, nil
}

//CurrentUser handles current user requests
func (c *usersClient) CurrentUser(ctx context.Context) (*model.CommonUserResponse, error) {
	userUUID := util.GetUserUUIDFromContext(ctx)
	if userUUID == nil {
		return nil, nil
	}
	userReq := model.CurrentUserRequestToPBGetUserRequest(userUUID)
	res, resErr := c.usersClient.GetUser(ctx, userReq)
	userRes := model.PBGetUserResponseToCommonUserResponse(res, resErr)
	return userRes, nil
}

func (c *usersClient) FollowUser(ctx context.Context, input model.FollowRequest) (bool, error) {
	return c.performFollowReq(ctx, input, c.usersClient.Follow, userspb.FollowRequest_USER)
}

func (c *usersClient) UnfollowUser(ctx context.Context, input model.FollowRequest) (bool, error) {
	return c.performFollowReq(ctx, input, c.usersClient.UnFollow, userspb.FollowRequest_USER)
}

func (c *usersClient) FollowSource(ctx context.Context, input model.FollowRequest) (bool, error) {
	return c.performFollowReq(ctx, input, c.usersClient.Follow, userspb.FollowRequest_SOURCE)
}

func (c *usersClient) UnfollowSource(ctx context.Context, input model.FollowRequest) (bool, error) {
	return c.performFollowReq(ctx, input, c.usersClient.UnFollow, userspb.FollowRequest_SOURCE)
}

type pbFollowFunc func(ctx context.Context, in *userspb.FollowRequest, opts ...grpc.CallOption) (*userspb.FollowResponse, error)

func (c *usersClient) performFollowReq(ctx context.Context, input model.FollowRequest, followFunc pbFollowFunc, followType userspb.FollowRequest_FollowedType) (bool, error) {
	followerUserUUID := util.GetUserUUIDFromContext(ctx)
	if followerUserUUID == nil {
		return false, errors.New("current user does not exist")
	}
	followedUserUUID, err := uuid.FromString(input.FollowedID)
	if err != nil {
		return false, errors.Wrap(err, "followed user uuid is not valid")
	}
	followReq := model.ToPBFollowRequest(followerUserUUID, followedUserUUID.Bytes(), followType)
	_, err = followFunc(ctx, followReq)
	if err != nil {
		return false, errors.Wrap(err, "failed to do follow request")
	}
	return true, nil
}
