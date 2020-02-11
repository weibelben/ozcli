package ozcli

import (
	"fmt"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	rootOzcliPath = "../"
)

func init() {
	rootCmd.AddCommand(k8sTunnelCmd)
	
	k8sTunnelCmd.Flags().StringVar(&rootOzcliPath, "path", rootOzcliPath,
		"path to run the tunnel command in")
}

var k8sTunnelCmd = &cobra.Command{
	Use:   "k8s-tunnel",
	Short: "Create a tunnel to k8s",
	Long: `Create a tunnel to the k8s api of the config that is
			currently sourced.`,
	Run: func(cmd *cobra.Command, args []string) {
		createK8sTunnel()
	},
}

func createK8sTunnel() {
	dir, err := os.Getwd()
	log.Info(dir)
	
	log.Info("Creating tunnel to k8s api...")
	makeTunnelCmd := exec.Command(fmt.Sprintf("bash %s/infrastructure/dev/k8s/make-tunnel.sh k8sapi", rootOzcliPath))
	err = makeTunnelCmd.Start()
	if err != nil {
		log.WithError(err).Error("Failed to run make-tunnel script.")
		return
	}
	err = makeTunnelCmd.Wait()
	if err != nil {
		log.WithError(err).Error("Failed to create tunnel to k8s api.")
	}
}
