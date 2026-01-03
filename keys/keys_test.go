package keys

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"
)

var ErrKeyring = errors.New("keyring error")

func TestNewDict(t *testing.T) {
	t.Run("uninitialised", func(t *testing.T) {
		keyring.MockInit()

		d, err := NewDict()
		require.NoError(t, err)

		require.NotNil(t, d)
		require.Len(t, d.keys, 0)
	})

	t.Run("uninitialised", func(t *testing.T) {
		keyring.MockInit()
		err := keyring.Set(keysSecret, user, `{"key":true}`)
		require.NoError(t, err)

		d, err := NewDict()
		require.NoError(t, err)
		require.NotNil(t, d)

		require.Len(t, d.keys, 1)
		assert.Contains(t, d.keys, "key")
		assert.True(t, d.keys["key"])
	})

	t.Run("keyring fail", func(t *testing.T) {
		keyring.MockInitWithError(ErrKeyring)

		_, err := NewDict()
		assert.ErrorIs(t, err, ErrKeyring)
	})
}

func TestLoadKeys(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		keyring.MockInit()
		err := keyring.Set(keysSecret, user, `{"key":true}`)
		require.NoError(t, err)

		d := &Dict{keys: map[string]bool{}}
		err = d.loadKeys()
		require.NoError(t, err)

		require.Len(t, d.keys, 1)
		assert.Contains(t, d.keys, "key")
		assert.True(t, d.keys["key"])
	})

	t.Run("unmarshal error", func(t *testing.T) {
		keyring.MockInit()
		err := keyring.Set(keysSecret, user, `{`)
		require.NoError(t, err)

		d := &Dict{keys: map[string]bool{}}
		err = d.loadKeys()

		var syntaxError *json.SyntaxError
		require.ErrorAs(t, err, &syntaxError)
	})

	t.Run("keyring fail", func(t *testing.T) {
		keyring.MockInitWithError(ErrKeyring)

		d := &Dict{keys: map[string]bool{}}
		err := d.loadKeys()
		assert.ErrorIs(t, err, ErrKeyring)
	})
}

func TestSaveKeys(t *testing.T) {
	t.Run("success", func(t *testing.T) {
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
	})

	t.Run("keyring fail", func(t *testing.T) {
		keyring.MockInitWithError(ErrKeyring)

		d := &Dict{keys: map[string]bool{}}
		err := d.saveKeys()
		assert.ErrorIs(t, err, ErrKeyring)
	})
}

func TestGetAllKeys(t *testing.T) {
	t.Run("success", func(t *testing.T) {
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
	})

	t.Run("keyring fail", func(t *testing.T) {
		keyring.MockInitWithError(ErrKeyring)

		d := &Dict{keys: map[string]bool{}}
		_, err := d.GetAllKeys()
		assert.ErrorIs(t, err, ErrKeyring)
	})
}

func TestAddKey(t *testing.T) {
	t.Run("success", func(t *testing.T) {
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
	})

	t.Run("keyring fail", func(t *testing.T) {
		keyring.MockInitWithError(ErrKeyring)

		d := &Dict{keys: map[string]bool{}}
		err := d.AddKey("test")
		assert.ErrorIs(t, err, ErrKeyring)
	})
}

func TestDeleteKey(t *testing.T) {
	t.Run("success", func(t *testing.T) {
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
	})

	t.Run("keyring fail", func(t *testing.T) {
		keyring.MockInitWithError(ErrKeyring)

		d := &Dict{keys: map[string]bool{}}
		err := d.DeleteKey("test")
		assert.ErrorIs(t, err, ErrKeyring)
	})
}
