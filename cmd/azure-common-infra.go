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
var commonInfraName string

// Initializes the configuration for the given environment
func commonInfra(servicePrincipal string, tenant string, storageAccount string, accessKey string, containerName string) (err error) {
	commonInfraName, error := Init(COMMON, clusterName)
	if error != nil {
		return error
	}

	log.Info(emoji.Sprintf(":white_check_mark: To proceed, run 'bedrock simulate bedrock/cluster/environments/" + commonInfraName + "'"))

	return err
}

var commonInfraCmd = &cobra.Command{
	Use:   COMMON + " --sp service-principal-app-id --tenant tenant-id --storage-account storage-account-name --access-key access-key --container-name storage-container-name [--cluster-name name-of-AKS-cluster] [--region region-of-resource] [--keyvault name-of-keyvault] [--keyvault-rg name-of-resource-group-for-keyvault] [--address-space address-space] [--subnet-prefix subnet-prefixes]",
	Short: "Deploys the Bedrock Common Infra Environment",
	Long:  `Deploys the Bedrock Common Infra Environment`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		return commonInfra(servicePrincipal, tenant, storageAccount, accessKey, containerName)
	},
}

func init() {
	commonInfraCmd.Flags().StringVar(&servicePrincipal, "sp", "", "Service Principal App Id")
	commonInfraCmd.Flags().StringVar(&tenant, "tenant", "", "Tenant ID for the Service Principal")
	commonInfraCmd.Flags().StringVar(&storageAccount, "storage-account", "", "Storage Account Name")
	commonInfraCmd.Flags().StringVar(&accessKey, "access-key", "", "Acces Key for the Storage Account")
	commonInfraCmd.Flags().StringVar(&containerName, "container-name", "", "Storage Container Name")
	commonInfraCmd.Flags().StringVar(&clusterName, "cluster-name", "", "Name of AKS Cluster")
	commonInfraCmd.Flags().StringVar(&region, "region", "westus2", "Region of deployment")
	commonInfraCmd.Flags().StringVar(&addressSpace, "address-space", "10.39.0.0/24", "CIDR for cluster address space")
	commonInfraCmd.Flags().StringVar(&subnetPrefix, "subnet-prefix", "10.39.0.0/24", "Subnet prefixes")
	commonInfraCmd.Flags().StringVar(&keyvaultRG, "global-rg", "", "Resource group of Key Vault")
	commonInfraCmd.Flags().StringVar(&keyvaultName, "keyvault", "", "Name of Key Vault")
	if error := commonInfraCmd.MarkFlagRequired("sp"); error != nil {
		return
	}
	if error := commonInfraCmd.MarkFlagRequired("tenant"); error != nil {
		return
	}
	rootCmd.AddCommand(commonInfraCmd)
}
