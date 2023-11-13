package config

import (
	"log"
	"os"

	"github.com/FadyGamilM/distigo/pkg/shards"
	"gopkg.in/yaml.v3"
)

type ShardsConfig struct {
	Shards []*shards.Shard
}

func ParseShardsConfig(configFilePath string) (*ShardsConfig, error) {
	configFile, err := os.ReadFile(configFilePath)
	if err != nil {
		log.Printf("error trying to read shards config file ➜ %v \n", err)
		return nil, err
	}

	shards := &ShardsConfig{}
	err = yaml.Unmarshal(configFile, shards)
	if err != nil {
		log.Printf("error trying to unmarshal the shatds config file into shards type ➜ %v \n", err)
		return nil, err
	}

	for _, shard := range shards.Shards {
		log.Println(shard)
	}

	return shards, nil
}
