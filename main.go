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
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Hello world!")

	godotenv.Load()

	portString := os.Getenv("PORT")
	if portString == ""{
		log.Fatal("PORT is not found in the environment")
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

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr : ":" + portString,
	}

	log.Printf("Server starting on port %v", portString)
	err := srv.ListenAndServe()
	if err!=nil {
		log.Fatal(err)
	}
	fmt.Println("Port:", portString)
}