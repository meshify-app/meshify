package account

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/meshify-app/meshify/auth"
	core "github.com/meshify-app/meshify/core"
	model "github.com/meshify-app/meshify/model"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

// ApplyRoutes applies router to gin Router
func ApplyRoutes(r *gin.RouterGroup) {
	g := r.Group("/accounts")
	{

		g.POST("/", createAccount)
		g.POST("/:id/activate", activateAccount)
		g.PATCH("/:id/activate", activateAccount)
		g.GET("/:id", readAllAccounts)
		g.PATCH("/:id", updateAccount)
		g.DELETE("/:id", deleteAccount)
	}
}

func activateAccount(c *gin.Context) {
	id := c.Param("id")

	v, err := core.ActivateAccount(id)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to create account")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, v)
}

func createAccount(c *gin.Context) {
	var account model.Account

	if err := c.ShouldBindJSON(&account); err != nil {
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

	account.From = user.Email

	v, err := core.CreateAccount(&account)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to create account")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	log.Infof("emailUser account = %v %v %v", account.Email, account.Id, account.MeshId)

	err = core.EmailUser(account.Email, account.Id, account.MeshId)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to send email to client")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, v)
}

func readAllAccounts(c *gin.Context) {
	email := c.Param("id")

	accounts, err := core.ReadAllAccounts(email)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to read accounts")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, accounts)
}

func updateAccount(c *gin.Context) {
	var data model.Account
	id := c.Param("id")

	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
		log.Info("updateAccount - %s", string(bodyBytes))
	}

	// Restore the io.ReadCloser to its original state
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	log.Info("updateAccount - %s", string(bodyBytes))

	if err := c.ShouldBindJSON(&data); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to bind")
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}
	client, err := core.UpdateAccount(id, &data)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to update client")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, client)
}

func deleteAccount(c *gin.Context) {
	id := c.Param("id")

	err := core.DeleteAccount(id)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to remove client")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
