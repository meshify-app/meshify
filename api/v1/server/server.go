package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	auth "github.com/meshify-app/meshify/auth"
	core "github.com/meshify-app/meshify/core"
	model "github.com/meshify-app/meshify/model"
	version "github.com/meshify-app/meshify/version"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

// ApplyRoutes applies router to gin Router
func ApplyRoutes(r *gin.RouterGroup) {
	g := r.Group("/server")
	{
		g.GET("", readServer)
		g.PATCH("", updateServer)
		g.GET("/version", versionStr)
	}
}

func readServer(c *gin.Context) {
	client, err := core.ReadServer()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to read client")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, client)
}

func updateServer(c *gin.Context) {
	var data model.Server

	if err := c.ShouldBindJSON(&data); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to bind")
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	// get update user from token and add to server infos
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
	if user.Email == "" {
		log.WithFields(log.Fields{
			"user": user,
		}).Error("user has no email")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	server, err := core.UpdateServer(&data)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to update client")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, server)
}

func versionStr(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version": version.Version,
	})
}
