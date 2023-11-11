package shards

import (
	"errors"
	"log"
)

type Shard struct {
	Name string
	Idx  int
}

type Shards []*Shard

func (shards *Shards) CheckShardExistence(shardName string) (shardIdx int, err error) {
	for _, shard := range *shards {
		if shard.Name == shardName {
			log.Printf("shard with name = %v is found, the index of this shard = %v \n", shardName, shard.Idx)
			shardIdx := shard.Idx
			return shardIdx, nil
		}
	}

	return -1, errors.New("could not find a shard with name = " + shardName)
}
