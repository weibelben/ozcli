package ozcli

import (
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(k8sTunnelCmd)
	rootCmd.AddCommand(socksProxyCmd)
	rootCmd.AddCommand(portForwardVaultCmd)
	rootCmd.AddCommand(portForwardCeleriumCmd)
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

var socksProxyCmd = &cobra.Command{
	Use:   "socks-proxy",
	Short: "Create a SOCKS proxy",
	Long:  "Create a SOCKS proxy.",
	Run: func(cmd *cobra.Command, args []string) {
		createSocksProxy()
	},
}

var portForwardVaultCmd = &cobra.Command{
	Use:   "port-forward-vault",
	Short: "Port forward to vault",
	Long:  "Port forward to vault allowing you to access its web UI.",
	Run: func(cmd *cobra.Command, args []string) {
		portForwardVault()
	},
}

var portForwardCeleriumCmd = &cobra.Command{
	Use:   "port-forward-celerium",
	Short: "Port forward to celerium",
	Long:  "Port forward to celerium allowing you to access its web UI and the db.",
	Run: func(cmd *cobra.Command, args []string) {
		portForwardCelerium()
	},
}

func createK8sTunnel() {
	log.Info("Creating tunnel to k8s api...")
	path := "./infrastructure/dev/k8s/make-tunnel.sh"

	runBashScript(path, "k8sapi")
}

func createSocksProxy() {
	log.Info("Creating SOCKS proxy...")
	path := "./infrastructure/dev/k8s/make-tunnel.sh"

	runBashScript(path, "socks")
}

func portForwardVault() {
	log.Info("Port Forwarding to Vault...")
	path := "./infrastructure/dev/k8s/port-forward.sh"

	runBashScript(path, "vault")
}

func portForwardCelerium() {
	log.Info("Port Forwarding to Celerium...")
	path := "./infrastructure/dev/k8s/port-forward.sh"

	runBashScript(path, "celerium")
}

func runBashScript(path string, arg string) {
	makeTunnelCmd := exec.Command(path, arg)
	makeTunnelCmd.Stdout = os.Stdout
	makeTunnelCmd.Stderr = os.Stderr

	err := makeTunnelCmd.Start()
	if err != nil {
		log.WithError(err).Fatalf("Failed to run %s.", path)
	}

	err = makeTunnelCmd.Wait()
	if err != nil {
		log.WithError(err).Fatalf("An error occurred while running %s.", path)
	}
}
