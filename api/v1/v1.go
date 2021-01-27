package apiv1

import (
	"github.com/gin-gonic/gin"
	user "github.com/meshify-app/meshify/api/v1/Users"
	"github.com/meshify-app/meshify/api/v1/account"
	"github.com/meshify-app/meshify/api/v1/auth"
	host "github.com/meshify-app/meshify/api/v1/host"
	"github.com/meshify-app/meshify/api/v1/mesh"
	"github.com/meshify-app/meshify/api/v1/server"
)

// ApplyRoutes apply routes to gin router
func ApplyRoutes(r *gin.RouterGroup, private bool) {
	v1 := r.Group("/v1.0")
	{
		if private {
			account.ApplyRoutes(v1)
			account.ApplyOrgRoutes(v1)
			host.ApplyRoutes(v1)
			server.ApplyRoutes(v1)
			user.ApplyRoutes(v1)
			mesh.ApplyRoutes(v1)
		} else {
			auth.ApplyRoutes(v1)

		}
	}
}
