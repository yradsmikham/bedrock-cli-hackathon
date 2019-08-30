package cmd

import (
	"os/exec"

	"github.com/kyokomi/emoji"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var servicePrincipal string
var secret string

// Demo is a function that will automate all the steps to creating an Azure Simple Cluster
func Demo(servicePrincipal string, secret string) (err error) {

	// Check for prerequisites
	requiredSystemTools := []string{"git", "helm", "sh", "curl", "terraform", "az"}
	for _, tool := range requiredSystemTools {
		path, err := exec.LookPath(tool)
		if err != nil {
			return err
		}
		log.Info(emoji.Sprintf(":mag: Using %s: %s", tool, path))
	}

	// Generate .tfvars file
	log.Info(emoji.Sprintf(":checkered_flag: Initializing Azure Simple Environment"))
	if _, _, error := Init(SIMPLE, "bedrock-demo-cluster"); error != nil {
		log.Error(emoji.Sprintf(":no_entry_sign: %s", error))
		return error
	}

	// Run terraform init and terraform plan
	if error := Simulate("bedrock/cluster/environments/bedrock-demo-cluster"); error != nil {
		log.Error(emoji.Sprintf(":no_entry_sign: %s", error))
		return error
	}

	// Run terraform apply
	if error := Deploy("bedrock/cluster/environments/bedrock-demo-cluster"); error != nil {
		log.Error(emoji.Sprintf(":no_entry_sign: %s", error))
		return error
	}

	log.Info(emoji.Sprintf(":raised_hands: Cluster has been successfully created!"))
	return err
}

var demoCmd = &cobra.Command{
	Use:   "demo --sp service-principal-app-id --secret service-principal-password --gitops-ssh-url manifest-repo-url-in-ssh-format",
	Short: "Demo an Azure Kubernetes Service (AKS) cluster using Terraform",
	Long:  `Demo an Azure Kubernetes Service (AKS) cluster using Terraform`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return Demo(servicePrincipal, secret)
	},
}

func init() {
	demoCmd.Flags().StringVar(&servicePrincipal, "sp", "", "Service Principal App Id")
	demoCmd.Flags().StringVar(&secret, "secret", "", "Password for the Service Principal")
	demoCmd.Flags().StringVar(&gitopsSSHUrl, "gitops-ssh-url", "git@github.com:timfpark/fabrikate-cloud-native-manifests.git", "The git repo that contains the resource manifests that should be deployed in the cluster in ssh format.")
	if error := demoCmd.MarkFlagRequired("sp"); error != nil {
		return
	}
	if error := demoCmd.MarkFlagRequired("secret"); error != nil {
		return
	}
	rootCmd.AddCommand(demoCmd)
}
