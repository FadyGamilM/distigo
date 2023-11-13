package transport

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/FadyGamilM/distigo/business"
	"github.com/FadyGamilM/distigo/pkg/config"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	StorageService business.StorageService
	ShardsConfig   *config.ShardsConfig
	// the shard that our server host (one server hosts one shard)
	HostedShardIndex int
	// the ShardsAddrs is a map between the shard index and the shard host
	ShardsAddrs map[int]string
}

func NewHandler(ss business.StorageService, shardsConfig *config.ShardsConfig, currentShardIndex int) *Handler {
	shardsAddrs := make(map[int]string, 0)
	for _, shard := range shardsConfig.Shards {
		shardsAddrs[shard.Idx] = shard.Host
	}
	return &Handler{
		StorageService:   ss,
		ShardsConfig:     shardsConfig,
		HostedShardIndex: currentShardIndex,
		ShardsAddrs:      shardsAddrs,
	}
}

func (h *Handler) HandleGet(c *gin.Context) {

	// read query params
	key := ReadQueryParams(c, "key")

	// get the shard that this key is distributed on
	// is doesn't have to be the hosted shard on the current running server which received the request, for example the current server which received the request is 8080 which hosts shard 0 but the key is distributed on shard 1, so we should route the request to the appropriate shard host
	distributedShardIdx, err := h.ShardsConfig.DistributeKeyOnShards(key)

	log.Printf("the shard that this request is distributed on is ➜ %v , the shard that is hosted on the current runnign server is ➜ %v \n", h.ShardsConfig.Shards[distributedShardIdx].Name, h.ShardsConfig.Shards[h.HostedShardIndex].Name)

	if err != nil {
		log.Println("error finding the appropriate shard for getting the key,val from")
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": fmt.Sprintf("error distributing the key = %v to a shard", key),
			},
		)
		return
	}
	// if the key is distributed on a shard that is not the one we host over here ..
	if distributedShardIdx != h.HostedShardIndex {
		// send the request to the approriate server
		resp, err := http.Get("http://" + h.ShardsAddrs[distributedShardIdx] + c.Request.RequestURI)
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				gin.H{
					"error": fmt.Sprintf("error trying to forward the request to appropriate shard host ➜ &v \n", err),
				},
			)
			return
		}
		// close the body
		defer resp.Body.Close()
		// read the response
		// Read the response body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
			return
		}
		// copy the response
		c.JSON(
			http.StatusOK,
			gin.H{
				"from-shard": h.ShardsConfig.Shards[distributedShardIdx].Name,
				"from-host":  h.ShardsAddrs[distributedShardIdx],
				"response":   body,
			},
		)
		return
	}

	val, err := h.StorageService.Get([]byte(key))
	if err != nil {
		log.Printf("error trying to fetch val of key = %v ➜ &v \n", key, err)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": fmt.Sprintf("no value with key = %v", key),
			},
		)
	}

	if val == nil {
		log.Printf("key = %v is not stored before, so we couldn't find any value associated with it \n", key)
		c.JSON(
			http.StatusNotFound,
			gin.H{
				"error": fmt.Sprintf("key = %v is not stored before", key),
			},
		)
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"value":      string(val),
			"shard":      distributedShardIdx,
			"shard-host": h.ShardsAddrs[h.HostedShardIndex],
		},
	)
}

func (h *Handler) HandlePost(c *gin.Context) {
	// read query params
	key := ReadQueryParams(c, "key")
	value := ReadQueryParams(c, "val")

	// distribute the key (just find the hash for response right now)
	distributedShardIdx, err := h.ShardsConfig.DistributeKeyOnShards(key)

	if err != nil {
		log.Println("error finding the appropriate shard for setting the key,val into")
		c.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": fmt.Sprintf("error distributing the key = %v to a shard", key),
			},
		)
		return
	}

	if distributedShardIdx == h.HostedShardIndex {
		err = h.StorageService.Set([]byte(key), []byte(value))
		if err != nil {
			log.Printf("error trying to set new (key, val) pair in database, key = %v val = %v ➜ &v \n", key, value, err)
			c.JSON(
				http.StatusNotFound,
				gin.H{
					"error": fmt.Sprintf("couldn't set new (key, val) pair in database"),
				},
			)
			return
		}
	} else {
		resp, err := http.Post("http://"+h.ShardsAddrs[distributedShardIdx]+c.Request.RequestURI, "application/json", bytes.NewBuffer([]byte(value)))
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				gin.H{
					"error": fmt.Sprintf("error trying to forward the request to appropriate shard host ➜ &v \n", err),
				},
			)
			return
		}

		// close the body
		defer resp.Body.Close()
		// read the response
		// Read the response body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
			return
		}

		c.JSON(
			http.StatusCreated,
			gin.H{
				"response":   body,
				"from-shard": h.ShardsConfig.Shards[distributedShardIdx].Name,
				"from-host":  h.ShardsAddrs[distributedShardIdx],
			},
		)
		return
	}
}
