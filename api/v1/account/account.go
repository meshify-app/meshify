package account

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/meshify-app/meshify/auth"
	core "github.com/meshify-app/meshify/core"
	model "github.com/meshify-app/meshify/model"
	"github.com/meshify-app/meshify/util"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

// ApplyRoutes applies router to gin Router
func ApplyRoutes(r *gin.RouterGroup) {
	g := r.Group("/accounts")
	{

		g.POST("/:id", createAccount)
		g.GET("/:id", readAllAccountsForUser)
		g.DELETE("/:id", deleteAccount)
	}
}

// ApplyRoutes applies router to gin Router
func ApplyOrgRoutes(r *gin.RouterGroup) {
	g := r.Group("/organizations")
	{

		g.POST("/:id", createOrg)
		g.GET("/:id", readAllOrgsForUser)
		g.DELETE("/:id", deleteOrg)
	}
}

func createOrg(c *gin.Context) {
	var data model.Organization

	if err := c.ShouldBindJSON(&data); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to bind")
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	a := c.Request.Header.Get(util.AuthTokenHeaderName)

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

	org, err := core.CreateOrg(&data)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to create org")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, org)

}

func readAllOrgsForUser(c *gin.Context) {

}
func deleteOrg(c *gin.Context) {

}

func createAccount(c *gin.Context) {
	email := c.Param("id")

	var account model.Account
	account.Email = email

	v, err := core.CreateAccount(&account)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to create account")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, v)
}

func readAllAccountsForUser(c *gin.Context) {
	email := c.Param("id")

	accounts, err := core.ReadAllAccountsForUser(email)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to read accounts")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, accounts)
}

/*
func updateAccount(c *gin.Context) {
	var data model.Account
	id := c.Param("id")

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
*/
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
