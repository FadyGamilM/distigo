package transport

import "github.com/gin-gonic/gin"

func ReadQueryParams(c *gin.Context, key string) string {
	val := c.Request.URL.Query().Get(key)
	return val
}
