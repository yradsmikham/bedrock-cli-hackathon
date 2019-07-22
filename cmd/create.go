package cmd

import (
	"errors"
	"os/exec"

	"github.com/kyokomi/emoji"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Create a cluster environment (azure simple, multi-cluster, keyvault, etc.)
func Create(environment string) (err error) {
	// Make sure host system contains all utils needed by Fabrikate
	requiredSystemTools := []string{"git", "helm", "sh", "curl", "terraform", "az"}
	for _, tool := range requiredSystemTools {
		path, err := exec.LookPath(tool)
		if err != nil {
			return err
		}
		log.Info(emoji.Sprintf(":mag: Using %s: %s", tool, path))
	}

	// Clone Microsoft Bedrock Repo
	log.Info(emoji.Sprintf(":crystal_ball: Cloning Bedrock Repo"))
	if output, err := exec.Command("git", "clone", "https://github.com/microsoft/bedrock").CombinedOutput(); err != nil {
		log.Error(emoji.Sprintf(":no_entry_sign: %s: %s", err, output))
		return err
	}

	// Create an Azure SP
	/* 	log.Info(emoji.Sprintf(":cop: Creating New Service Principal"))
	   	if output, err := exec.Command("az", "ad", "sp", "create-for-rbac", "--subscription", "<subscription-id>").CombinedOutput(); err != nil {
	   		log.Error(emoji.Sprintf(":no_entry_sign: %s: %s", err, output))
	   		return err
	   	} */

	// Copy Terraform Template
	log.Info(emoji.Sprintf(":flashlight: Creating New Environment"))
	if output, err := exec.Command("cp", "-r", "bedrock/cluster/environments/azure-simple", "bedrock/cluster/environments/bedrock-cli-demo").CombinedOutput(); err != nil {
		log.Error(emoji.Sprintf(":no_entry_sign: %s: %s", err, output))
		return err
	}

	// Create SSH Keys
	log.Info(emoji.Sprintf(":closed_lock_with_key: Creating New SSH Keys"))
	if output, err := exec.Command("ssh-keygen", "-f", "bedrock/cluster/environments/bedrock-cli-demo/deploy_key", "-P", "''").CombinedOutput(); err != nil {
		log.Error(emoji.Sprintf(":no_entry_sign: %s: %s", err, output))
		return err
	}

	// Terraform Initialization
	log.Info(emoji.Sprintf(":checkered_flag: Terraform Init"))
	cmd := exec.Command("terraform", "init")
	cmd.Dir = "path/to/terraform/config"
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Error(emoji.Sprintf(":no_entry_sign: %s: %s", err, output))
		return err
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

var environment string

var createCmd = &cobra.Command{
	Use:   "create <config>",
	Short: "Create an Azure Kubernetes Service (AKS) cluster using Terraform",
	Long:  `Create an Azure Kubernetes Service (AKS) cluster using Terraform`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		if !((args[0] == "simple") || (args[0] == "multi")) {
			return errors.New("the environment you specified is not of the following: simple, multi, keyvault")
		}
		return Create(environment)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
