package service

import (
	"log"

	"github.com/codegangsta/negroni"
	r "github.com/dancannon/gorethink"
)

var (
	router  *negroni.Negroni
	session *r.Session
)

func init() {
	var err error

	session, err = r.Connect(r.ConnectOpts{
		Address:  "localhost:28015",
		Database: "chat",
		MaxOpen:  40,
	})
	if err != nil {
		log.Fatalln(err.Error())
	}

	r.DbCreate("chat").RunWrite(session)
	r.Db("chat").TableCreate("messages").Run(session)
}
