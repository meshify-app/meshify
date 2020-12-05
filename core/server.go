package core

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"

	model "github.com/grapid/meshify/model"
	storage "github.com/grapid/meshify/storage"
	template "github.com/grapid/meshify/template"
	util "github.com/grapid/meshify/util"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

// ReadServer object, create default one
func ReadServer() (*model.Server, error) {
	if !util.FileExists(filepath.Join(os.Getenv("WG_CONF_DIR"), "mesh.json")) {
		server := &model.Server{}

		key, err := wgtypes.GeneratePrivateKey()
		if err != nil {
			return nil, err
		}
		server.PrivateKey = key.String()
		server.PublicKey = key.PublicKey().String()

		server.Endpoint = ""
		server.ListenPort = 51820

		u := uuid.NewV4()
		server.MeshID = u.String()
		server.MeshName = "Meshify"

		server.Address = make([]string, 0)
		//		server.Address = append(server.Address, "fd9f:6666::10:6:6:1/64")
		server.Address = append(server.Address, "100.0.0.0/24")

		server.Dns = make([]string, 0)
		//		server.Dns = append(server.Dns, "fd9f::10:0:0:2")
		//		server.Dns = append(server.Dns, "10.0.0.2")

		server.AllowedIPs = make([]string, 0)
		//server.AllowedIPs = append(server.AllowedIPs, "0.0.0.0/0")
		//server.AllowedIPs = append(server.AllowedIPs, "::/0")

		server.PersistentKeepalive = 29
		server.Mtu = 0
		server.PreUp = "echo WireGuard PreUp"
		server.PostUp = "echo WireGuard PostUp"
		server.PreDown = "echo WireGuard PreDown"
		server.PostDown = "echo WireGuard PostDown"
		server.Created = time.Now().UTC()
		server.Updated = server.Created

		err = storage.Serialize("mesh.json", server)
		if err != nil {
			return nil, err
		}

		// mesh.json was missing, dump wg config after creation
		err = UpdateServerConfigWg()
		if err != nil {
			return nil, err
		}
	}

	c, err := storage.Deserialize("mesh.json")
	if err != nil {
		return nil, err
	}

	return c.(*model.Server), nil
}

// UpdateServer keep private values from existing one
func UpdateServer(server *model.Server) (*model.Server, error) {
	current, err := storage.Deserialize("mesh.json")
	if err != nil {
		return nil, err
	}

	// check if server is valid
	errs := server.IsValid()
	if len(errs) != 0 {
		for _, err := range errs {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("server validation error")
		}
		return nil, errors.New("failed to validate server")
	}

	server.PrivateKey = current.(*model.Server).PrivateKey
	server.PublicKey = current.(*model.Server).PublicKey
	//server.PresharedKey = current.(*model.Server).PresharedKey
	server.Updated = time.Now().UTC()

	err = storage.Serialize("mesh.json", server)
	if err != nil {
		return nil, err
	}

	v, err := storage.Deserialize("mesh.json")
	if err != nil {
		return nil, err
	}
	server = v.(*model.Server)

	return server, UpdateServerConfigWg()
}

// UpdateServerConfigWg in wg format
func UpdateServerConfigWg() error {
	clients, err := ReadHosts()
	if err != nil {
		return err
	}

	server, err := ReadServer()
	if err != nil {
		return err
	}

	_, err = template.DumpServerWg(clients, server)
	if err != nil {
		return err
	}

	return nil
}

// GetAllReservedIps the list of all reserved IPs, client and server
func GetAllReservedIps() ([]string, error) {
	clients, err := ReadHosts()
	if err != nil {
		return nil, err
	}

	server, err := ReadServer()
	if err != nil {
		return nil, err
	}

	reserverIps := make([]string, 0)

	for _, client := range clients {
		for _, cidr := range client.Address {
			ip, err := util.GetIpFromCidr(cidr)
			if err != nil {
				log.WithFields(log.Fields{
					"err":  err,
					"cidr": cidr,
				}).Error("failed to ip from cidr")
			} else {
				reserverIps = append(reserverIps, ip)
			}
		}
	}

	for _, cidr := range server.Address {
		ip, err := util.GetIpFromCidr(cidr)
		if err != nil {
			log.WithFields(log.Fields{
				"err":  err,
				"cidr": err,
			}).Error("failed to ip from cidr")
		} else {
			reserverIps = append(reserverIps, ip)
		}
	}

	return reserverIps, nil
}

// ReadWgConfigFile return content of wireguard config file
func ReadWgConfigFile() ([]byte, error) {
	return util.ReadFile(filepath.Join(os.Getenv("WG_CONF_DIR"), os.Getenv("WG_INTERFACE_NAME")))
}
