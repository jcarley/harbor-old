package main

import (
	"log"
	"net/http"
	"os"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/jcarley/harbor/controllers"
	"github.com/jcarley/harbor/web"
)

var (
	router *negroni.Negroni
)

func NewServer() *http.Server {
	wd, _ := os.Getwd()
	log.Println("Working directory", wd)

	router := NewRouter()

	n := negroni.Classic()
	n.UseHandler(router)

	return &http.Server{
		Addr:    ":3000",
		Handler: n,
	}
}

func StartServer(server *http.Server) {
	log.Println("Starting server")

	// err := server.ListenAndServeTLS("ssl/cert.pem", "ssl/key.pem")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalln("Error: %v", err)
	}

}

func NewRouter() http.Handler {

	// Setup session store
	store := web.NewCookieStore()
	appContext := web.NewContext(store)

	router := mux.NewRouter()

	// Add handlers for routes
	authenticationController := controllers.NewAuthenticationController(appContext)
	authenticationController.Register(router)

	return router
}
