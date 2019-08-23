package cmd

import (
	"io/ioutil"
	"os/exec"

	"github.com/kyokomi/emoji"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	util "github.com/yradsmikham/bedrock-cli/util"
)

// Deploy a bedrock environment by executing `terraform apply`
func Deploy(name string) (err error) {
	log.Info(emoji.Sprintf(":rocket: Starting Environment Deployment!"))

	files, err := ioutil.ReadDir(name)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		log.Info(emoji.Sprintf(":eyes: Searching for azure-common-infra environment..."))
		if f.Name() == COMMON {
			log.Info(emoji.Sprintf(":round_pushpin: Azure-common-infra environment found!"))
			setEnv(name, COMMON)

			// Terraform Init
			if error := util.TerraformInitBackend(name + "/azure-common-infra"); error != nil {
				return error
			}

			// Terraform Plan
			if error := util.TerraformApply(name + "/azure-common-infra"); error != nil {
				return error
			}

			break
		}
	}

	// Run Terraform Init on everything else (e.g. azure-single-keyvault, azure-multi-cluster)
	for _, f := range files {
		if f.Name() == SIMPLE {
			log.Info(emoji.Sprintf(":dancers: Deploying Azure-Simple Environment"))
			setEnv(name, SIMPLE)
			// Terraform Init
			if error := util.TerraformInit(name + "/azure-simple"); error != nil {
				return error
			}

			// Terraform Plan
			if error := util.TerraformApply(name + "/azure-simple"); error != nil {
				return error
			}

			// Add to local kubeconfig
			log.Info(emoji.Sprintf(":mailbox_with_mail: Found Kubeconfig output. Merging into local kubeconfig."))
			mergeConfigCmd := exec.Command("/bin/sh", "-c", "KUBECONFIG=./output/bedrock_kube_config:~/.kube/config kubectl config view --flatten > merged-config && mv merged-config ~/.kube/config")
			mergeConfigCmd.Dir = name + "/azure-simple"
			if output, err := mergeConfigCmd.CombinedOutput(); err != nil {
				log.Error(emoji.Sprintf(":no_entry_sign: %s: %s", err, output))
				return err
			}

			break
		}
		if f.Name() == KEYVAULT {
			log.Info(emoji.Sprintf(":dancers: Deploying Azure-Single-Keyvault Environment"))
			setEnv(name, KEYVAULT)

			// Terraform Init
			if error := util.TerraformInitBackend(name + "/azure-single-keyvault"); error != nil {
				return error
			}

			// Terraform Plan
			if error := util.TerraformApply(name + "/azure-single-keyvault"); error != nil {
				return error
			}

			// Add to local kubeconfig
			log.Info(emoji.Sprintf(":mailbox_with_mail: Found Kubeconfig output. Merging into local kubeconfig."))
			mergeConfigCmd := exec.Command("/bin/sh", "-c", "KUBECONFIG=./output/bedrock_kube_config:~/.kube/config kubectl config view --flatten > merged-config && mv merged-config ~/.kube/config")
			mergeConfigCmd.Dir = name + "/azure-single-keyvault"
			if output, err := mergeConfigCmd.CombinedOutput(); err != nil {
				log.Error(emoji.Sprintf(":no_entry_sign: %s: %s", err, output))
				return err
			}

			break
		}
		if f.Name() == MULTIPLE {
			log.Info(emoji.Sprintf(":dancers: Deploying Azure-Multiple-Clusters Environment"))
			setEnv(name, MULTIPLE)

			// Terraform Init
			if error := util.TerraformInitBackend(name + "/azure-multiple-clusters"); error != nil {
				return error
			}

			// Terraform Plan
			if error := util.TerraformApply(name + "/azure-multiple-clusters"); error != nil {
				return error
			}

			// TO-DO: For multiple cluster, must add each cluster individually
			log.Info(emoji.Sprintf(":mailbox_with_mail: Found Kubeconfig output. Merging into local kubeconfig."))
			mergeConfigCmd := exec.Command("/bin/sh", "-c", "for f in output/*kube_config; do KUBECONFIG=$f:~/.kube/config kubectl config view --flatten > merged-config && mv merged-config ~/.kube/config; done")
			mergeConfigCmd.Dir = name + "/azure-multiple-clusters"
			if output, err := mergeConfigCmd.CombinedOutput(); err != nil {
				log.Error(emoji.Sprintf(":no_entry_sign: %s: %s", err, output))
				return err
			}

			break
		}
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
