package db

import (
	"database/sql"
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v2"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mobyos/mobyos-admin-app/server/types"
	"github.com/twinj/uuid"
)

var dbPath = fmt.Sprintf("%s./db/mobyos.db", os.Getenv("DB_PATH"))

func init() {
	uuid.SwitchFormat(uuid.FormatHex)
}

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

func CreateApplication(appDesc types.AppDescriptor) error {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	desc, err := appDesc.GetBytes()
	if err != nil {
		return err
	}
	_, err = db.Exec("insert into application values (?, ?)", uuid.NewV5(uuid.NameSpaceURL, uuid.Name(appDesc.Name)).String(), desc)
	if err != nil {
		return err
	}

	return nil
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
