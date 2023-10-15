package main

import (
	"context"

	"github.com/Collaborators/configs"
	"github.com/Collaborators/handlers"

	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {

	// controllers.CreateStorage()
	Environment := "Development"
	configs.EnvironmentVariables(Environment) // Load .env file in the root directory of project, if it exists

	// Create logger
	loggerUserRegistry := log.New(os.Stdout, "Github Repo contributors Registry", log.LstdFlags|log.Lshortfile)

	// // Create UserHandler
	contributorHandlerV1 := handlers.NewContributorsHandler(loggerUserRegistry)

	// // Create router
	router := mux.NewRouter()

	/// Registered routes ///
	PostUserRegistry := router.Methods("POST").Subrouter()
	PostUserRegistry.HandleFunc("/v1/fetch/contributors", contributorHandlerV1.FetchGithub_ContributorsHandlerV1)

	//// Cron functions ////
	// movie := db.Movie{
	// 	Id:          "avc!23",
	// 	Name:        "Winstonge Bacha",
	// 	Description: "Hey Sexy",
	// }
	// err := db.InitDatabase().CreateMovie(movie)
	// fmt.Println(err)

	// contributorHandlerV1.FetchGithubRepo_Contributors()

	ch := gohandlers.CORS(
		//allowOrigins,
		gohandlers.AllowedMethods([]string{"POST", "GET", "PUT"}),
		gohandlers.AllowedHeaders([]string{"Authorization", "Content-Type", "Access-Control-Allow-Origin", "access-control-allow-origin", "access-control-allow-headers"}),
		gohandlers.AllowedOrigins([]string{"*", "localhost"}),
		gohandlers.AllowCredentials(),
	)
	// routes.Router(ch(router))
	srv := &http.Server{
		Handler: ch(router),
		Addr:    ":8080",
		// Good practice: enforce timeouts for servers you create!
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  55 * time.Second,
		WriteTimeout: 55 * time.Second,
	}

	go func() {
		log.Fatal(srv.ListenAndServe())
	}()

	//This is for gracefully shutdown
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	log.Println("Received request to terminate the server", sig)

	tc, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	srv.Shutdown(tc)
}

/*
	Date: 2023-04-10
	Description: Load environment variables
*/

func LoadEnvironmentVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
}
