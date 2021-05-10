package core

import (
	"errors"
	"reflect"
	"sort"
	"time"

	model "github.com/meshify-app/meshify/model"
	mongo "github.com/meshify-app/meshify/mongo"
	template "github.com/meshify-app/meshify/template"
	util "github.com/meshify-app/meshify/util"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

// CreateMesh mesh with all necessary data
func CreateMesh(mesh *model.Mesh) (*model.Mesh, error) {

	u := uuid.NewV4()
	mesh.Id = u.String()

	ips := make([]string, 0)
	for _, network := range mesh.Default.Address {
		ip, err := util.GetNetworkAddress(network)
		if err != nil {
			return nil, err
		}
		if util.IsIPv6(ip) {
			ip = ip + "/64"
		} else {
			ip = ip + "/24"
		}
		ips = append(ips, ip)
	}

	mesh.Default.Address = ips
	if len(mesh.Default.AllowedIPs) == 0 {
		mesh.Default.AllowedIPs = ips
	}

	mesh.Created = time.Now().UTC()
	mesh.Updated = mesh.Created

	if mesh.Default.EnableDns {
		if len(mesh.Default.Dns) == 0 {
			mesh.Default.Dns = ips
		}
	}

	if mesh.Default.PresharedKey == "" {
		presharedKey, err := wgtypes.GenerateKey()
		if err != nil {
			return nil, err
		}
		mesh.Default.PresharedKey = presharedKey.String()
	}

	// check if mesh is valid
	errs := mesh.IsValid()
	if len(errs) != 0 {
		for _, err := range errs {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("mesh validation error")
		}
		return nil, errors.New("failed to validate mesh")
	}

	err := mongo.Serialize(mesh.Id, "id", "mesh", mesh)
	if err != nil {
		return nil, err
	}

	v, err := mongo.Deserialize(mesh.Id, "id", "mesh", reflect.TypeOf(model.Mesh{}))
	if err != nil {
		return nil, err
	}
	mesh = v.(*model.Mesh)

	// data modified, dump new config
	return mesh, UpdateServerConfigWg()
}

// ReadMesh mesh by id
func ReadMesh(id string) (*model.Mesh, error) {
	v, err := mongo.Deserialize(id, "id", "mesh", reflect.TypeOf(model.Mesh{}))
	if err != nil {
		return nil, err
	}
	mesh := v.(*model.Mesh)

	return mesh, nil
}

// UpdateMesh preserve keys
func UpdateMesh(Id string, mesh *model.Mesh) (*model.Mesh, error) {
	v, err := mongo.Deserialize(Id, "id", "mesh", reflect.TypeOf(model.Mesh{}))
	if err != nil {
		return nil, err
	}
	//	current := v.(*model.Mesh)

	if v == nil {
		return nil, errors.New("Mesh is nil")
		//		x: = fmt.Sprintf("could not retrieve mesh %s", Id)
		//		return nil, errors.New(x)
	}

	//	if current.ID != Id {
	//		return nil, errors.New("records Id mismatch")
	//	}

	// check if mesh is valid
	errs := mesh.IsValid()
	if len(errs) != 0 {
		for _, err := range errs {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("mesh validation error")
		}
		return nil, errors.New("failed to validate mesh")
	}

	mesh.Updated = time.Now().UTC()

	if mesh.Default.EnableDns {
		if len(mesh.Default.Dns) == 0 {
			mesh.Default.Dns = mesh.Default.Address
		}
	}

	err = mongo.Serialize(mesh.Id, "id", "mesh", mesh)
	if err != nil {
		return nil, err
	}

	v, err = mongo.Deserialize(Id, "id", "mesh", reflect.TypeOf(model.Mesh{}))
	if err != nil {
		return nil, err
	}
	mesh = v.(*model.Mesh)

	// data modified, dump new config
	return mesh, UpdateServerConfigWg()
}

// DeleteMesh from disk
func DeleteMesh(id string) error {

	err := mongo.Delete(id, "id", "mesh")
	//	path := filepath.Join(os.Getenv("WG_CONF_DIR"), id)
	//	err := os.Remove(path)
	if err != nil {
		return err
	}

	// data modified, dump new config
	return UpdateServerConfigWg()
}

// ReadMeshes all clients
func ReadMeshes(email string) ([]*model.Mesh, error) {

	accounts, err := mongo.ReadAllAccounts(email)

	results := make([]*model.Mesh, 0)

	for _, account := range accounts {
		if account.Status != "Pending" {
			meshes := make([]*model.Mesh, 0)
			meshes = mongo.ReadAllMeshes("accountid", account.Parent)
			for _, mesh := range meshes {
				results = append(results, mesh)
			}
		}
	}

	/*
		files, err := ioutil.ReadDir(filepath.Join(os.Getenv("WG_CONF_DIR")))
		if err != nil {
			return nil, err
		}

		for _, f := range files {
			// clients file name is an uuid
			_, err := uuid.FromString(f.Name())
			if err == nil {
				c, err := mongo.Deserialize(f.Name())
				if err != nil {
					log.WithFields(log.Fields{
						"err":  err,
						"path": f.Name(),
					}).Error("failed to deserialize client")
				} else {
					clients = append(clients, c.(*model.Host))
				}
			}
		}
	*/

	sort.Slice(results, func(i, j int) bool {
		return results[i].Created.After(results[j].Created)
	})

	return results, err
}

// ReadMeshConfig in wg format
func ReadMeshConfig(id string) ([]byte, error) {
	client, err := ReadHost(id)
	if err != nil {
		return nil, err
	}

	server, err := ReadServer()
	if err != nil {
		return nil, err
	}

	configDataWg, err := template.DumpClientWg(client, server)
	if err != nil {
		return nil, err
	}

	return configDataWg, nil
}
