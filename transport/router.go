package transport

import "github.com/gin-gonic/gin"

type DistigoRouter struct {
	Handler *gin.Engine
}

func HttpRouter() *gin.Engine {
	return gin.Default()
}

func (dr *DistigoRouter) SetupEndpoints() {
	dr.Handler.Group("/api/distigo")
	dr.Handler.GET("/get", HandleGet)
	dr.Handler.POST("/set", HandlePost)
}
