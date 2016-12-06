package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/go-zoo/bone"
	"github.com/mobyos/mobyos-admin-app/server/db"
	"github.com/mobyos/mobyos-admin-app/server/docker"
)

func GetApps(rw http.ResponseWriter, req *http.Request) {
	installations, err := db.GetApplications()
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
	}

	err = docker.SetContainerStates(installations)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := docker.SetApplicationsAreLocal(installations); err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(rw).Encode(installations)
}

func GetStoreApps(rw http.ResponseWriter, req *http.Request) {

	storeApps, err := db.GetStoreApps()
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(rw).Encode(storeApps)
}

func GetApp(rw http.ResponseWriter, req *http.Request) {
	appId := bone.GetValue(req, "id")

	app, err := db.GetApplication(appId)

	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	} else if app == nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	err = docker.SetContainerState(app)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(rw).Encode(app)
}

func InstallApp(rw http.ResponseWriter, req *http.Request) {
	name := bone.GetValue(req, "name")
	err := db.CreateApplication(name)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)

}

func StartApplication(rw http.ResponseWriter, req *http.Request) {
	appId := bone.GetValue(req, "id")
	app, err := db.GetApplication(appId)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	} else if app == nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = docker.RunApp(app)
	if err != nil && strings.Contains(err.Error(), "in use") {
		log.Println(err)
		rw.WriteHeader(http.StatusConflict)
		return
	} else if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)
	json.NewEncoder(rw).Encode(app)
}

func StopApp(rw http.ResponseWriter, req *http.Request) {
	appId := bone.GetValue(req, "id")
	app, err := db.GetApplication(appId)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	} else if app == nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	err = docker.StopApp(app)
	if err != nil && strings.Contains(err.Error(), "No such container") {
		log.Println(err)
		rw.WriteHeader(http.StatusConflict)
		return
	} else if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
}

func DeleteApp(rw http.ResponseWriter, req *http.Request) {
	appId := bone.GetValue(req, "id")
	app, err := db.GetApplication(appId)

	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	} else if app == nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	err = docker.SetContainerState(app)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if app.IsRunning {
		// can't delete apps that are running need to stop them first
		rw.WriteHeader(http.StatusConflict)
		return
	}

	err = db.DeleteApplication(app.Id)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

}
