package boot

import (
	"github.com/pkg/errors"
	"github.com/srcabl/gateway/internal/server"
	"github.com/srcabl/gateway/internal/services"
	"github.com/srcabl/services/pkg/config"
)

// Strap initializes the gateway service
type Strap struct {
	Config        *config.Gateway
	UsersClient   services.UsersClient
	PostsClient   services.PostsClient
	SourcesClient services.SourcesClient
	GraphServer   server.GraphQL

	onconnect  map[string](func() (func() error, error))
	onshutdown map[string](func() error)
}

// New news up a boot strap
func New(cfg *config.Gateway) (*Strap, error) {
	usersClient, err := services.NewUsersClient(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to new up users client")
	}

	sourcesClient, err := services.NewSourcesClient(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to new up sources client")
	}

	postsClient, err := services.NewPostsClient(cfg, sourcesClient)
	if err != nil {
		return nil, errors.Wrap(err, "failed to new up posts client")
	}

	server, err := server.New(cfg, usersClient, postsClient, sourcesClient)
	if err != nil {
		return nil, errors.Wrap(err, "failed to new up the graph ql server")
	}

	return &Strap{
		Config:        cfg,
		UsersClient:   usersClient,
		PostsClient:   postsClient,
		SourcesClient: sourcesClient,
		GraphServer:   server,

		onconnect: map[string](func() (func() error, error)){
			"users client run":   usersClient.Run,
			"posts client run":   postsClient.Run,
			"sources client run": sourcesClient.Run,
			"server run":         server.Run,
		},
		onshutdown: map[string](func() error){},
	}, nil
}

// Connect connects all application services
func (s *Strap) Connect() error {
	for name, connect := range s.onconnect {
		os, err := connect()
		if err != nil {
			return errors.Wrapf(err, "%s failed", name)
		}
		s.onshutdown[name] = os
	}
	return nil
}

// Shutdown shuts down all application srvices
func (s *Strap) Shutdown() []error {
	var errs []error
	for name, shutdown := range s.onshutdown {
		err := shutdown()
		if err != nil {
			errs = append(errs, errors.Wrapf(err, "%s failed", name))
		}
	}
	return errs
}
