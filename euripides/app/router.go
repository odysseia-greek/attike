package app

import (
	"github.com/gorilla/mux"
	"github.com/graphql-go/handler"
	"github.com/odysseia-greek/attike/euripides/handlers"
	"github.com/odysseia-greek/attike/euripides/middleware"
	"github.com/odysseia-greek/attike/euripides/schemas"
	plato "github.com/odysseia-greek/plato/middleware"
)

// InitRoutes to start up a mux router and return the routes
func InitRoutes() *mux.Router {
	serveMux := mux.NewRouter()

	srv := handler.New(&handler.Config{
		Schema:   &schemas.EuripidesSchema,
		Pretty:   true,
		GraphiQL: false,
	})

	serveMux.HandleFunc("/homeros/v1/health", plato.Adapt(handlers.Health, plato.ValidateRestMethod("GET"), plato.SetCorsHeaders()))
	serveMux.Handle("/graphql", middleware.Adapt(srv, middleware.SetCorsHeaders()))

	return serveMux
}
