package storage

import (
	"encoding/json"
	"os"
	"path/filepath"

	model "github.com/meshify-app/meshify/model"
	util "github.com/meshify-app/meshify/util"
)

// Serialize write interface to disk
func Serialize(id string, c interface{}) error {
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return util.WriteFile(filepath.Join(os.Getenv("WG_CONF_DIR"), id), b)
}

// Deserialize read interface from disk
func Deserialize(id string) (interface{}, error) {
	path := filepath.Join(os.Getenv("WG_CONF_DIR"), id)

	data, err := util.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if id == "mesh.json" {
		var s *model.Server
		err = json.Unmarshal(data, &s)
		if err != nil {
			return nil, err
		}
		return s, nil
	}

	// if not the server, must be client
	var c *model.Host
	err = json.Unmarshal(data, &c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
