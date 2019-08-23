package cmd

import (
	"github.com/kyokomi/emoji"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var tenant string
var storageAccount string
var accessKey string
var containerName string
var clusterName string

// Initializes the configuration for the given environment
func commonInfra(servicePrincipal string, storageAccount string, accessKey string, containerName string) (err error) {
	commonInfraName, error := Init(COMMON, clusterName)
	if error != nil {
		return error
	}

	log.Info(emoji.Sprintf(":white_check_mark: To proceed, run 'bedrock simulate bedrock/cluster/environments/" + commonInfraName + "'"))

	return err
}

var commonInfraCmd = &cobra.Command{
	Use:   COMMON + " [--subscription subscription-id] [--sp service-principal-app-id] [--secret service-principal-password] [--tenant serice-principal-tenant-id] [--storage-account storage-account-name] [--access-key access-key] [--container-name storage-container-name] [--cluster-name name-of-AKS-cluster] [--region region-of-resource] [--keyvault name-of-keyvault] [--keyvault-rg name-of-resource-group-for-keyvault] [--address-space address-space] [--subnet-prefix subnet-prefixes]",
	Short: "Deploys the Bedrock Common Infra Environment",
	Long:  `Deploys the Bedrock Common Infra Environment`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		return commonInfra(servicePrincipal, storageAccount, accessKey, containerName)
	},
}

func init() {
	commonInfraCmd.Flags().StringVar(&resourceGroup, "resource-group", "", "An existing Azure Resource Group")
	commonInfraCmd.Flags().StringVar(&servicePrincipal, "sp", "", "Service Principal App ID")
	commonInfraCmd.Flags().StringVar(&secret, "secret", "", "Password for  Service Principal")
	commonInfraCmd.Flags().StringVar(&subscription, "subscription", "", "Azure Subscription ID")
	commonInfraCmd.Flags().StringVar(&tenant, "tenant", "", "Tenant ID for the Service Principal")
	commonInfraCmd.Flags().StringVar(&storageAccount, "storage-account", "", "Storage Account Name")
	commonInfraCmd.Flags().StringVar(&accessKey, "access-key", "", "Acces Key for the Storage Account")
	commonInfraCmd.Flags().StringVar(&containerName, "container-name", "", "Storage Container Name")
	commonInfraCmd.Flags().StringVar(&clusterName, "cluster-name", "", "Name of AKS Cluster")
	commonInfraCmd.Flags().StringVar(&region, "region", "westus2", "Region of deployment")
	commonInfraCmd.Flags().StringVar(&addressSpace, "address-space", "10.39.0.0/24", "CIDR for cluster address space")
	commonInfraCmd.Flags().StringVar(&subnetPrefix, "subnet-prefix", "10.39.0.0/24", "Subnet prefixes")
	commonInfraCmd.Flags().StringVar(&keyvaultName, "keyvault", "", "Name of Key Vault")
	rootCmd.AddCommand(commonInfraCmd)
}
