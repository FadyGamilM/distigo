package transport

import "github.com/gin-gonic/gin"

type DistigoRouter struct {
	r       *gin.Engine
	handler *Handler
}

func NewDistigoRouter(router *gin.Engine, handler *Handler) *DistigoRouter {
	return &DistigoRouter{
		r:       router,
		handler: handler,
	}
}

func HttpRouter() *gin.Engine {
	return gin.Default()
}

func (dr *DistigoRouter) SetupEndpoints() {
	routerGroup := dr.r.Group("/api/distigo")
	routerGroup.GET("/get", dr.handler.HandleGet)
	routerGroup.POST("/set", dr.handler.HandlePost)
	routerGroup.DELETE("/purge", dr.handler.HandlePurge)
	routerGroup.GET("/logs", dr.handler.GetAllShardKeys)
}
