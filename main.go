package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/adreasnow/keychain-cli/keys"
	"github.com/urfave/cli/v3"
)

var (
	ErrMissingKey    = errors.New("missing key")
	ErrMissingSecret = errors.New("missing secret")
)

func main() {
	var key string
	var secret string

	cmd := &cli.Command{
		Name:                  "keychain-cli",
		EnableShellCompletion: true,
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List all secrets",
				Action: func(_ context.Context, cmd *cli.Command) error {
					return list(cmd)
				},
			},
			{
				Name:  "get",
				Usage: "Get a specific secret by key",
				Arguments: []cli.Argument{
					&cli.StringArg{Name: "key", Destination: &key},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					return get(key, cmd)
				},
				ShellComplete: func(_ context.Context, cmd *cli.Command) { completion(cmd) },
			},
			{
				Name:  "set",
				Usage: "Create/update a secret",
				Arguments: []cli.Argument{
					&cli.StringArg{Name: "key", Destination: &key},
					&cli.StringArg{Name: "secret", Destination: &secret},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					return set(key, secret, cmd)
				},
			},
			{
				Name:  "delete",
				Usage: "Delete a secret",
				Arguments: []cli.Argument{
					&cli.StringArg{Name: "key", Destination: &key},
				},
				Action: func(_ context.Context, cmd *cli.Command) error {
					return delete(key, cmd)
				},
				ShellComplete: func(_ context.Context, cmd *cli.Command) { completion(cmd) },
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintln(cmd.ErrWriter, err) //nolint:errcheck
	}
}

func obfuscate(secret string) string {
	out := strings.Builder{}
	for range secret {
		out.WriteString("*")
	}
	return out.String()
}

func set(key, secret string, cmd *cli.Command) error {
	dict, err := keys.NewDict()
	if err != nil {
		return err
	}

	if key == "" {
		return ErrMissingKey
	}

	if secret == "" {
		return ErrMissingSecret
	}

	err = dict.SetSecret(key, secret)
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.Writer, "Set secret %s=%s\n", key, obfuscate(secret)) //nolint:errcheck
	return nil
}

func list(cmd *cli.Command) error {
	dict, err := keys.NewDict()
	if err != nil {
		return err
	}

	secrets, err := dict.GetAllKeys()
	if err != nil {
		return err
	}

	for _, secret := range secrets {
		fmt.Fprintln(cmd.Writer, secret) //nolint:errcheck
	}

	return nil
}

func get(key string, cmd *cli.Command) error {
	dict, err := keys.NewDict()
	if err != nil {
		return err
	}

	if key == "" {
		return ErrMissingKey
	}

	secret, err := dict.GetSecret(key)
	if err != nil {
		return err
	}

	fmt.Fprintln(cmd.Writer, secret) //nolint:errcheck
	return nil
}

func delete(key string, cmd *cli.Command) error {
	dict, err := keys.NewDict()
	if err != nil {
		return err
	}

	if key == "" {
		return ErrMissingKey
	}

	err = dict.DeleteSecret(key)
	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.Writer, "Deleted secret %s\n", key) //nolint:errcheck
	return nil
}

func completion(cmd *cli.Command) {
	dict, err := keys.NewDict()
	if err != nil {
		fmt.Fprintf(cmd.ErrWriter, "failed to get keys: %v\n", err) //nolint:errcheck
		os.Exit(1)
	}

	keys, err := dict.GetAllKeys()
	if err != nil {
		fmt.Fprintf(cmd.ErrWriter, "failed to get keys: %v\n", err) //nolint:errcheck
		os.Exit(1)
	}

	args := cmd.Args().Slice()
	for _, key := range keys {
		if slices.Contains(args, key) {
			return
		}
	}

	for _, key := range keys {
		fmt.Println(key)
	}
}
