package cmd

import (
	"github.com/spf13/cobra"
)

// Initializes the configuration for the given environment
func azureSimple(servicePrincipal string, secret string) (err error) {
	Init("simple")
	return err
}

var servicePrincipal string
var secret string

var azureSimpleCmd = &cobra.Command{
	Use:   "azure-simple [--sp service-principal-app-id] [--secret service-principal-password]",
	Short: "Deploys a Bedrock Simple Azure Kubernetes Service (AKS) cluster configuration",
	Long:  `Deploys a Bedrock Simple Azure Kubernetes Service (AKS) cluster configuration`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		return azureSimple(servicePrincipal, secret)
	},
}

func init() {
	azureSimpleCmd.Flags().StringVar(&servicePrincipal, "sp", "", "Service Principal App Id")
	azureSimpleCmd.Flags().StringVar(&secret, "secret", "", "Password for the Service Principal")
	azureSimpleCmd.MarkFlagRequired("sp")
	azureSimpleCmd.MarkFlagRequired("secret")
	//azureSimpleCmd.MarkFlagRequired("gitops-ssh-url")

	rootCmd.AddCommand(azureSimpleCmd)
}
