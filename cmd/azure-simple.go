package cmd

import (
	"github.com/spf13/cobra"
)

var gitopsSSHUrl string

// Initializes the configuration for the given environment
func azureSimple(servicePrincipal string, secret string) (err error) {
	Init(SIMPLE)
	return err
}

var azureSimpleCmd = &cobra.Command{
	Use:   SIMPLE + " --sp service-principal-app-id --secret service-principal-password [--gitops-ssh-url manifest repo url in ssh format]",
	Short: "Deploys a Bedrock Simple Azure Kubernetes Service (AKS) cluster configuration",
	Long:  `Deploys a Bedrock Simple Azure Kubernetes Service (AKS) cluster configuration`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		return azureSimple(servicePrincipal, secret)
	},
}

func init() {
	azureSimpleCmd.Flags().StringVar(&servicePrincipal, "sp", "", "Service Principal App Id")
	azureSimpleCmd.Flags().StringVar(&secret, "secret", "", "Password for the Service Principal")
	azureSimpleCmd.Flags().StringVar(&gitopsSSHUrl, "gitops-ssh-url", "git@github.com:timfpark/fabrikate-cloud-native-manifests.git", "The git repo that contains the resource manifests that should be deployed in the cluster in ssh format.")
	azureSimpleCmd.MarkFlagRequired("sp")
	azureSimpleCmd.MarkFlagRequired("secret")

	rootCmd.AddCommand(azureSimpleCmd)
}
