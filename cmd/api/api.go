package api

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"github.com/yuichi10/edinet_search/db"
	"github.com/yuichi10/edinet_search/graph"
)

const defaultPort = "8080"

//go:embed out/*
var uiEmbedStaticFiles embed.FS

func runServer() {
	// 'out' ディレクトリの内容をサブディレクトリとして取得
	uiStaticFiles, _ := fs.Sub(uiEmbedStaticFiles, "out")
	uiFS := http.FileServer(http.FS(uiStaticFiles))
	nextStaticFiles, _ := fs.Sub(uiEmbedStaticFiles, "out/_next/static")
	nextStaticFS := http.FileServer(http.FS(nextStaticFiles))

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := chi.NewRouter()

	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:8080"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	}).Handler)

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
	srv.AddTransport(&transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	})

	router.Handle("/api", playground.Handler("GraphQL playground", "/query"))
	router.Handle("/api/query", srv)
	router.Handle("/", uiFS)
	router.Handle("/_next/static/*", http.StripPrefix("/_next/static/", nextStaticFS)) // _next/static以下のファイルの配信

	log.Printf("connect to http://0.0.0.0:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func New() *cobra.Command {
	c := &cobra.Command{
		Use:   "api",
		Short: "検索用のAPIを立てます。",
		Run: func(cmd *cobra.Command, args []string) {
			db.OpenDB()
			runServer()
		},
	}
	return c
}
