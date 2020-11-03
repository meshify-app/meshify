package apiv1

import (
	user "github.com/alan-grapid/meshify/api/v1/Users"
	"github.com/alan-grapid/meshify/api/v1/auth"
	"github.com/alan-grapid/meshify/api/v1/client"
	"github.com/alan-grapid/meshify/api/v1/mesh"
	"github.com/alan-grapid/meshify/api/v1/server"
	"github.com/gin-gonic/gin"
)

// ApplyRoutes apply routes to gin router
func ApplyRoutes(r *gin.RouterGroup, private bool) {
	v1 := r.Group("/v1.0")
	{
		if private {
			client.ApplyRoutes(v1)
			server.ApplyRoutes(v1)
			user.ApplyRoutes(v1)
			mesh.ApplyRoutes(v1)
		} else {
			auth.ApplyRoutes(v1)

		}
	}
}
