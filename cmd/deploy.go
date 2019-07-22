package cmd

import (
	"errors"
	"os/exec"

	"github.com/kyokomi/emoji"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Create a cluster environment (azure simple, multi-cluster, keyvault, etc.)
func Deploy() (err error) {
	// Make sure host system contains all utils needed by Fabrikate
	requiredSystemTools := []string{"git", "helm", "sh", "curl", "terraform", "az"}
	for _, tool := range requiredSystemTools {
		path, err := exec.LookPath(tool)
		if err != nil {
			return err
		}
		log.Info(emoji.Sprintf(":mag: Using %s: %s", tool, path))
	}

	// Terraform Apply
	log.Info(emoji.Sprintf(":rocket: Terraform Apply"))
	cmd2 := exec.Command("terraform", "apply", "-auto-approve", "--var", "resource_group_name=bedrock-cli-demo-simple-cluster", "--var", "cluster_name=bedrock-cli-demo-simple-cluster", "--var", "dns_prefix=bedrock-cli-demo-simple-cluster", "--var", "service_principal_id=<app-Id>", "--var", "service_principal_secret=<password>", "--var", "ssh_public_key=<ssh public key>", "--var", "vnet_name=bedrock-cli-demo-simple-cluster", "--var", "gitops_ssh_key=/path/to/private/key")
	cmd2.Dir = "path/to/terraform/config"
	if output, err := cmd2.CombinedOutput(); err != nil {
		log.Error(emoji.Sprintf(":no_entry_sign: %s: %s", err, output))
		return err
	}

	// Copy to KUBECONFIG
	log.Info(emoji.Sprintf(":heavy_plus_sign: Download Credentials for Kubernetes Cluster"))
	if output, err := exec.Command("az", "aks", "get-credentials", "--resource-group", "bedrock-cli-demo-simple-cluster", "--name", "bedrock-cli-demo-simple-cluster").CombinedOutput(); err != nil {
		log.Error(emoji.Sprintf(":no_entry_sign: %s: %s", err, output))
		return err
	}

	if err == nil {
		log.Info(emoji.Sprintf(":raised_hands: Cluster has been successfully created!"))
	}

	return err
}

var deployCmd = &cobra.Command{
	Use:   "create <config>",
	Short: "Create an Azure Kubernetes Service (AKS) cluster using Terraform",
	Long:  `Create an Azure Kubernetes Service (AKS) cluster using Terraform`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		if !((args[0] == "simple") || (args[0] == "multi")) {
			return errors.New("the environment you specified is not of the following: simple, multi, keyvault")
		}
		return Deploy()
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
}
