package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/joho/godotenv"
	"github.com/rezaif79-ri/graphqltutor1/graph"
	"github.com/rezaif79-ri/graphqltutor1/graph/config/bun"
)

const defaultPort = "8081"

func main() {
	initEnv()

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
	bun.OpenBunDBConn()
	defer bun.CloseBunDBConn()

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe("127.0.0.1:"+port, nil))
}

func initEnv() {
	if os.Getenv("ENV") != "PRODUCTION" {
		godotenv.Load("dev.env")
	}
}
