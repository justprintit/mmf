package main

import (
	"os"

	"github.com/sethvargo/go-password/password"
	"github.com/spf13/cobra"

	"github.com/justprintit/mmf"
)

func newPasswordGenerator() (*password.Generator, error) {
	input := &password.GeneratorInput{
		Symbols: "!@#$%^()",
	}

	return password.NewGenerator(input)
}

func generateClientId(gen *password.Generator) (string, error) {
	return gen.Generate(16, 4, 0, false, true)
}

func GenerateClientId() (string, error) {
	gen, err := newPasswordGenerator()
	if err != nil {
		return "", err
	}

	s, err := generateClientId(gen)
	if err != nil {
		return "", err
	}

	return s, nil
}

func generateClientSecret(gen *password.Generator) (string, error) {
	return gen.Generate(32, 8, 4, false, true)
}

func GenerateClientSecret() (string, error) {
	gen, err := newPasswordGenerator()
	if err != nil {
		return "", err
	}

	s, err := generateClientSecret(gen)
	if err != nil {
		return "", err
	}

	return s, nil
}

func GenerateClientPair() (string, string, error) {
	gen, err := newPasswordGenerator()
	if err != nil {
		return "", "", err
	}

	id, err := generateClientId(gen)
	if err != nil {
		return "", "", err
	}

	secret, err := generateClientSecret(gen)
	if err != nil {
		return "", "", err
	}

	return id, secret, nil
}

var generateCmd = &cobra.Command{
	Use: "generate outputs a new config file",
	RunE: func(cmd *cobra.Command, args []string) error {

		id, secret, err := GenerateClientPair()
		if err != nil {
			return err
		}

		cfg := mmf.Config{
			ClientId:     id,
			ClientSecret: secret,
		}

		if buf, err := mmf.MarshalConfig(&cfg); err != nil {
			return err
		} else if _, err := os.Stdout.Write(buf); err != nil {
			return err
		} else {
			return nil
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
