package transport

import (
	"fmt"
	"log"
	"net/http"

	"github.com/FadyGamilM/distigo/business"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	StorageService business.StorageService
}

func NewHandler(ss business.StorageService) *Handler {
	return &Handler{
		StorageService: ss,
	}
}

func (h *Handler) HandleGet(c *gin.Context) {

	// read query params
	key := ReadQueryParams(c, "key")
	val, err := h.StorageService.Get([]byte(key))
	if err != nil {
		log.Printf("error trying to fetch val of key = %v ➜ &v \n", key, err)
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"error": fmt.Sprintf("no value with key = %v", key),
			},
		)
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"response": val,
		},
	)
}

func (h *Handler) HandlePost(c *gin.Context) {
	// read query params
	key := ReadQueryParams(c, "key")
	value := ReadQueryParams(c, "val")

	err := h.StorageService.Set([]byte(key), []byte(value))
	if err != nil {
		log.Printf("error trying to set new (key, val) pair in database, key = %v val = %v ➜ &v \n", key, value, err)
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"error": fmt.Sprintf("couldn't set new (key, val) pair in database"),
			},
		)
	}

	c.JSON(
		http.StatusCreated,
		gin.H{
			"response": "post is called",
		},
	)
}
