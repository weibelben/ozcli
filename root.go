package ozcli

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const defaultPath = "../"

var rootOzPath string

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

	pathHelp := "path to run the ozcli command in"
	rootCmd.PersistentFlags().StringVar(&rootOzPath, "path", defaultPath, pathHelp)
}

func setWorkingDirectory() {
	rootOzPath = os.Getenv("ROOT_DIR")
	if rootOzPath == "" {
		log.Fatal("ROOT_DIR not defined. Have you sourced a config?")
	}

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	log.Info(dir)
	log.Info("root oz path = " + rootOzPath)

	if err := os.Chdir(rootOzPath); err != nil {
		log.WithError(err).Fatalf(
			"could not set ozcli's working directory to %s", rootOzPath)
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
