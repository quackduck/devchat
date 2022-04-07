package v2

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

const (
	defaultAdminsFileName = "admins.json"
)

type (
	AdminID      = string
	AdminInfo    = string
	AdminInfoMap = map[AdminID]AdminInfo
)

func (s *Server) GetAdmins() (AdminInfoMap, error) {
	if _, err := os.Stat(defaultAdminsFileName); err == os.ErrNotExist {
		return nil, errors.New("make an admins.json file to add admins")
	}

	data, err := ioutil.ReadFile("admins.json")
	if err != nil {
		return nil, fmt.Errorf("error reading admins.json: %s", err)
	}

	adminsList := make(AdminInfoMap)

	if err = json.Unmarshal(data, &adminsList); err != nil {
		return nil, fmt.Errorf("bad json: %v", err)
	}

	return adminsList, nil
}
