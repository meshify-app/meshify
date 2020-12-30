package core

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"time"

	model "github.com/grapid/meshify/model"
	mongo "github.com/grapid/meshify/mongo"
	template "github.com/grapid/meshify/template"
	util "github.com/grapid/meshify/util"
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

	key, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return nil, err
	}
	host.Current.PrivateKey = key.String()
	host.Current.PublicKey = key.PublicKey().String()

	presharedKey, err := wgtypes.GenerateKey()
	if err != nil {
		return nil, err
	}
	host.Current.PresharedKey = presharedKey.String()

	reserverIps, err := GetAllReservedIps()
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

// ReadHosts all hosts
func ReadHosts() ([]*model.Host, error) {
	hosts := make([]*model.Host, 0)
	/*
		files, err := ioutil.ReadDir(filepath.Join(os.Getenv("WG_CONF_DIR")))
		if err != nil {
			return nil, err
		}

		for _, f := range files {
			// hosts file name is an uuid
			_, err := uuid.FromString(f.Name())
			if err == nil {
				c, err := mongo.Deserialize(f.Name())
				if err != nil {
					log.WithFields(log.Fields{
						"err":  err,
						"path": f.Name(),
					}).Error("failed to deserialize host")
				} else {
					hosts = append(hosts, c.(*model.Host))
				}
			}
		}
	*/
	hosts = mongo.ReadAllHosts()

	sort.Slice(hosts, func(i, j int) bool {
		return hosts[i].Created.After(hosts[j].Created)
	})

	return hosts, nil
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
