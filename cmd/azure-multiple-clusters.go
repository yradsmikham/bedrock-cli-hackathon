package cmd

import (
	"github.com/kyokomi/emoji"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Initializes the configuration for the given environment
func azureMultiCluster(servicePrincipal string, secret string) (err error) {
	if tenant == "" {
		log.Error(emoji.Sprintf(":confounded: Please specify the Tenant ID for your Service Principal using '--tenant' argument"))
		return err
	}
	if subscription == "" {
		log.Error(emoji.Sprintf(":confounded: Please specify the Subcription ID using '--subscription' argument"))
		return err
	}

	Init(MULTIPLE, clusterName)
	return err
}

var azureMultiClusterCmd = &cobra.Command{
	Use:   MULTIPLE + " --subscription subscription-id --sp service-principal-app-id --secret service-principal-password --tenant tenant-id [--gitops-ssh-url manifest-repo-url-in-ssh-format | --cluster-name cluster-name ]",
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
	azureMultiClusterCmd.Flags().StringVar(&tenant, "tenant", "", "Tenant ID for the Service Principal")
	azureMultiClusterCmd.Flags().StringVar(&subscription, "subscription", "", "Subscription ID")
	azureMultiClusterCmd.Flags().StringVar(&clusterName, "cluster-name", "", "Name of AKS Cluster")
	if error := azureMultiClusterCmd.MarkFlagRequired("sp"); error != nil {
		return
	}
	if error := azureMultiClusterCmd.MarkFlagRequired("secret"); error != nil {
		return
	}
	rootCmd.AddCommand(azureMultiClusterCmd)
}
