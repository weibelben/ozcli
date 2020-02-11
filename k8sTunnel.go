package ozcli

import (
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(k8sTunnelCmd)
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
	log.Info("Creating tunnel to k8s api...")
	path := "./infrastructure/dev/k8s/make-tunnel.sh"

	runBashScript(path, "k8sapi")
}

func runBashScript(path string, arg string) {
	makeTunnelCmd := exec.Command(path, arg)
	makeTunnelCmd.Stdout = os.Stdout
	makeTunnelCmd.Stderr = os.Stderr

	err := makeTunnelCmd.Start()
	if err != nil {
		log.WithError(err).Errorf("Failed to run %s.", path)
		return
	}

	err = makeTunnelCmd.Wait()
	if err != nil {
		log.WithError(err).Errorf("An error occurred while running %s.", path)
	}
}
