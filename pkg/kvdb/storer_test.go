package kvdb

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetKey(t *testing.T) {
	storage := New_KV[string, int]()

	if err := storage.Set("one", 1); err != nil {
		t.Errorf("failed setting new key = %v", "one")
	}

	// test that the key is set correctly ..
	val, err := storage.Get("one")
	if err != nil {
		t.Errorf("error fetching val of key = %v", "one")
	}

	require.Equal(t, 1, val)
}

func TestSetExistingKey(t *testing.T) {
	storage := New_KV[string, int]()

	if err := storage.Set("one", 1); err != nil {
		t.Errorf("failed setting new key = %v", "one")
	}

	err := storage.Set("one", 11)
	t.Logf("error ➜ %v \n", err)
	require.Error(t, err)
}
