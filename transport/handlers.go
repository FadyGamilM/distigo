package transport

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleGet(c *gin.Context) {

	// read query params
	key := ReadQueryParams(c, "key")

	c.JSON(
		http.StatusOK,
		gin.H{
			"response": "get is called",
		},
	)
}

func HandlePost(c *gin.Context) {
	// read query params
	key := ReadQueryParams(c, "key")
	value := ReadQueryParams("val")

	c.JSON(
		http.StatusCreated,
		gin.H{
			"response": "post is called",
		},
	)
}
