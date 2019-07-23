package cmd

import (
	"errors"
	"os/exec"

	"github.com/kyokomi/emoji"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Create a cluster environment (azure simple, multi-cluster, keyvault, etc.)
func Simulate() (err error) {
	// Make sure host system contains all utils needed by Fabrikate
	requiredSystemTools := []string{"git", "helm", "sh", "curl", "terraform", "az"}
	for _, tool := range requiredSystemTools {
		path, err := exec.LookPath(tool)
		if err != nil {
			return err
		}
		log.Info(emoji.Sprintf(":mag: Using %s: %s", tool, path))
	}

	// Terraform Initialization
	log.Info(emoji.Sprintf(":checkered_flag: Terraform Init"))
	cmd := exec.Command("terraform", "init")
	cmd.Dir = "path/to/terraform/config"
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Error(emoji.Sprintf(":no_entry_sign: %s: %s", err, output))
		return err
	}

	// Terraform Plan

	if err == nil {
		log.Info(emoji.Sprintf(":raised_hands: Cluster has been successfully created!"))
	}

	return err
}

var simulateCmd = &cobra.Command{
	Use:   "simulate <config>",
	Short: "Simulate an Azure Kubernetes Service (AKS) cluster using Terraform",
	Long:  `Simulate an Azure Kubernetes Service (AKS) cluster using Terraform`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		if !((args[0] == "simple") || (args[0] == "multi")) {
			return errors.New("the environment you specified is not of the following: simple, multi, keyvault")
		}
		return Simulate()
	},
}

func init() {
	rootCmd.AddCommand(simulateCmd)
}
