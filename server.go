package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"

	"github.com/3XBAT/coments-system/graph"
	"github.com/3XBAT/coments-system/pkg/repository"
	"github.com/3XBAT/coments-system/pkg/repository/postgres"
	"github.com/3XBAT/coments-system/pkg/service"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

const defaultPort = "8080"

func main() {
	input := bufio.NewScanner(os.Stdin)

	fmt.Println("Enter type of repository (db/cache)");
	input.Scan()
	repoType := input.Text()
	fmt.Printf("you entered this - %s", repoType)
	
	if err := input.Err(); err != nil {
		fmt.Println("Undefined type", err)
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("Error loading .env file: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	

	var repos *repository.Repository

	switch repoType {
	case "db":
		db, err := postgres.NewPostgresDB(postgres.Config{
			Host:     "localhost",
			Port:     "5436",
			Username: "postgres",
			Password: os.Getenv("db_password"),
			DBName:   "commentDb",
			SSLMode:  "disable",
		})
		if err != nil {
			logrus.Fatalf(fmt.Sprintf("failed to initialized db:%s", err.Error()))
		}
		repos = repository.NewRepositoryPostgres(db)
	case "cache":
		repos = repository.NewRepositoryCache()
	}
	
	service := service.NewService(repos)
	resolver := graph.NewResolver(service)

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	logrus.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	logrus.Fatal(http.ListenAndServe(":"+port, nil))
}
