package server

import (
	"fmt"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
	"github.com/srcabl/gateway/graph"
	"github.com/srcabl/gateway/graph/generated"
	"github.com/srcabl/gateway/internal/middleware"
	"github.com/srcabl/gateway/internal/services"
	"github.com/srcabl/services/pkg/config"
)

// GraphQL defines the behavior of the graphql server
type GraphQL interface {
	Run() (func() error, error)
}

// GraphQLServer is the graphql server
type GraphQLServer struct {
	address    string
	port       int
	sessionkey string

	server *handler.Server
}

// New news up a graphql server
func New(cfg *config.Gateway, usersClient services.UsersClient, postsClient services.PostsClient, sourceClient services.SourcesClient) (GraphQL, error) {
	resolver, err := graph.New(usersClient, postsClient, sourceClient)
	if err != nil {
		return nil, errors.Wrap(err, "failed to new the graphql resolver")
	}
	config := generated.Config{Resolvers: resolver}
	schema := generated.NewExecutableSchema(config)
	srv := handler.NewDefaultServer(schema)

	return &GraphQLServer{
		address:    cfg.Server.Address,
		port:       cfg.Server.Port,
		sessionkey: cfg.Server.SessionKey,
		server:     srv,
	}, nil
}

// Run starts up the server
func (g GraphQLServer) Run() (func() error, error) {
	//initialize session store
	store := sessions.NewCookieStore([]byte(g.sessionkey))

	//create router to inject middleware
	router := chi.NewRouter()
	router.Use(middleware.InjectSession(store))
	router.Use(middleware.InjectCors())

	//set up graphql endpoints
	router.Handle("/graphql", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/query", g.server)

	fullAddr := fmt.Sprintf("%s:%d", g.address, g.port)
	fmt.Printf("Listening on %s\n", fullAddr)
	err := http.ListenAndServe(fullAddr, router)
	return func() error {
		return errors.Wrap(err, "server ended")
	}, nil
}
