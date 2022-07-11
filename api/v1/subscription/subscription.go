package subscription

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	auth "github.com/meshify-app/meshify/auth"
	"github.com/meshify-app/meshify/core"
	model "github.com/meshify-app/meshify/model"
	"github.com/meshify-app/meshify/mongo"
	"github.com/meshify-app/meshify/util"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

// ApplyRoutes applies router to gin Router
func ApplyRoutes(r *gin.RouterGroup) {
	g := r.Group("/subscriptions")
	{

		g.POST("", createSubscription)
		g.GET("/:id", readSubscription)
		g.PATCH("/:id", updateSubscription)
		g.DELETE("/:id", deleteSubscription)
		g.GET("", readSubscriptions)
	}
}

func createSubscription(c *gin.Context) {
	var body string
	var sub map[string]interface{}

	// read and log the request body

	bytes, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to read request body")
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	body = string(bytes)
	// remove all the backslashes from the body
	body = strings.Replace(body, "\\", "", -1)
	log.Info(body)
	bytes = []byte(body)

	err = json.Unmarshal(bytes, &sub)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to unmarshal request body")
		// c.AbortWithStatus(http.StatusUnprocessableEntity)
		// return with no error so webhook doesn't get disabled
		c.JSON(http.StatusOK, body)
		return
	}

	log.Info(sub)

	// walk the json and find the customer href
	links := sub["_links"].(map[string]interface{})
	log.Info(links)

	customer := links["customer"].([]interface{})
	log.Info(customer)

	customer0 := customer[0].(map[string]interface{})
	log.Info(customer0)

	href := customer0["href"].(string)
	log.Info(href)

	go func() {

		// make http request with basic authentication using href as url to get the customer object
		req, err := http.NewRequest("GET", href, nil)
		if err != nil {
			return
		}

		req.SetBasicAuth(os.Getenv("WC_USERNAME"), os.Getenv("WC_PASSWORD"))
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Errorf("http.client.Do = %v", err)
			return
		}

		if resp.StatusCode != 200 {
			log.Errorf("http status %s expect 200 OK", resp.Status)
			return
		}

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Error(err)
			return
		}
		defer resp.Body.Close()

		var data map[string]interface{}
		err = json.Unmarshal(bodyBytes, &data)
		if err != nil {
			log.Error(err)
			return
		}

		log.Info(data)

		// get the email from the customer object
		email := data["email"].(string)
		log.Info(email)

		// generate a random subscription id
		id, err := util.RandomString(8)
		if err != nil {
			log.Error(err)
		}

		// construct a subscription object
		subscription := model.Subscription{
			Id:          id,
			Email:       email,
			Issued:      time.Now(),
			LastUpdated: time.Now(),
		}

		errs := subscription.IsValid()
		if len(errs) != 0 {
			for _, err := range errs {
				log.WithFields(log.Fields{
					"err": err,
				}).Error("subscription validation error")
			}
			return
		}

		// save subscription to mongodb
		mongo.Serialize(subscription.Id, "id", "subscriptions", subscription)

	}()

	c.JSON(http.StatusOK, body)
}

func readSubscription(c *gin.Context) {
	id := c.Param("id")

	client, err := core.ReadSubscription(id)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to read client")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, client)
}

func updateSubscription(c *gin.Context) {
	var data model.Subscription
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

	client, err := core.UpdateSubscription(id, &data)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to update client")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, client)
}

func deleteSubscription(c *gin.Context) {
	id := c.Param("id")

	err := core.DeleteSubscription(id)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to remove client")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func readSubscriptions(c *gin.Context) {
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

	if user.Email == "" {
		log.Error("security alert: Email empty on authenticated token")
		c.AbortWithStatus(http.StatusForbidden)
	}

	subscriptions, err := core.ReadSubscriptions(user.Email)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("failed to list clients")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, subscriptions)
}
