package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"

	"github.com/redtoad/xcom-editor/lib/api"
)

var root, port string

const defaultPort = "8080"

func main() {

	flag.StringVar(&port, "port", os.Getenv("PORT"), "port for server")
	flag.Parse()

	if root = flag.Arg(0); root == "" {
		root = "."
	}

	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(
		api.NewExecutableSchema(
			api.Config{
				Resolvers: &api.Resolver{
					RootPath: root,
				}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
