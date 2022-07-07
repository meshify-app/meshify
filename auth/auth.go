package auth

import (
	"fmt"
	"os"

	"github.com/meshify-app/meshify/auth/fake"
	"github.com/meshify-app/meshify/auth/github"
	"github.com/meshify-app/meshify/auth/microsoft"
	"github.com/meshify-app/meshify/auth/oauth2oidc"
	model "github.com/meshify-app/meshify/model"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

// Auth interface to implement as auth provider
type Auth interface {
	Setup() error
	CodeUrl(state string) string
	Exchange(code string) (*oauth2.Token, error)
	UserInfo(oauth2Token *oauth2.Token) (*model.User, error)
}

// GetAuthProvider  get an instance of auth provider based on config
func GetAuthProvider() (Auth, error) {
	var oauth2Client Auth
	var err error

	switch os.Getenv("OAUTH2_PROVIDER_NAME") {
	case "fake":
		log.Warn("Oauth is set to fake, no actual authentication will be performed")
		oauth2Client = &fake.Fake{}

	case "oauth2oidc":
		log.Warn("Oauth is set to oauth2oidc, must be RFC implementation on server side")
		oauth2Client = &oauth2oidc.Oauth2idc{}

	case "microsoft":
		log.Warn("Oauth is set to Microsoft")
		oauth2Client = &microsoft.Oauth2Msft{}

	case "github":
		log.Warn("Oauth is set to github, no openid will be used")
		oauth2Client = &github.Github{}

	case "google":
		return nil, fmt.Errorf("auth provider name %s not yet implemented", os.Getenv("OAUTH2_PROVIDER_NAME"))
	default:
		return nil, fmt.Errorf("auth provider name %s unknown", os.Getenv("OAUTH2_PROVIDER_NAME"))
	}

	err = oauth2Client.Setup()

	return oauth2Client, err
}
