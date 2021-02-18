package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/srcabl/gateway/graph/generated"
	"github.com/srcabl/gateway/graph/model"
)

func (r *mutationResolver) ChangePassword(ctx context.Context, input model.ChangePasswordRequest) (*model.CommonUserResponse, error) {
	return r.usersClient.ChangePassword(ctx, input)
}

func (r *mutationResolver) ForgotPassword(ctx context.Context, email string) (bool, error) {
	return r.usersClient.ForgotPassword(ctx, email)
}

func (r *mutationResolver) Register(ctx context.Context, input model.RegisterUserRequest) (*model.CommonUserResponse, error) {
	return r.usersClient.Register(ctx, input)
}

func (r *mutationResolver) Login(ctx context.Context, input model.LoginUserRequest) (*model.CommonUserResponse, error) {
	return r.usersClient.Login(ctx, input)
}

func (r *mutationResolver) Logout(ctx context.Context) (bool, error) {
	return r.usersClient.Logout(ctx)
}

func (r *mutationResolver) FollowUser(ctx context.Context, input model.FollowRequest) (bool, error) {
	return r.usersClient.FollowUser(ctx, input)
}

func (r *mutationResolver) UnfollowUser(ctx context.Context, input model.FollowRequest) (bool, error) {
	return r.usersClient.UnfollowUser(ctx, input)
}

func (r *mutationResolver) FollowSource(ctx context.Context, input model.FollowRequest) (bool, error) {
	return r.usersClient.FollowSource(ctx, input)
}

func (r *mutationResolver) UnfollowSource(ctx context.Context, input model.FollowRequest) (bool, error) {
	return r.usersClient.UnfollowSource(ctx, input)
}

func (r *mutationResolver) CreatePost(ctx context.Context, input model.CreatePostRequest) (*model.CommonPostResponse, error) {
	return r.postsClient.CreatePost(ctx, input)
}

func (r *mutationResolver) UpdatePost(ctx context.Context, input model.UpdatePostRequest) (*model.CommonPostResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeletePost(ctx context.Context, input model.DeletePostRequest) (*model.CommonPostResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) CurrentUser(ctx context.Context) (*model.CommonUserResponse, error) {
	return r.usersClient.CurrentUser(ctx)
}

func (r *queryResolver) CurrentUserUsersFollowed(ctx context.Context) (*model.CommonUserResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) CurrentUserSourcesFollowed(ctx context.Context) (*model.CommonSourceResponse, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) CurrentUsersPosts(ctx context.Context) (*model.CommonPostsResponse, error) {
	return r.postsClient.CurrentUsersPosts(ctx)
}

func (r *queryResolver) Posts(ctx context.Context, input model.PostsRequest) (*model.CommonPostsResponse, error) {
	return r.postsClient.Posts(ctx, input)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
