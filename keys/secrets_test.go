package keys

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"
)

func TestSetSecret(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		key := "testKey"
		secret := "testSecret"

		keyring.MockInit()

		err := keyring.Set(keysSecret, user, `{}`)
		require.NoError(t, err)

		err = keyring.Set(key, user, secret)
		require.NoError(t, err)

		d := &Dict{keys: map[string]bool{}}
		err = d.SetSecret(key, secret)
		require.NoError(t, err)

		retrievedSecret, err := keyring.Get(key, user)
		require.NoError(t, err)
		assert.Equal(t, secret, retrievedSecret)

		keys, err := keyring.Get(keysSecret, user)
		require.NoError(t, err)
		assert.Equal(t, `{"`+key+`":true}`, keys)

		require.Len(t, d.keys, 1)
	})

	t.Run("keyring fail", func(t *testing.T) {
		keyring.MockInitWithError(ErrKeyring)

		d := &Dict{keys: map[string]bool{}}
		err := d.SetSecret("test", "test")
		assert.ErrorIs(t, err, ErrKeyring)
	})
}

func TestGetSecret(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		key := "testKey"
		secret := "testSecret"

		keyring.MockInit()
		err := keyring.Set(key, user, secret)
		require.NoError(t, err)

		d := &Dict{keys: map[string]bool{}}
		retrievedSecret, err := d.GetSecret(key)
		require.NoError(t, err)

		assert.Equal(t, secret, retrievedSecret)
	})

	t.Run("keyring fail", func(t *testing.T) {
		keyring.MockInitWithError(ErrKeyring)

		d := &Dict{keys: map[string]bool{}}
		_, err := d.GetSecret("test")
		assert.ErrorIs(t, err, ErrKeyring)
	})
}

func TestDeleteSecret(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		key := "testKey"
		secret := "testSecret"

		keyring.MockInit()

		err := keyring.Set(keysSecret, user, `{"`+key+`":true}`)
		require.NoError(t, err)

		err = keyring.Set(key, user, secret)
		require.NoError(t, err)

		d := &Dict{keys: map[string]bool{}}
		err = d.DeleteSecret(key)
		require.NoError(t, err)

		_, err = keyring.Get(key, user)
		require.ErrorIs(t, err, keyring.ErrNotFound)

		keys, err := keyring.Get(keysSecret, user)
		require.NoError(t, err)
		assert.Equal(t, `{}`, keys)

		assert.Len(t, d.keys, 0)
	})

	t.Run("keyring fail", func(t *testing.T) {
		keyring.MockInitWithError(ErrKeyring)

		d := &Dict{keys: map[string]bool{}}
		err := d.DeleteSecret("test")
		assert.ErrorIs(t, err, ErrKeyring)
	})
}
