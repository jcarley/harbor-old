package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

var (
	router *negroni.Negroni
)

type appHandler struct {
	*ctx
	handler func(context, http.ResponseWriter, *http.Request)
}

func (this appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	this.handler(this.ctx, w, r)
}

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
	authKey := securecookie.GenerateRandomKey(64)
	encryptionKey := securecookie.GenerateRandomKey(32)
	store := sessions.NewCookieStore(
		authKey,
		encryptionKey,
	)

	appContext := &ctx{
		sessionStore: store,
	}

	router := mux.NewRouter()

	// Add handlers for routes
	router.Handle("/login", appHandler{appContext, login}).Methods("POST")
	router.Handle("/register", appHandler{appContext, register}).Methods("POST")

	return router
}
