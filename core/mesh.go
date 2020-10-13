package core

import (
	"errors"
	"sort"
	"time"

	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"gitlab.127-0-0-1.fr/vx3r/wg-gen-web/model"
	"gitlab.127-0-0-1.fr/vx3r/wg-gen-web/mongo"
	"gitlab.127-0-0-1.fr/vx3r/wg-gen-web/template"
	"gitlab.127-0-0-1.fr/vx3r/wg-gen-web/util"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

// CreateMesh mesh with all necessary data
func CreateMesh(client *model.Client) (*model.Client, error) {

	u := uuid.NewV4()
	client.Id = u.String()

	key, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return nil, err
	}
	client.PrivateKey = key.String()
	client.PublicKey = key.PublicKey().String()

	presharedKey, err := wgtypes.GenerateKey()
	if err != nil {
		return nil, err
	}
	client.PresharedKey = presharedKey.String()

	reserverIps, err := GetAllReservedIps()
	if err != nil {
		return nil, err
	}

	ips := make([]string, 0)
	for _, network := range client.Address {
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
	client.Address = ips
	client.AllowedIPs = ips
	client.Created = time.Now().UTC()
	client.Updated = client.Created

	// check if client is valid
	errs := client.IsValid()
	if len(errs) != 0 {
		for _, err := range errs {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("client validation error")
		}
		return nil, errors.New("failed to validate client")
	}

	err = mongo.Serialize(client.Id, "mesh", client)
	if err != nil {
		return nil, err
	}

	v, err := mongo.Deserialize(client.Id, "mesh")
	if err != nil {
		return nil, err
	}
	client = v.(*model.Client)

	// data modified, dump new config
	return client, UpdateServerConfigWg()
}

// ReadMesh client by id
func ReadMesh(id string) (*model.Client, error) {
	v, err := mongo.Deserialize(id, "mesh")
	if err != nil {
		return nil, err
	}
	client := v.(*model.Client)

	return client, nil
}

// UpdateMesh preserve keys
func UpdateMesh(Id string, client *model.Client) (*model.Client, error) {
	v, err := mongo.Deserialize(Id, "mesh")
	if err != nil {
		return nil, err
	}
	current := v.(*model.Client)

	if current.Id != client.Id {
		return nil, errors.New("records Id mismatch")
	}

	// check if client is valid
	errs := client.IsValid()
	if len(errs) != 0 {
		for _, err := range errs {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("client validation error")
		}
		return nil, errors.New("failed to validate client")
	}

	// keep keys
	client.PrivateKey = current.PrivateKey
	client.PublicKey = current.PublicKey
	client.Updated = time.Now().UTC()

	err = mongo.Serialize(client.Id, "mesh", client)
	if err != nil {
		return nil, err
	}

	v, err = mongo.Deserialize(Id, "mesh")
	if err != nil {
		return nil, err
	}
	client = v.(*model.Client)

	// data modified, dump new config
	return client, UpdateServerConfigWg()
}

// DeleteMesh from disk
func DeleteMesh(id string) error {

	err := mongo.DeleteClient(id, "mesh")
	//	path := filepath.Join(os.Getenv("WG_CONF_DIR"), id)
	//	err := os.Remove(path)
	if err != nil {
		return err
	}

	// data modified, dump new config
	return UpdateServerConfigWg()
}

// ReadMeshes all clients
func ReadMeshes() ([]*model.Client, error) {
	clients := make([]*model.Client, 0)
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
					clients = append(clients, c.(*model.Client))
				}
			}
		}
	*/
	clients = mongo.ReadAllMeshes()

	sort.Slice(clients, func(i, j int) bool {
		return clients[i].Created.After(clients[j].Created)
	})

	return clients, nil
}

// ReadMeshConfig in wg format
func ReadMeshConfig(id string) ([]byte, error) {
	client, err := ReadClient(id)
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
