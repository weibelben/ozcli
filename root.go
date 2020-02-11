package ozcli

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const defaultPath = "../"

var RootPath string

var rootCmd = &cobra.Command{
	Use:   "ozcli",
	Short: "OzCLI is your portal to Oz",
	Long: `A command driven interface for launching, modifying, and destroying
                Oz environments and related services.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func init() {
	cobra.OnInitialize(setWorkingDirectory)

	pathHelp := "path to run the hopcli command in"
	rootCmd.PersistentFlags().StringVar(&RootPath, "path", defaultPath, pathHelp)
}

func setWorkingDirectory() {
	RootPath = "."
	if err := os.Chdir(RootPath); err != nil {
		log.WithError(err).Fatalf(
			"could not set ozcli's working directory to %s", RootPath)
	}
}

// Execute is the main entry point into the Cobra cli. It parses the options
// provided and executes the desired functions based on commands and flags.
func Execute() {
	wd, _ := os.Getwd()
	defer func() { _ = os.Chdir(wd) }()
	if err := rootCmd.Execute(); err != nil {
		log.WithError(err).Fatal("Failed to run ozcli")
	}
}
