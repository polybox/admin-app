package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	yaml "gopkg.in/yaml.v2"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mobyos/mobyos-admin-app/server/types"
)

var dbPath = fmt.Sprintf("%s./db/mobyos.db", os.Getenv("DB_PATH"))

type Scanner interface {
	Scan(...interface{}) error
}

func GetApplications() ([]*types.Application, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query("select id, descriptor from application")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	apps := []*types.Application{}
	for rows.Next() {

		app, err := createApplicationFromRow(rows)
		if err != nil {
			return nil, err
		}
		apps = append(apps, app)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return apps, nil
}

func createApplicationFromRow(row Scanner) (*types.Application, error) {
	app := &types.Application{}
	var desc []byte
	err := row.Scan(&app.Id, &desc)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(desc, &app.Descriptor)
	if err != nil {
		return nil, err
	}

	app.Name = app.Descriptor.Name
	app.IconUrl = app.Descriptor.IconUrl
	app.Description = app.Descriptor.Description
	app.RemotePath = app.Descriptor.RemotePath

	return app, nil
}

func GetApplication(appId string) (*types.Application, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	row := db.QueryRow("select id, descriptor from application where id = ?", appId)

	return createApplicationFromRow(row)
}

func CreateApplication(name string) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	storeApp := app_descriptors[name]
	appDesc := types.AppDescriptor{}
	err = yaml.Unmarshal(storeApp, &appDesc)

	desc, err := appDesc.GetBytes()
	if err != nil {
		return err
	}
	_, err = db.Exec("insert into application values (?, ?)", appDesc.GetId(), desc)
	if err != nil {
		return err
	}

	return nil
}

var app_descriptors map[string][]byte = map[string][]byte{"Spotify": []byte(`
name: "Spotify"
description: "Music for everyone"
icon_url: "http://icons.iconarchive.com/icons/osullivanluke/orb-os-x/128/Spotify-icon.png"
remote_path: "/musicbox_webclient"
services:
  app:
    image: mopidy
    sound: true
    ports:
        - "6680"
`),
	"Retropie": []byte(`
name: "Retropie"
description: "Play your favourite Arcade, home-console, and classic PC games with the minimum set-up."
icon_url: "https://retroresolution.files.wordpress.com/2016/03/retropie_logo_300x300.png"
remote_path: "/"
services:
  app:
    image: retropie
    sound: true
    input: true
    ports:
        - "8080"
`),
	"Kodi": []byte(`
name: "Kodi"
description: "Kodi, the one and only media center"
icon_url: "http://www.homemediatech.net/wp-content/uploads/2015/11/kodi-logo.png"
remote_path: "/"
services:
  app:
    image: marcosnils/kodi
    ui: true
    sound: true
    ports:
        - "8080"
    volumes:
        - "/root/.kodi"
`),
	"Motion": []byte(`
name: "Motion"
description: "See everything"
icon_url: "https://cdn4.iconfinder.com/data/icons/technology-devices-1/500/security-camera-128.png"
remote_path: "/"
services:
  app:
    image: jshridha/motioneye
    ui: true
    sound: true
    ports:
        - "8081"
    volumes:
        - "/home/nobody/media"
        - "/config"
`),
	"Hotspot": []byte(`
name: "Hotspot"
description: "Share your internet connection with your guests"
icon_url: "http://www.montclair-hostel.com/wp-content/uploads/2015/03/wifi.png"
remote_path: "/"
`),
	"BlueJeans": []byte(`
name: "BlueJeans"
description: "BlueJeans conference"
icon_url: "https://logo.clearbit.com/bluejeans.com"
remote_path: "/"
services:
  app:
    image: resin/rpi-raspbian
`),
	"Mantika VPN": []byte(`
name: "Mantika VPN"
description: "Connect to company VPN"
icon_url: "http://soporte.fen.uchile.cl/mw/images/1/18/Vpnfen.png"
remote_path: "/"
`),
	"Metrics": []byte(`
name: "Metrics"
description: "Grafana dashboards"
icon_url: "https://demo.lightbend.com/grafana/public/img/grafana_icon.svg"
remote_path: "/"
`)}

// GetStoreApps is a hardcoded API that returns store apps
func GetStoreApps() ([]*types.AppDescriptor, error) {
	apps := []*types.AppDescriptor{}
	for _, descriptor := range app_descriptors {
		desc := types.AppDescriptor{}
		err := yaml.Unmarshal(descriptor, &desc)

		if err != nil {
			log.Fatal(err)
		}

		apps = append(apps, &desc)

	}
	return apps, nil
}

func DeleteApplication(appId string) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("delete from application where id = (?)", appId)
	if err != nil {
		return err
	}

	return nil
}
