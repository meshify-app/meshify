package core

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"time"

	model "github.com/meshify-app/meshify/model"
	mongo "github.com/meshify-app/meshify/mongo"
	template "github.com/meshify-app/meshify/template"
	util "github.com/meshify-app/meshify/util"
	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"github.com/skip2/go-qrcode"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"gopkg.in/gomail.v2"
)

// CreateHost host with all necessary data
func CreateHost(host *model.Host) (*model.Host, error) {

	u := uuid.NewV4()
	host.Id = u.String()

	if host.HostGroup == "" {
		host.HostGroup = host.Id
	}

	// read the meshes and configure the default values
	meshes, err := ReadMeshes(host.CreatedBy)
	if err != nil {
		return nil, err
	}

	for _, mesh := range meshes {
		if mesh.MeshName == host.MeshName {
			host.Default = mesh.Default
			host.Current = mesh.Default
			host.MeshId = mesh.Id
			host.AccountId = mesh.AccountId
		}
	}

	key, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return nil, err
	}
	host.Current.PrivateKey = key.String()
	host.Current.PublicKey = key.PublicKey().String()

	reserverIps, err := GetAllReservedMeshIps(host.MeshName)
	if err != nil {
		return nil, err
	}

	ips := make([]string, 0)
	for _, network := range host.Default.Address {
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
	host.Current.Address = ips
	host.Current.AllowedIPs = ips
	if host.Current.EnableDns && len(host.Current.Dns) == 0 {
		host.Current.Dns = ips
	}

	if host.Current.SubnetRouting && len(host.Current.PostUp) == 0 {
		host.Current.PostUp = fmt.Sprintf("iptables -A FORWARD -i %s -j ACCEPT; iptables -A FORWARD -o %s -j ACCEPT; iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE", host.MeshName, host.MeshName)
	}
	if host.Current.SubnetRouting && len(host.Current.PostDown) == 0 {
		host.Current.PostDown = fmt.Sprintf("iptables -D FORWARD -i %s -j ACCEPT; iptables -D FORWARD -o %s -j ACCEPT; iptables -t nat -D POSTROUTING -o eth0 -j MASQUERADE", host.MeshName, host.MeshName)
	}

	host.Created = time.Now().UTC()
	host.Updated = host.Created

	// check if host is valid
	errs := host.IsValid()
	if len(errs) != 0 {
		for _, err := range errs {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("host validation error")
		}
		return nil, errors.New("failed to validate host")
	}

	err = mongo.Serialize(host.Id, "id", "hosts", host)
	if err != nil {
		return nil, err
	}

	v, err := mongo.Deserialize(host.Id, "id", "hosts", reflect.TypeOf(model.Host{}))
	if err != nil {
		return nil, err
	}
	host = v.(*model.Host)

	// data modified, dump new config
	return host, UpdateServerConfigWg()
}

