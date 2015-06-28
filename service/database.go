package service

import (
	"log"

	r "github.com/dancannon/gorethink"
)

var (
	session *r.Session
)

func init() {
	var err error

	session, err = r.Connect(r.ConnectOpts{
		Address:  "localhost:28015",
		Database: "harbor",
		MaxIdle:  40,
		MaxOpen:  40,
	})
	if err != nil {
		log.Fatalln(err.Error())
	}

	r.DbCreate("harbor").RunWrite(session)
	r.Db("harbor").TableCreate("users").Run(session)
}

func Session() *r.Session {
	return session
}

func Users() r.Term {
	return r.Db("harbor").Table("users")
}
