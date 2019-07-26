package cmd

import (
	"github.com/kyokomi/emoji"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Initializes the configuration for the given environment
func azureMultiCluster(servicePrincipal string, secret string) (err error) {
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
	Init(MULTIPLE)
	return err
}

var azureMultiClusterCmd = &cobra.Command{
	Use:   MULTIPLE + " --sp service-principal-app-id --secret service-principal-password --storage-account storage-account-name --access-key storage-account-access-key --container-name storage-container-name [--gitops-ssh-url manifest repo url in ssh format] [--tenant tenant-id | --common-infra-path common-infra-path]",
	Short: "Deploys Bedrock Multiple Azure Kubernetes Service (AKS) cluster configuration",
	Long:  `Deploys Bedrock Multiple Azure Kubernetes Service (AKS) cluster configuration`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		return azureMultiCluster(servicePrincipal, secret)
	},
}

func init() {
	azureMultiClusterCmd.Flags().StringVar(&servicePrincipal, "sp", "", "Service Principal App Id")
	azureMultiClusterCmd.Flags().StringVar(&secret, "secret", "", "Password for the Service Principal")
	azureMultiClusterCmd.Flags().StringVar(&gitopsSSHUrl, "gitops-ssh-url", "git@github.com:timfpark/fabrikate-cloud-native-manifests.git", "The git repo that contains the resource manifests that should be deployed in the cluster in ssh format.")
	azureMultiClusterCmd.Flags().StringVar(&commonInfraPath, "common-infra-path", "", "Common infra path for a successful deployment")
	azureMultiClusterCmd.Flags().StringVar(&tenant, "tenant", "", "Tenant ID for the Service Principal")
	azureMultiClusterCmd.Flags().StringVar(&storageAccount, "storage-account", "", "Storage Account Name")
	azureMultiClusterCmd.Flags().StringVar(&accessKey, "access-key", "", "Storage Account Access Key")
	azureMultiClusterCmd.Flags().StringVar(&containerName, "container-name", "", "Storage Container Name")
	azureMultiClusterCmd.MarkFlagRequired("sp")
	azureMultiClusterCmd.MarkFlagRequired("secret")
	rootCmd.AddCommand(azureMultiClusterCmd)
}
