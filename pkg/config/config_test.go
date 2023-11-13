package config

import (
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseShardsConfig(t *testing.T) {
	configs, err := ParseShardsConfig("../../shards.yaml")
	require.NoError(t, err)
	log.Println(configs.Shards[0].Name)
	require.Equal(t, 3, len(configs.Shards))
}
