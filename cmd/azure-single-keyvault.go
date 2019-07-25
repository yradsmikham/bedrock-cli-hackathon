package cmd

import (
	"github.com/kyokomi/emoji"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Initializes the configuration for the given environment
func azureSingleKeyvault(servicePrincipal string, secret string) (err error) {
	if tenant == "" && commonInfraPath == "" {
		log.Error(emoji.Sprintf(":confounded: One of common-infra-path and tenant need to be specified"))
		return err
	}
	Init(KEYVAULT)
	return err
}

var commonInfraPath string
var subscription string

var azureSingleKeyvaultCmd = &cobra.Command{
	Use:   KEYVAULT + " --subscription sub --sp service-principal-app-id --secret service-principal-password [--gitops-ssh-url manifest repo url in ssh format] [--tenant tenant-id | --common-infra-path common-infra-path]",
	Short: "Deploys a Bedrock Azure Kubernetes Service (AKS) cluster with an Azure Key Vault",
	Long:  `Deploys a Bedrock Azure Kubernetes Service (AKS) cluster with an Azure Key Vault. Make sure a successful deployment of ` + COMMON + ` is complete before attempting to deploy this one`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		return azureSingleKeyvault(servicePrincipal, secret)
	},
}

func init() {
	azureSingleKeyvaultCmd.Flags().StringVar(&servicePrincipal, "sp", "", "Service Principal App Id")
	azureSingleKeyvaultCmd.Flags().StringVar(&subscription, "subscription", "", "Subscription Id")
	azureSingleKeyvaultCmd.Flags().StringVar(&secret, "secret", "", "Password for the Service Principal")
	azureSingleKeyvaultCmd.Flags().StringVar(&gitopsSSHUrl, "gitops-ssh-url", "git@github.com:timfpark/fabrikate-cloud-native-manifests.git", "The git repo that contains the resource manifests that should be deployed in the cluster in ssh format")
	azureSingleKeyvaultCmd.Flags().StringVar(&commonInfraPath, "common-infra-path", "", "Common infra path for a successful deployment")
	azureSingleKeyvaultCmd.Flags().StringVar(&tenant, "tenant", "", "Tenant ID for the Service Principal")
	azureSingleKeyvaultCmd.MarkFlagRequired("sp")
	azureSingleKeyvaultCmd.MarkFlagRequired("secret")
	azureSingleKeyvaultCmd.MarkFlagRequired("subscription")

	rootCmd.AddCommand(azureSingleKeyvaultCmd)
}
