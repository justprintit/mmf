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

	return generateClientId(gen)
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "generate outputs a new config file",
	RunE: func(cmd *cobra.Command, args []string) error {

		// ClientId
		id, err := GenerateClientId()
		if err != nil {
			return err
		}

		// Config
		cfg := &Config{
			Auth: mmf.Config{
				ClientId: id,
			},
		}

		_, err = cfg.WriteTo(os.Stdout)
		return err
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
