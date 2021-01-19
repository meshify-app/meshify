package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	auth "github.com/meshify-app/meshify/auth"
	core "github.com/meshify-app/meshify/core"
	model "github.com/meshify-app/meshify/model"
	log "github.com/sirupsen/logrus"
	"github.com/skip2/go-qrcode"
	"golang.org/x/oauth2"
)

// ApplyRoutes applies router to gin Router
func ApplyRoutes(r *gin.RouterGroup) {
	g := r.Group("/users")
	{

		g.POST("", inviteUser)
		g.GET("/:id/invite", emailUser)
		g.GET("/:id", readUser)
		g.PATCH("/:id", updateUser)
		g.DELETE("/:id", deleteUser)
		g.GET("", readUsers)
	}
}

func inviteUser(c *gin.Context) {
	var data model.User

	if err := c.ShouldBindJSON(&data); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to bind")
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

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

	u, err := core.CreateUser(&data)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to create user")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, u)
}

func readUser(c *gin.Context) {
	id := c.Param("id")

	client, err := core.ReadUser(id)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to read client")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, client)
}

func updateUser(c *gin.Context) {
	var data model.User
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

	client, err := core.UpdateUser(id, &data)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to update client")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, client)
}

func deleteUser(c *gin.Context) {
	id := c.Param("id")

	err := core.DeleteUser(id)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to remove client")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func readUsers(c *gin.Context) {
	meshes, err := core.ReadUsers()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to list clients")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, meshes)
}

func configUsers(c *gin.Context) {
	configData, err := core.ReadUserConfig(c.Param("id"))
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to read client config")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	formatQr := c.DefaultQuery("qrcode", "false")
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
	c.Data(http.StatusOK, "image/png", png)
	return
}

func emailUser(c *gin.Context) {
	id := c.Param("id")

	err := core.EmailUser(id)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to send email to client")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
