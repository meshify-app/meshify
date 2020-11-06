package core

import (
	"errors"
	"reflect"
	"sort"
	"time"

	model "github.com/alan-grapid/meshify/model"
	mongo "github.com/alan-grapid/meshify/mongo"
	template "github.com/alan-grapid/meshify/template"
	util "github.com/alan-grapid/meshify/util"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

// CreateMesh mesh with all necessary data
func CreateMesh(mesh *model.Mesh) (*model.Mesh, error) {

	u := uuid.NewV4()
	mesh.MeshID = u.String()

	reserverIps, err := GetAllReservedIps()
	if err != nil {
		return nil, err
	}

	ips := make([]string, 0)
	for _, network := range mesh.Default.Address {
		ip, err := util.GetAvailableIp(network, reserverIps)
		if err != nil {
			return nil, err
		}
		if util.IsIPv6(ip) {
			ip = ip + "/128"
		} else {
			ip = ip + "/32"
		}
		ips = append(ips, ip)
	}
	mesh.Default.Address = ips
	mesh.Default.AllowedIPs = ips
	mesh.Created = time.Now().UTC()
	mesh.Updated = mesh.Created

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

	err = mongo.Serialize(mesh.MeshID, "meshid", "mesh", mesh)
	if err != nil {
		return nil, err
	}

	v, err := mongo.Deserialize(mesh.MeshID, "meshid", "mesh", reflect.TypeOf(model.Mesh{}))
	if err != nil {
		return nil, err
	}
	mesh = v.(*model.Mesh)

	// data modified, dump new config
	return mesh, UpdateServerConfigWg()
}

// ReadMesh mesh by id
func ReadMesh(id string) (*model.Mesh, error) {
	v, err := mongo.Deserialize(id, "meshid", "mesh", reflect.TypeOf(model.Mesh{}))
	if err != nil {
		return nil, err
	}
	mesh := v.(*model.Mesh)

	return mesh, nil
}

// UpdateMesh preserve keys
func UpdateMesh(Id string, mesh *model.Mesh) (*model.Mesh, error) {
	v, err := mongo.Deserialize(Id, "meshid", "mesh", reflect.TypeOf(model.Mesh{}))
	if err != nil {
		return nil, err
	}
	//	current := v.(*model.Mesh)

	if v == nil {
		return nil, errors.New("Mesh is nil")
		//		x: = fmt.Sprintf("could not retrieve mesh %s", Id)
		//		return nil, errors.New(x)
	}

	//	if current.MeshID != Id {
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

	err = mongo.Serialize(mesh.MeshID, "meshid", "mesh", mesh)
	if err != nil {
		return nil, err
	}

	v, err = mongo.Deserialize(Id, "meshid", "mesh", reflect.TypeOf(model.Mesh{}))
	if err != nil {
		return nil, err
	}
	mesh = v.(*model.Mesh)

	// data modified, dump new config
	return mesh, UpdateServerConfigWg()
}

// DeleteMesh from disk
func DeleteMesh(id string) error {

	err := mongo.Delete(id, "meshid", "mesh")
	//	path := filepath.Join(os.Getenv("WG_CONF_DIR"), id)
	//	err := os.Remove(path)
	if err != nil {
		return err
	}

	// data modified, dump new config
	return UpdateServerConfigWg()
}

// ReadMeshes all clients
func ReadMeshes() ([]*model.Mesh, error) {
	meshes := make([]*model.Mesh, 0)
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
	meshes = mongo.ReadAllMeshes()

	sort.Slice(meshes, func(i, j int) bool {
		return meshes[i].Created.After(meshes[j].Created)
	})

	return meshes, nil
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
