package main

import (
	"log"

	"github.com/spf13/cobra"
)

const (
	CmdName           = "mmf"
	DefaultConfigFile = CmdName + ".hcl"
)

var (
	cfg          = NewConfig()
	cfgFile      string
	cfgReadError error
)

var rootCmd = &cobra.Command{
	Use:   CmdName,
	Short: "mmf gives command line access to MyMiniFactory.com",
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	// root level flags
	pflags := rootCmd.PersistentFlags()
	pflags.StringVarP(&cfgFile, "config-file", "f", DefaultConfigFile, "config file (HCL format)")

	// load config-file before cobra commands
	cobra.OnInitialize(func() {
		if cfgFile != "" {
			if err := cfg.ReadInFile(cfgFile); err != nil {
				cfgReadError = err
				log.Println(err)
			}
		}
	})
}
