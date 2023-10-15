package main

import (
	"AECS_Assignment/configs"
	"AECS_Assignment/handlers"
	"context"

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
	loggerUserRegistry := log.New(os.Stdout, "Github Commit Registry", log.LstdFlags|log.Lshortfile)

	// // Create UserHandler
	commitHandlerV1 := handlers.NewPullsHandler(loggerUserRegistry)

	// // Create router
	router := mux.NewRouter()

	log.Println("Starting Service")

	// // Cron Jobs //
	// c := cron.New(cron.WithSeconds())
	// _, err := c.AddFunc("0 0 * * * *", func() { // This will run every day at midnight
	// 	// Initialize your PullsHandler

	// 	// Fetch previous day's pulls
	// 	commitHandlerV1.FetchGithub_PullsV1()
	// })

	// if err != nil {
	// 	fmt.Printf("Could not initialize cron job: %v\n", err)
	// 	return
	// }

	// c.Start()

	// commitHandlerV1.FetchGithub_PullsV1()

	/// Registered routes ///
	PostUserRegistry := router.Methods("POST").Subrouter()
	PostUserRegistry.HandleFunc("/v1/fetch/pulls", commitHandlerV1.FetchGithub_PullsHandlerV1)

	//// Cron functions ////

	ch := gohandlers.CORS(
		//allowOrigins,
		gohandlers.AllowedMethods([]string{"POST", "GET", "PUT"}),
		gohandlers.AllowedHeaders([]string{"Authorization", "Content-Type", "Access-Control-Allow-Origin", "access-control-allow-origin", "access-control-allow-headers"}),
		gohandlers.AllowedOrigins([]string{"*", "localhost"}),
		gohandlers.AllowCredentials(),
	)
	// routes.Router(ch(router))
	//
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
