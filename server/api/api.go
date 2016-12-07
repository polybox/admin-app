package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/mobyos/mobyos-admin-app/server/types"
)

func GetApplications() ([]*types.Application, error) {
	r, err := http.Get(fmt.Sprintf("%s/apps", os.Getenv("UBIQ_REMOTE_API")))
	if err != nil {
		return nil, err
	}
	if r.StatusCode != 200 {
		return nil, fmt.Errorf("Error returning apps from remote endpoint. Status [%d].", r.StatusCode)
	}
	var apps []*types.Application

	jsonErr := json.NewDecoder(r.Body).Decode(&apps)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return apps, nil
}
