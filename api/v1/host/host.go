package client

import (
	"net/http"

	"github.com/gin-gonic/gin"
	auth "github.com/meshify-app/meshify/auth"
	core "github.com/meshify-app/meshify/core"
	model "github.com/meshify-app/meshify/model"
	util "github.com/meshify-app/meshify/util"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

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
		c.AbortWithStatus(http.StatusInternalServerError)
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

	if err := c.ShouldBindJSON(&data); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to bind")
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	// get update user from token and add to client infos
	oauth2Token := c.MustGet("oauth2Token").(*oauth2.Token)
	oauth2Client := c.MustGet("oauth2Client").(auth.Auth)
	user, err := oauth2Client.UserInfo(oauth2Token)
	if err != nil {
		log.WithFields(log.Fields{
			"oauth2Token": oauth2Token,
			"err":         err,
		}).Error("failed to get user with oauth token")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	data.UpdatedBy = user.Name

	client, err := core.UpdateHost(id, &data)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to update client")
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
		c.AbortWithStatus(http.StatusInternalServerError)
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
		c.AbortWithStatus(http.StatusInternalServerError)
		return
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

	if c.Param("id") == "" {
		log.Error("hostgroup cannot be empty")
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	meshes, err := core.ReadHost2("hostGroup", c.Param("id"))
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to read client config")
		c.AbortWithStatus(http.StatusInternalServerError)
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
				msg.Config[i].Hosts = append(msg.Config[i].Hosts, *client)
			}
		}
	}

	c.JSON(http.StatusOK, msg)
}

func configHost(c *gin.Context) {
	configData, err := core.ReadHost2("id", c.Param("id"))
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to read client config")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, configData)

	/*	formatQr := c.DefaultQuery("qrcode", "false")
		if formatQr == "false" {
			// return config as txt file
			c.Header("Content-Disposition", "attachment; filename=wg0.conf")
			c.Data(http.StatusOK, "application/config", configData)
			return
		}
		// return config as png qrcode
		png, err := qrcode.Encode(string(configData), qrcode.Medium, 250)
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Error("failed to create qrcode")
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		c.Data(http.StatusOK, "image/png", png)*/

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
