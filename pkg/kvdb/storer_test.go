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

func TestUpdateKey(t *testing.T) {
	storage := New_KV[string, int]()

	if err := storage.Set("one", 1); err != nil {
		t.Errorf("failed setting new key = %v", "one")
	}

	val, err := storage.Get("one")
	if err != nil {
		t.Logf("error fetching val of key = %v", "one")
	}
	require.NoError(t, err)
	// before update
	require.Equal(t, 1, val)

	err = storage.Update("one", 11)
	if err != nil {
		t.Logf("failed updating key = %v with value = %v", "one", 11)
	}
	require.NoError(t, err)

	// test that the key is set correctly ..
	val, err = storage.Get("one")
	if err != nil {
		t.Errorf("error fetching val of key = %v", "one")
	}
	require.NoError(t, err)
	require.Equal(t, 11, val)
}

func TestDeleteKey(t *testing.T) {
	storage := New_KV[string, int]()

	if err := storage.Set("one", 1); err != nil {
		t.Errorf("failed setting new key = %v", "one")
	}

	val, err := storage.Get("one")
	if err != nil {
		t.Logf("error fetching val of key = %v", "one")
	}
	require.NoError(t, err)
	// before update
	require.Equal(t, 1, val)

	err = storage.Delete("one")
	if err != nil {
		t.Logf("error deleting val of key = %v", "one")
	}
	require.NoError(t, err)

	val, err = storage.Get("one")
	if err != nil {
		t.Logf("error fetching val of key = %v", "one")
	}
	t.Logf("the value is : %v", val) // => from here i knew that the go system will set the value as the zero value of this type
	require.Error(t, err)
}
