package client

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	auth "github.com/meshify-app/meshify/auth"
	core "github.com/meshify-app/meshify/core"
	model "github.com/meshify-app/meshify/model"
	util "github.com/meshify-app/meshify/util"
	log "github.com/sirupsen/logrus"
	"github.com/skip2/go-qrcode"
	"golang.org/x/oauth2"
)

//var statusCache *cache.Cache

// ApplyRoutes applies router to gin Router
func ApplyRoutes(r *gin.RouterGroup) {
	g := r.Group("/host")
	{

		g.POST("", createHost)
		g.GET("/:id", readHost)
		g.PATCH("/:id", updateHost)
		g.DELETE("/:id", deleteHost)
		g.GET("", readHosts)
		g.GET("/:id/config", configHost)
		g.GET("/:id/status", statusHost)
		g.GET("/:id/email", emailHost)
	}

	//	statusCache = cache.New(1*time.Minute, 10*time.Minute)
}

func createHost(c *gin.Context) {
	var data model.Host

	if err := c.ShouldBindJSON(&data); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to bind")
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	a := util.GetCleanAuthToken(c)
	log.Infof("%v", a)
	// get creation user from token and add to client infos
	oauth2Token := c.MustGet("oauth2Token").(*oauth2.Token)
	oauth2Client := c.MustGet("oauth2Client").(auth.Auth)
	user, err := oauth2Client.UserInfo(oauth2Token)
	if err != nil {
		log.WithFields(log.Fields{
			"oauth2Token": oauth2Token,
			"err":         err,
		}).Error("failed to get user with oauth token")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	data.CreatedBy = user.Email
	if data.AccountId == "" {
		data.AccountId = user.AccountId
	}
	data.APIKey, err = util.RandomString(32)

	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to generate state random string")
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	client, err := core.CreateHost(&data)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to create client")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, client)
}

func readHost(c *gin.Context) {
	id := c.Param("id")

	client, err := core.ReadHost(id)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to read client")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, client)
}

func updateHost(c *gin.Context) {
	var data model.Host
	id := c.Param("id")
	if id == "" {
		log.Error("hostid cannot be empty")
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	if err := c.ShouldBindJSON(&data); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to bind")
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	apikey := c.Request.Header.Get("X-API-KEY")

	if apikey != "" {

		host, err := core.ReadHost(id)
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("failed to read client config")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		authorized := false

		if host.APIKey == apikey {
			authorized = true
		}

		if !authorized {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	} else {
		// get update user from token and add to client infos
		oauth2Token := c.MustGet("oauth2Token").(*oauth2.Token)
		oauth2Client := c.MustGet("oauth2Client").(auth.Auth)
		user, err := oauth2Client.UserInfo(oauth2Token)
		if err != nil {
			log.WithFields(log.Fields{
				"oauth2Token": oauth2Token,
				"err":         err,
			}).Error("failed to get user with oauth token")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		data.UpdatedBy = user.Name
	}

	client, err := core.UpdateHost(id, &data)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to update host")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, client)
}

func deleteHost(c *gin.Context) {
	id := c.Param("id")
	// get update user from token and add to client infos

	oauth2Token := c.MustGet("oauth2Token").(*oauth2.Token)
	oauth2Client := c.MustGet("oauth2Client").(auth.Auth)
	user, err := oauth2Client.UserInfo(oauth2Token)
	if err != nil {
		log.WithFields(log.Fields{
			"oauth2Token": oauth2Token,
			"err":         err,
		}).Error("failed to get user with oauth token")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	log.Infof("User %s deleted host %s", user.Name, id)

	err = core.DeleteHost(id)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to remove client")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func readHosts(c *gin.Context) {
	value, exists := c.Get("oauth2Token")
	if !exists {
		c.AbortWithStatus(401)
		return
	}
	oauth2Token := value.(*oauth2.Token)
	oauth2Client := c.MustGet("oauth2Client").(auth.Auth)
	user, err := oauth2Client.UserInfo(oauth2Token)
	if err != nil {
		log.WithFields(log.Fields{
			"oauth2Token": oauth2Token,
			"err":         err,
		}).Error("failed to get user with oauth token")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if user.Email == "" {
		log.Error("security alert: Email empty on authenticated token")
		c.AbortWithStatus(http.StatusForbidden)
	}
	clients, err := core.ReadHostsForUser(user.Email)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to list clients")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, clients)
}

func statusHost(c *gin.Context) {

	//	id := c.Param("id")
	if c.Param("id") == "" {
		log.Error("hostgroup cannot be empty")
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	hostGroup := c.Param("id")

	apikey := c.Request.Header.Get("X-API-KEY")
	etag := c.Request.Header.Get("If-None-Match")

	/*
		m, _ := statusCache.Get(id)
		if m != nil {
			msg := m.(model.Message)
			authorized := false

			for _, config := range msg.Config {
				for _, mesh := range config.Hosts {
					if mesh.HostGroup == id && mesh.APIKey == apikey {
						authorized = true
						break
					}
				}
			}
			if !authorized {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			c.JSON(http.StatusOK, m)
			return
		}
	*/

	meshes, err := core.ReadHost2("hostGroup", c.Param("id"))
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to read client config")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	authorized := false

	for _, mesh := range meshes {
		if mesh.APIKey == apikey {
			authorized = true
			break
		}
	}
	if !authorized {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var msg model.Message
	hconfig := make([]model.HostConfig, len(meshes))

	msg.Id = c.Param("id")
	msg.Config = hconfig

	for i, mesh := range meshes {
		clients, err := core.ReadHost2("meshid", mesh.MeshId)
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("failed to list clients")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		msg.Config[i] = model.HostConfig{}
		msg.Config[i].MeshName = mesh.MeshName
		msg.Config[i].MeshId = mesh.MeshId

		for _, client := range clients {
			// They should all match
			if client.MeshId == msg.Config[i].MeshId {
				// If this config isn't explicitly for this host, remove the private key
				// and api key from the results
				if client.HostGroup != hostGroup {
					client.Current.PrivateKey = ""
					client.APIKey = ""
				}
				msg.Config[i].Hosts = append(msg.Config[i].Hosts, *client)
			} else {
				log.Errorf("internal error")
			}
		}
	}
	bytes, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("cannot marshal msg %v", err)
	}
	md5 := fmt.Sprintf("%x", md5.Sum(bytes))
	if md5 == etag {
		c.AbortWithStatus(http.StatusNotModified)
	} else {
		c.Header("ETag", md5)
		c.JSON(http.StatusOK, msg)
	}

	//	statusCache.Set(id, msg, 0)
}

func configHost(c *gin.Context) {

	formatQr := c.DefaultQuery("qrcode", "false")
	data, err := core.ReadHostConfig(c.Param("id"))
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to read host config")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	sdata := string(data)

	if formatQr == "false" {
		// return config as txt file
		c.Header("Content-Disposition", "attachment; filename=meshify.conf")
		c.Data(http.StatusOK, "application/config", data)
		return
	}
	// return config as png qrcode
	png, err := qrcode.Encode(sdata, qrcode.Medium, 250)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to create qrcode")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.Data(http.StatusOK, "image/png", png)

	return
}

func emailHost(c *gin.Context) {
	id := c.Param("id")

	err := core.EmailHost(id)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to send email to client")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
