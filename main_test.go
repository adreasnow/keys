package main

import (
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/adreasnow/keychain-cli/keys"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v3"
	"github.com/zalando/go-keyring"
)

var ErrKeyring = errors.New("keyring error")

func TestObfuscate(t *testing.T) {
	for _, s := range []string{
		"testString1",
		"short1",
		"veryLongStringThatIsProbablyAComplexAPIKeyOrSomethingLikeThat...1",
	} {
		t.Run(s, func(t *testing.T) {
			key := obfuscate(s)
			assert.Len(t, key, len(s))
			assert.NotContains(t, key, "1")
		})
	}
}

func TestSet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		keyring.MockInit()

		buf := &bytes.Buffer{}
		cmd := &cli.Command{
			Writer: buf,
		}

		dict, err := keys.NewDict()
		require.NoError(t, err)

		key := "testKey"
		secret := "testSecret"

		err = set(key, secret, cmd)
		assert.NoError(t, err)

		keys, err := dict.GetAllKeys()
		assert.NoError(t, err)
		assert.Len(t, keys, 1)
		assert.Contains(t, keys, "testKey")

		assert.Contains(t, buf.String(), fmt.Sprintf("Set secret %s=%s", key, obfuscate(secret)))
	})

	t.Run("missing key", func(t *testing.T) {
		keyring.MockInit()

		cmd := &cli.Command{}

		err := set("", "testSecret", cmd)
		assert.ErrorIs(t, err, ErrMissingKey)
	})

	t.Run("missing secret", func(t *testing.T) {
		keyring.MockInit()

		cmd := &cli.Command{}

		err := set("testKey", "", cmd)
		assert.ErrorIs(t, err, ErrMissingSecret)
	})

	t.Run("keyring fail", func(t *testing.T) {
		keyring.MockInitWithError(ErrKeyring)

		cmd := &cli.Command{}

		err := set("testKey", "testSecret", cmd)
		assert.ErrorIs(t, err, ErrKeyring)
	})
}

func TestGet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		keyring.MockInit()

		buf := &bytes.Buffer{}
		cmd := &cli.Command{
			Writer: buf,
		}

		dict, err := keys.NewDict()
		require.NoError(t, err)

		key := "testKey"
		secret := "testSecret"
		err = dict.SetSecret(key, secret)
		require.NoError(t, err)

		err = get(key, cmd)
		assert.NoError(t, err)

		keys, err := dict.GetAllKeys()
		require.NoError(t, err)
		require.Len(t, keys, 1)

		assert.Contains(t, buf.String(), secret)
	})

	t.Run("missing key", func(t *testing.T) {
		keyring.MockInit()

		cmd := &cli.Command{}

		err := get("", cmd)
		assert.ErrorIs(t, err, ErrMissingKey)
	})

	t.Run("keyring fail", func(t *testing.T) {
		keyring.MockInitWithError(ErrKeyring)

		cmd := &cli.Command{}

		err := get("testKey", cmd)
		assert.ErrorIs(t, err, ErrKeyring)
	})
}

func TestDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		keyring.MockInit()

		buf := &bytes.Buffer{}
		cmd := &cli.Command{
			Writer: buf,
		}

		dict, err := keys.NewDict()
		require.NoError(t, err)

		key := "testKey"
		secret := "testSecret"
		err = dict.SetSecret(key, secret)
		require.NoError(t, err)

		err = delete(key, cmd)
		assert.NoError(t, err)

		dict, err = keys.NewDict()
		require.NoError(t, err)

		keys, err := dict.GetAllKeys()
		assert.NoError(t, err)
		assert.Len(t, keys, 0)

		assert.Contains(t, buf.String(), "Deleted secret "+key)
	})

	t.Run("missing key", func(t *testing.T) {
		keyring.MockInit()

		cmd := &cli.Command{}

		err := delete("", cmd)
		assert.ErrorIs(t, err, ErrMissingKey)
	})

	t.Run("keyring fail", func(t *testing.T) {
		keyring.MockInitWithError(ErrKeyring)

		cmd := &cli.Command{}

		err := delete("testKey", cmd)
		assert.ErrorIs(t, err, ErrKeyring)
	})
}

func TestList(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		keyring.MockInit()

		buf := &bytes.Buffer{}
		cmd := &cli.Command{
			Writer: buf,
		}

		dict, err := keys.NewDict()
		require.NoError(t, err)

		key1 := "testKey1"
		err = dict.SetSecret(key1, "testSecret1")
		require.NoError(t, err)

		key2 := "testKey2"
		err = dict.SetSecret(key2, "testSecret2")
		require.NoError(t, err)

		err = list(cmd)
		assert.NoError(t, err)

		assert.Contains(t, buf.String(), key1)
		assert.Contains(t, buf.String(), key2)
	})

	t.Run("keyring fail", func(t *testing.T) {
		keyring.MockInitWithError(ErrKeyring)

		cmd := &cli.Command{}

		err := get("testKey", cmd)
		assert.ErrorIs(t, err, ErrKeyring)
	})
}
