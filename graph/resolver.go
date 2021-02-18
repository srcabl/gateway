package graph

import (
	"github.com/srcabl/gateway/graph/generated"
	"github.com/srcabl/gateway/internal/services"
)

//go:generate go run github.com/99designs/gqlgen

// Resolver is the graphql resolver
type Resolver struct {
	usersClient   services.UsersClient
	postsClient   services.PostsClient
	sourcesClient services.SourcesClient
}

// New news up the graphql resolvers
func New(usersClient services.UsersClient, postsClient services.PostsClient, sourcesClient services.SourcesClient) (generated.ResolverRoot, error) {
	return &Resolver{
		usersClient:   usersClient,
		postsClient:   postsClient,
		sourcesClient: sourcesClient,
	}, nil
}
