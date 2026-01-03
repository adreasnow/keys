package keys

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"
)

func TestLoadKeys(t *testing.T) {
	keyring.MockInit()
	err := keyring.Set(keysSecret, user, `{"key":true}`)
	require.NoError(t, err)

	d := &Dict{keys: map[string]bool{}}
	err = d.loadKeys()
	require.NoError(t, err)

	require.Len(t, d.keys, 1)
	assert.Contains(t, d.keys, "key")
	assert.True(t, d.keys["key"])
}

func TestSaveKeys(t *testing.T) {
	keyring.MockInit()
	err := keyring.Set(keysSecret, user, `{}`)
	require.NoError(t, err)

	d := &Dict{keys: map[string]bool{
		"key": true,
	}}
	err = d.saveKeys()
	require.NoError(t, err)

	require.Len(t, d.keys, 1)
	assert.Contains(t, d.keys, "key")
	assert.True(t, d.keys["key"])
}

func TestNewDict(t *testing.T) {
	t.Run("uninitialised", func(t *testing.T) {
		keyring.MockInit()

		d := NewDict()
		require.NotNil(t, d)
		require.Len(t, d.keys, 0)
	})

	t.Run("uninitialised", func(t *testing.T) {
		keyring.MockInit()
		err := keyring.Set(keysSecret, user, `{"key":true}`)
		require.NoError(t, err)

		d := NewDict()
		require.NotNil(t, d)
		require.Len(t, d.keys, 1)
		assert.Contains(t, d.keys, "key")
		assert.True(t, d.keys["key"])
	})
}

func TestGetAllKeys(t *testing.T) {
	keyring.MockInit()
	err := keyring.Set(keysSecret, user, `{"key1":true,"key2":true,"key3":true}`)
	require.NoError(t, err)

	d := &Dict{keys: map[string]bool{}}
	keys, err := d.GetAllKeys()
	require.NoError(t, err)

	require.Len(t, keys, 3)
	assert.Contains(t, keys, "key1")
	assert.Contains(t, keys, "key2")
	assert.Contains(t, keys, "key3")
}

func TestAddKey(t *testing.T) {
	keyring.MockInit()
	err := keyring.Set(keysSecret, user, `{"key1":true}`)
	require.NoError(t, err)

	d := &Dict{keys: map[string]bool{}}
	keys, err := d.GetAllKeys()
	require.NoError(t, err)
	require.Len(t, keys, 1)

	err = d.AddKey("key2")
	require.NoError(t, err)

	keys, err = d.GetAllKeys()
	require.NoError(t, err)
	require.Len(t, keys, 2)
	assert.Contains(t, keys, "key1")
	assert.Contains(t, keys, "key2")
}

func TestDeleteKey(t *testing.T) {
	keyring.MockInit()
	err := keyring.Set(keysSecret, user, `{"key1":true,"key2":true}`)
	require.NoError(t, err)

	d := &Dict{keys: map[string]bool{}}
	keys, err := d.GetAllKeys()
	require.NoError(t, err)
	require.Len(t, keys, 2)

	err = d.DeleteKey("key1")
	require.NoError(t, err)

	keys, err = d.GetAllKeys()
	require.NoError(t, err)
	require.Len(t, keys, 1)
	assert.Contains(t, keys, "key2")
}
