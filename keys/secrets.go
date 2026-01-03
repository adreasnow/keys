package keys

import (
	"fmt"

	"github.com/zalando/go-keyring"
)

func (d Dict) GetSecret(key string) (secret string, err error) {
	secret, err = keyring.Get(key, user)
	if err != nil {
		err = fmt.Errorf("failed to get secret: %w", err)
	}

	return
}

func (d *Dict) DeleteSecret(key string) (err error) {
	err = keyring.Delete(key, user)
	if err != nil {
		err = fmt.Errorf("failed to delete secret: %w", err)
		return
	}

	err = d.DeleteKey(key)
	if err != nil {
		err = fmt.Errorf("failed to delete key: %w", err)
	}

	return
}

func (d *Dict) SetSecret(key, secret string) (err error) {
	err = keyring.Set(key, user, secret)
	if err != nil {
		err = fmt.Errorf("failed to set secret: %w", err)
		return
	}

	err = d.AddKey(key)
	if err != nil {
		err = fmt.Errorf("failed to add key: %w", err)
	}

	return
}
