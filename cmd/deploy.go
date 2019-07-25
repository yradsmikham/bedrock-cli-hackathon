package cmd

import (
	"os/exec"
	// "io/ioutil"

	"github.com/kyokomi/emoji"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/yradsmikham/bedrock-cli/utils"
)

// Apply and deploy a bedrock environment (azure simple, multi-cluster, keyvault, etc.)
func Deploy(name string) (err error) {
	log.Info(emoji.Sprintf(":eyes: Starting environment deployment!"))
	
	// TODO: For each subdirectory inside the named environment directory, run terraform init and plan?
	// Alternatively, just look for *common*, and run that directory first?

	// Terraform Init
	utils.TerraformInit(name)

	// Terraform Plan (terraform plan -var-file=./bedrock-terraform.tfvars)
	// TODO: Check that ./bedrock-terraform.tfvars exists. Throw error if it doesn't, (or default to terraform.tfvars?)
	utils.TerraformApply(name)

	// KUBECONFIG=./output/bedrock_kube_config:~/.kube/config kubectl config view --flatten > merged-config && mv merged-config ~/.kube/config
	// TODO: Check if a bedrock kubeconfig output exists, then add that file to the local kubeconfig.
	log.Info(emoji.Sprintf(":mailbox_with_mail: Found Kubeconfig output. Merging into local kubeconfig."))
	mergeConfigCmd := exec.Command("/bin/sh", "-c", "KUBECONFIG=./output/bedrock_kube_config:~/.kube/config kubectl config view --flatten > merged-config && mv merged-config ~/.kube/config")
	mergeConfigCmd.Dir = name
	if output, err := mergeConfigCmd.CombinedOutput(); err != nil {
		log.Error(emoji.Sprintf(":no_entry_sign: %s: %s", err, output))
		return err
	}

	if err == nil {
		log.Info(emoji.Sprintf(":raised_hands: Completed Terraform environment deployment!"))
	}
	
	return err
}

var deployCmd = &cobra.Command{
	Use:   "deploy <environment-name>",
	Short: "Deploy the bedrock environment using Terraform",
	Long:  `Deploy the bedrock environment deployment using terraform init and apply and adds the cluster credentials to the local kubeconfig.`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		var name = "unique-environment-name"

		if len(args) > 0 {
			name = args[0]
		} 
		return Deploy(name)
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
}
