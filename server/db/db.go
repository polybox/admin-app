package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mobyos/mobyos-admin-app/server/types"
	"github.com/twinj/uuid"
)

var dbPath = fmt.Sprintf("%s./db/mobyos.db", os.Getenv("DB_PATH"))

func init() {
	uuid.SwitchFormat(uuid.FormatHex)
}

func GetApplications() ([]*types.Application, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	rows, err := db.Query("select id, name, icon_url from application")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	apps := []*types.Application{}
	for rows.Next() {

		app := types.Application{}
		err = rows.Scan(&app.Id, &app.Name, &app.IconUrl)
		if err != nil {
			return nil, err
		}
		apps = append(apps, &app)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return apps, nil
}

func GetApplication(appId string) (*types.Application, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	row := db.QueryRow("select id, name, descriptor from application where id = ?", appId)

	app := &types.Application{}
	err = row.Scan(&app.Id, &app.Name, &app.Descriptor)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return app, nil
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
	_, err = db.Exec("insert into application values (?, ?, ?,  ?)", uuid.NewV5(uuid.NameSpaceURL, uuid.Name(appDesc.Name)).String(), appDesc.Name, appDesc.IconUrl, desc)
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
