package core

import (
	"errors"
	"reflect"
	"sort"
	"time"

	model "github.com/meshify-app/meshify/model"
	mongo "github.com/meshify-app/meshify/mongo"
	template "github.com/meshify-app/meshify/template"
	log "github.com/sirupsen/logrus"
)

// CreateService service with all necessary data
func CreateService(service *model.Service) (*model.Service, error) {
	/*
		u := uuid.NewV4()
		service.Id = u.String()

		ips := make([]string, 0)
		for _, network := range service.Default.Address {
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

		service.Default.Address = ips
		if len(service.Default.AllowedIPs) == 0 {
			service.Default.AllowedIPs = ips
		}

		service.Created = time.Now().UTC()
		service.Updated = service.Created

		if service.Default.PresharedKey == "" {
			presharedKey, err := wgtypes.GenerateKey()
			if err != nil {
				return nil, err
			}
			service.Default.PresharedKey = presharedKey.String()
		}

		// check if service is valid
		errs := service.IsValid()
		if len(errs) != 0 {
			for _, err := range errs {
				log.WithFields(log.Fields{
					"err": err,
				}).Error("service validation error")
			}
			return nil, errors.New("failed to validate service")
		}

		err := mongo.Serialize(service.Id, "id", "service", service)
		if err != nil {
			return nil, err
		}

		v, err := mongo.Deserialize(service.Id, "id", "service", reflect.TypeOf(model.Service{}))
		if err != nil {
			return nil, err
		}
		service = v.(*model.Service)

		// data modified, dump new config
		return service, UpdateServerConfigWg()
	*/
	return nil, errors.New("not implemented")
}

// ReadService service by id
func ReadService(id string) (*model.Service, error) {
	v, err := mongo.Deserialize(id, "id", "service", reflect.TypeOf(model.Service{}))
	if err != nil {
		return nil, err
	}
	service := v.(*model.Service)

	return service, nil
}

// UpdateService preserve keys
func UpdateService(Id string, service *model.Service) (*model.Service, error) {
	v, err := mongo.Deserialize(Id, "id", "service", reflect.TypeOf(model.Service{}))
	if err != nil {
		return nil, err
	}
	//	current := v.(*model.Service)

	if v == nil {
		return nil, errors.New("Service is nil")
		//		x: = fmt.Sprintf("could not retrieve service %s", Id)
		//		return nil, errors.New(x)
	}

	//	if current.ID != Id {
	//		return nil, errors.New("records Id mismatch")
	//	}

	// check if service is valid
	errs := service.IsValid()
	if len(errs) != 0 {
		for _, err := range errs {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("service validation error")
		}
		return nil, errors.New("failed to validate service")
	}

	service.Updated = time.Now().UTC()

	err = mongo.Serialize(service.Id, "id", "service", service)
	if err != nil {
		return nil, err
	}

	v, err = mongo.Deserialize(Id, "id", "service", reflect.TypeOf(model.Service{}))
	if err != nil {
		return nil, err
	}
	service = v.(*model.Service)

	// data modified, dump new config
	return service, UpdateServerConfigWg()
}

// DeleteService from disk
func DeleteService(id string) error {

	err := mongo.Delete(id, "id", "service")
	//	path := filepath.Join(os.Getenv("WG_CONF_DIR"), id)
	//	err := os.Remove(path)
	if err != nil {
		return err
	}

	// data modified, dump new config
	return UpdateServerConfigWg()
}

// ReadServicees all clients
func ReadServicees(email string) ([]*model.Service, error) {

	accounts, err := mongo.ReadAllAccounts(email)

	results := make([]*model.Service, 0)

	for _, account := range accounts {
		if account.Id == account.Parent && account.Status == "Active" {
			servicees, err := mongo.ReadAllServices(account.Email)
			if err != nil {
				results = append(results, servicees...)
			}
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Created.After(results[j].Created)
	})

	return results, err
}

// ReadServiceConfig in wg format
func ReadServiceConfig(id string) ([]byte, error) {
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
