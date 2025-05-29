package main

//Notes
//Go mod vendor: It vendors your dependencies — i.e., makes local copies — so your project doesn't need to fetch them from the internet.
//A dependency is any external package or module that your code relies on to work.

//For example, if you write a Go web server using the gorilla/mux router, then gorilla/mux is a dependency — you’re depending on it to provide functionality.

//direct dependency : A package you explicitly import in your code.

//indirect dependency : A package that is required by one of your direct dependencies, not by you directly.

//go mod tidy :
//Adds missing dependencies used in your code but not listed in go.mod.

// Removes unused dependencies listed in go.mod and go.sum

//Terminal represents a session ; if u press the dustbin button -- it kills the current session

//origin: scheme + host + port
//cross-origin resource sharing --> one script running on one origin may access resources on another origin; eg. Js code on one origin tries to access api on another origin
import (
	
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/teadrinke/Go/internal/database"
	

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	fmt.Println("Hello world!")

	godotenv.Load()

	portString := os.Getenv("PORT")
	if portString == ""{
		log.Fatal("PORT is not found in the environment")
	}

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == ""{
		log.Fatal("DB_URL is not found in the environment")
	}

	// sql.Open("postgres", dbUrl)

	conn, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Can't connect to database:", err)
	}

	queries := database.New(conn)
	

	apiCfg := apiConfig{
		DB : queries,
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{"Link"},
		AllowCredentials: false,
		MaxAge: 300,
	}))

	v1Router := chi.NewRouter()

	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handleErr)
	v1Router.Post("/users", apiCfg.handlerCreateUser)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))
	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	v1Router.Get("/feeds", apiCfg.handlerGetFeeds)
	v1Router.Post("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))
	v1Router.Get("/feed_follows", apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollow))
	router.Mount("/v1", v1Router)
	// v1Router.Use(middleware.Logger)


	srv := &http.Server{
		Handler: router,
		Addr : ":" + portString,
	}
	log.Println("Try opening: http://localhost:" + portString + "/v1/healthz")
	log.Printf("Server starting on port %v", portString)
	

	err = srv.ListenAndServe()
	if err!=nil {
		log.Fatal(err)
	}
	fmt.Println("Port:", portString)
}