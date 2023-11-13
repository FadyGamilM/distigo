package config

import (
	"errors"
	"hash/fnv"
	"log"
	"os"

	"github.com/FadyGamilM/distigo/pkg/shards"
	"gopkg.in/yaml.v3"
)

type ShardsConfig struct {
	Shards []*shards.Shard
}

// read a yaml config file and parse its content to ShardConfig
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

func (shardsConfigs *ShardsConfig) CheckShardExistence(shardName string) (shardIdx int, err error) {
	for _, shard := range shardsConfigs.Shards {
		if shard.Name == shardName {
			log.Printf("shard with name = %v is found, the index of this shard = %v \n", shardName, shard.Idx)
			shardIdx := shard.Idx
			return shardIdx, nil
		}
	}

	return -1, errors.New("could not find a shard with name = " + shardName)
}

// hash the key to know the index of the appropriate shard to store this key-val pair
func (sc *ShardsConfig) DistributeKeyOnShards(key string) (int, error) {
	// define a hash
	hash := fnv.New64()
	// Hash(Key)
	hash.Write([]byte(key))
	// get the reminder
	return int(hash.Sum64() % uint64(len(sc.Shards))), nil
}
