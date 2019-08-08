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
	if storageAccount == "" {
		log.Error(emoji.Sprintf(":confounded: Please specify a Storage Account Name using '--storage-account' argument"))
		return err
	}

	if accessKey == "" {
		log.Error(emoji.Sprintf(":confounded: Please specify the Storage Access Key using '--access-key' argument"))
		return err
	}

	if containerName == "" {
		log.Error(emoji.Sprintf(":confounded: Please specify the Storage Container Name using '--container-name' argument"))
		return err
	}

	if _, error := Init(KEYVAULT, clusterName); error != nil {
		return error
	}
	return err
}

var commonInfraPath string
var subscription string

var azureSingleKeyvaultCmd = &cobra.Command{
	Use:   KEYVAULT + " --subscription subscription-id --sp service-principal-app-id --secret service-principal-password --storage-account storage-account-name --access-key storage-account-access-key --container-name storage-container-name [--gitops-ssh-url manifest-repo-url-in-ssh-format | --cluster-name cluster-name | --tenant tenant-id | --common-infra-path common-infra-path]",
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
	azureSingleKeyvaultCmd.Flags().StringVar(&storageAccount, "storage-account", "", "Storage Account Name")
	azureSingleKeyvaultCmd.Flags().StringVar(&accessKey, "access-key", "", "Storage Account Access Key")
	azureSingleKeyvaultCmd.Flags().StringVar(&containerName, "container-name", "", "Storage Container Name")
	azureSingleKeyvaultCmd.Flags().StringVar(&clusterName, "cluster-name", "", "Name of AKS Cluster")
	if error := azureSingleKeyvaultCmd.MarkFlagRequired("sp"); error != nil {
		return
	}
	if error := azureSingleKeyvaultCmd.MarkFlagRequired("secret"); error != nil {
		return
	}
	if error := azureSingleKeyvaultCmd.MarkFlagRequired("subscription"); error != nil {
		return
	}

	rootCmd.AddCommand(azureSingleKeyvaultCmd)
}