// GetAllReservedIps the list of all reserved IPs, client and server
func GetAllReservedMeshIps(meshName string) ([]string, error) {
	clients := mongo.ReadAllHosts("meshName", meshName)
	reserverIps := make([]string, 0)

	for _, client := range clients {
		if client.MeshName == meshName {
			for _, cidr := range client.Current.Address {
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
	}

	return reserverIps, nil
}

// ReadHost host by id
func ReadHost(id string) (*model.Host, error) {
	v, err := mongo.Deserialize(id, "id", "hosts", reflect.TypeOf(model.Host{}))
	if err != nil {
		return nil, err
	}
	host := v.(*model.Host)

	return host, nil
}

// UpdateHost preserve keys
func UpdateHost(Id string, host *model.Host) (*model.Host, error) {
	v, err := mongo.Deserialize(Id, "id", "hosts", reflect.TypeOf(model.Host{}))
	if err != nil {
		return nil, err
	}
	current := v.(*model.Host)

	if current.Id != host.Id {
		return nil, errors.New("records Id mismatch")
	}

	if len(host.Current.Address) == 0 ||
		(len(host.Default.Address) > 0 && len(current.Default.Address) > 0 &&
			(host.Default.Address[0] != current.Default.Address[0])) {
		reserverIps, err := GetAllReservedMeshIps(host.MeshName)
		if err != nil {
			return nil, err
		}

		ips := make([]string, 0)
		for _, network := range host.Default.Address {
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
		host.Current.Address = ips
		host.Current.AllowedIPs = ips
		if host.Current.EnableDns && len(host.Current.Dns) == 0 {
			host.Current.Dns = ips
		}
	}

	if host.Current.SubnetRouting && len(host.Current.PostUp) == 0 {
		host.Current.PostUp = fmt.Sprintf("iptables -A FORWARD -i %s -j ACCEPT; iptables -A FORWARD -o %s -j ACCEPT; iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE", host.MeshName, host.MeshName)
	}
	if host.Current.SubnetRouting && len(host.Current.PostDown) == 0 {
		host.Current.PostDown = fmt.Sprintf("iptables -D FORWARD -i %s -j ACCEPT; iptables -D FORWARD -o %s -j ACCEPT; iptables -t nat -D POSTROUTING -o eth0 -j MASQUERADE", host.MeshName, host.MeshName)
	}

	// check if host is valid
	errs := host.IsValid()
	if len(errs) != 0 {
		for _, err := range errs {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("host validation error")
		}
		return nil, errors.New("failed to validate host")
	}

	// keep keys
	host.Current.PrivateKey = current.Current.PrivateKey
	host.Current.PublicKey = current.Current.PublicKey
	host.Updated = time.Now().UTC()

	err = mongo.Serialize(host.Id, "id", "hosts", host)
	if err != nil {
		return nil, err
	}

	v, err = mongo.Deserialize(Id, "id", "hosts", reflect.TypeOf(model.Host{}))
	if err != nil {
		return nil, err
	}
	host = v.(*model.Host)

	// data modified, dump new config
	return host, UpdateServerConfigWg()
}

// DeleteHost from disk
func DeleteHost(id string) error {

	err := mongo.DeleteHost(id, "hosts")
	//	path := filepath.Join(os.Getenv("WG_CONF_DIR"), id)
	//	err := os.Remove(path)
	if err != nil {
		return err
	}

	// data modified, dump new config
	return UpdateServerConfigWg()
}

// ReadHost2 host by param and id
func ReadHost2(param string, id string) ([]*model.Host, error) {
	hosts := make([]*model.Host, 0)
	hosts = mongo.ReadAllHosts(param, id)
	return hosts, nil
}

// ReadHosts all hosts
func ReadHosts() ([]*model.Host, error) {
	hosts := make([]*model.Host, 0)
	hosts = mongo.ReadAllHosts("", "")
	return hosts, nil
}

// ReadHosts all hosts
func ReadHostsForUser(email string) ([]*model.Host, error) {
	accounts, err := mongo.ReadAllAccounts(email)

	results := make([]*model.Host, 0)

	for _, account := range accounts {
		if account.Status != "Pending" {
			hosts := make([]*model.Host, 0)
			hosts = mongo.ReadAllHosts("accountid", account.Parent)
			for _, host := range hosts {
				results = append(results, host)
			}
		}
	}

	return results, err
}

// ReadHostConfig in wg format
func ReadHostConfig(id string) ([]byte, error) {
	host, err := ReadHost(id)
	if err != nil {
		return nil, err
	}

	server, err := ReadServer()
	if err != nil {
		return nil, err
	}

	configDataWg, err := template.DumpClientWg(host, server)
	if err != nil {
		return nil, err
	}

	return configDataWg, nil
}

// EmailHost send email to host
func EmailHost(id string) error {
	host, err := ReadHost(id)
	if err != nil {
		return err
	}

	configData, err := ReadHostConfig(id)
	if err != nil {
		return err
	}

	// conf as .conf file
	tmpfileCfg, err := ioutil.TempFile("", "wireguard-vpn-*.conf")
	if err != nil {
		return err
	}
	if _, err := tmpfileCfg.Write(configData); err != nil {
		return err
	}
	if err := tmpfileCfg.Close(); err != nil {
		return err
	}
	defer os.Remove(tmpfileCfg.Name()) // clean up

	// conf as png image
	png, err := qrcode.Encode(string(configData), qrcode.Medium, 280)
	if err != nil {
		return err
	}
	tmpfilePng, err := ioutil.TempFile("", "qrcode-*.png")
	if err != nil {
		return err
	}
	if _, err := tmpfilePng.Write(png); err != nil {
		return err
	}
	if err := tmpfilePng.Close(); err != nil {
		return err
	}
	defer os.Remove(tmpfilePng.Name()) // clean up

	// get email body
	emailBody, err := template.DumpEmail(host, filepath.Base(tmpfilePng.Name()))
	if err != nil {
		return err
	}

	// port to int
	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		return err
	}

	d := gomail.NewDialer(os.Getenv("SMTP_HOST"), port, os.Getenv("SMTP_USERNAME"), os.Getenv("SMTP_PASSWORD"))
	s, err := d.Dial()
	if err != nil {
		return err
	}
	m := gomail.NewMessage()

	m.SetHeader("From", os.Getenv("SMTP_FROM"))
	m.SetAddressHeader("To", host.Email, host.Name)
	m.SetHeader("Subject", "WireGuard VPN Configuration")
	m.SetBody("text/html", string(emailBody))
	m.Attach(tmpfileCfg.Name())
	m.Embed(tmpfilePng.Name())

	err = gomail.Send(s, m)
	if err != nil {
		return err
	}

	return nil
}
