package main

import (
	"net/http"

	log "github.com/Innovatube/log4go"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
)

func main() {
	log.AddFilter("stdout", log.DEBUG, log.NewConsoleLogWriter())
	goji.Abandon(middleware.Logger)
	goji.Use(log.NewGojiLog("my-app-name"))
	goji.Get("/ping", yourHandler)
	goji.Serve()
}

func yourHandler(c web.C, w http.ResponseWriter, r *http.Request) {
}
