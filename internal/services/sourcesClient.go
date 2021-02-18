package services

import (
	"fmt"
	"log"

	"github.com/pkg/errors"
	"github.com/srcabl/protos/sources"
	sourcespb "github.com/srcabl/protos/sources"
	"github.com/srcabl/services/pkg/config"
	"google.golang.org/grpc"
)

// SourcesClient defines the behavior of a sources client
type SourcesClient interface {
	Run() (func() error, error)
	Service() sourcespb.SourcesServiceClient
}

type sourcesClient struct {
	sourcesPort    int
	sourcesConn    *grpc.ClientConn
	sourcesService sourcespb.SourcesServiceClient
}

// NewSourcesClient news up the sources client
func NewSourcesClient(config *config.Gateway) (SourcesClient, error) {
	return &sourcesClient{
		sourcesPort: config.Services.SourcesPort,
	}, nil
}

// Run starts up the clients
func (c *sourcesClient) Run() (func() error, error) {
	log.Printf("Starting Sources Client Connection on: %d\n", c.sourcesPort)
	sourcesConn, err := grpc.Dial(fmt.Sprintf("localhost:%d", c.sourcesPort), grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to dial to sources port: %d", c.sourcesPort)
	}
	c.sourcesConn = sourcesConn
	c.sourcesService = sourcespb.NewSourcesServiceClient(sourcesConn)

	return c.close(), nil
}

// Close closes the grpc connection
func (c *sourcesClient) close() func() error {
	return func() error {
		if err := c.sourcesConn.Close(); err != nil {
			return errors.Wrap(err, "failed to close posts connection")
		}
		return nil
	}
}

func (c *sourcesClient) Service() sources.SourcesServiceClient {
	return c.sourcesService
}
