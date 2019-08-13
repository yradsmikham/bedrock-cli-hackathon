package cmd

import (
	"github.com/spf13/cobra"
)

var tenant string

// Initializes the configuration for the given environment
func commonInfra(servicePrincipal string, tenant string, storageAccount string, accessKey string, containerName string) (err error) {
	if error := Init(COMMON, clusterName); error != nil {
		return error
	}
	return err
}

var storageAccount string
var accessKey string
var containerName string

var commonInfraCmd = &cobra.Command{
	Use:   COMMON + " --sp service-principal-app-id --tenant tenant-id --storage-account storage-account-name --access-key access-key --container-name storage-container-name [--cluster-name name-of-AKS-cluster] ",
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
	if error := commonInfraCmd.MarkFlagRequired("sp"); error != nil {
		return
	}
	if error := commonInfraCmd.MarkFlagRequired("tenant"); error != nil {
		return
	}
	rootCmd.AddCommand(commonInfraCmd)
}
