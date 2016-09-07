package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	yaml "gopkg.in/yaml.v2"

	"github.com/go-zoo/bone"
	"github.com/mobyos/admin-app/server/db"
	"github.com/mobyos/admin-app/server/docker"
	"github.com/mobyos/admin-app/server/types"
)

func GetInstalledApps(rw http.ResponseWriter, req *http.Request) {
	installations, err := db.GetInstallations()
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
	}

	err = docker.SetInstallationStatuses(installations)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
	}

	json.NewEncoder(rw).Encode(installations)
}

func InstallApp(rw http.ResponseWriter, req *http.Request) {
	appYaml, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	desc := types.AppDescriptor{}
	err = yaml.Unmarshal(appYaml, &desc)
	if err != nil {
		// TODO return a meaningful error
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	err = db.CreateApplication(desc)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)

}

func StartApp(rw http.ResponseWriter, req *http.Request) {
	appId := bone.GetValue(req, "id")
	app, err := db.GetInstallation(appId)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	} else if app == nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	desc := types.AppDescriptor{}
	err = yaml.Unmarshal(app.Descriptor, &desc)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = docker.RunApp(app.Id, desc)
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
}

func StopApp(rw http.ResponseWriter, req *http.Request) {
	appId := bone.GetValue(req, "id")
	app, err := db.GetInstallation(appId)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	} else if app == nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	err = docker.SetInstallationStatus(app)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if app.Status == "Not running" {
		rw.WriteHeader(http.StatusConflict)
		return
	}

	err = docker.StopApp(app.Id)
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
	app, err := db.GetInstallation(appId)

	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	} else if app == nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}

	err = docker.SetInstallationStatus(app)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	if app.Status != "Not running" {
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
