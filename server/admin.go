package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-zoo/bone"
	"github.com/mobyos/mobyos-admin-app/server/handlers"
	"github.com/urfave/negroni"
)

func main() {
	mux := bone.New()

	mux.Get("/store", http.HandlerFunc(handlers.GetStoreApps))
	mux.Get("/apps", http.HandlerFunc(handlers.GetApps))
	mux.Post("/apps/:name", http.HandlerFunc(handlers.InstallApp))
	mux.Delete("/apps/:id", http.HandlerFunc(handlers.DeleteApp))
	mux.Get("/apps/:id", http.HandlerFunc(handlers.GetApp))
	mux.Post("/apps/:id/start", http.HandlerFunc(handlers.StartApplication))
	mux.Post("/apps/:id/stop", http.HandlerFunc(handlers.StopApp))
	mux.Get("/*", http.FileServer(http.Dir("./www")))

	n := negroni.Classic()
	n.UseHandler(mux)

	port := "3000"
	if os.Getenv("UBIQ_PORT") != "" {
		port = os.Getenv("UBIQ_PORT")
	}
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), n))
}
