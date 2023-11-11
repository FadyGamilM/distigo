package config

import (
	"io/ioutil"
	"log"

	"github.com/FadyGamilM/distigo/pkg/shards"
	"gopkg.in/yaml.v3"
)

func ParseShardsConfig(configFilePath string) ([]*shards.Shard, error) {
	configFile, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Printf("error trying to read shatds config file ➜ %v \n", err)
		return nil, err
	}

	var shards []*shards.Shard
	err = yaml.Unmarshal(configFile, shards)
	if err != nil {
		log.Printf("error trying to unmarshal the shatds config file into shards type ➜ %v \n", err)
		return nil, err
	}

	for _, shard := range shards {
		log.Println(shard)
	}

	return shards, nil
}
