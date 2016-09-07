package main

import (
	"log"
	"net/http"

	"github.com/go-zoo/bone"
	"github.com/urfave/negroni"

	"github.com/mobyos/admin/handlers"
)

func main() {
	mux := bone.New()

	mux.Get("/apps/installed", http.HandlerFunc(handlers.GetInstalledApps))
	mux.Post("/apps", http.HandlerFunc(handlers.InstallApp))
	mux.Delete("/apps/:id", http.HandlerFunc(handlers.DeleteApp))
	mux.Post("/apps/:id/start", http.HandlerFunc(handlers.StartApp))
	mux.Post("/apps/:id/stop", http.HandlerFunc(handlers.StopApp))
	mux.Get("/*", http.FileServer(http.Dir("./www")))

	n := negroni.Classic()
	n.UseHandler(mux)

	log.Fatal(http.ListenAndServe(":3000", n))
}
