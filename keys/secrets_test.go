package keys

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zalando/go-keyring"
)

func TestSetSecret(t *testing.T) {
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
}

func TestGetSecret(t *testing.T) {
	key := "testKey"
	secret := "testSecret"

	keyring.MockInit()
	err := keyring.Set(key, user, secret)
	require.NoError(t, err)

	d := &Dict{keys: map[string]bool{}}
	retrievedSecret, err := d.GetSecret(key)
	require.NoError(t, err)

	assert.Equal(t, secret, retrievedSecret)
}

func TestDeleteSecret(t *testing.T) {
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
}
