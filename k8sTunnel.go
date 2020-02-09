package ozcli

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(k8sTunnelCmd)
}

var k8sTunnelCmd = &cobra.Command{
	Use:   "k8s-tunnel",
	Short: "Create a tunnel to k8s",
	Long: `Create a tunnel to the k8s api of the config that is currently
			sourced.`,
	Run: func(cmd *cobra.Command, args []string) {
		createK8sTunnel()
	},
}

func createK8sTunnel() {
	log.Info("Creating tunnel to k8s api")
}
