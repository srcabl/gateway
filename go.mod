module github.com/srcabl/gateway

go 1.15

replace github.com/srcabl/services => /home/kero/automata/srcabl/services

replace github.com/srcabl/protos => /home/kero/automata/srcabl/protos

require (
	github.com/99designs/gqlgen v0.13.0
	github.com/go-chi/chi v3.3.2+incompatible
	github.com/gofrs/uuid v4.0.0+incompatible
	github.com/gorilla/sessions v1.2.1
	github.com/pkg/errors v0.9.1
	github.com/rs/cors v1.6.0
	github.com/smartystreets/assertions v1.0.0 // indirect
	github.com/srcabl/protos v0.1.0
	github.com/srcabl/services v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.5.1 // indirect
	github.com/vektah/gqlparser/v2 v2.1.0
	golang.org/x/crypto v0.0.0-20191011191535-87dc89f01550
	golang.org/x/net v0.0.0-20200324143707-d3edc9973b7e // indirect
	google.golang.org/grpc v1.32.0
	gopkg.in/yaml.v2 v2.2.8 // indirect
)
