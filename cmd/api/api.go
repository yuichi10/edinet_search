package api

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/spf13/cobra"
	"github.com/yuichi10/edinet_search/graph"
)

const defaultPort = "8080"

func runServer() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}



func New() *cobra.Command {
	c := &cobra.Command{
		Use:   "api",
		Short: "検索用のAPIを立てます。",
		Run: func(cmd *cobra.Command, args []string) {
			runServer()
		},
	}
	return c
}
