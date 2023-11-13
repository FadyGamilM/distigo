package shards

type Shard struct {
	Name string
	Idx  int
	Host string
}

type Shards []*Shard
