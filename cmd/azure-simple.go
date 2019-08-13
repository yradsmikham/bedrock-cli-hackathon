package cmd

import (
	"github.com/spf13/cobra"
)

var gitopsSSHUrl string

// Initializes the configuration for the given environment
func azureSimple(servicePrincipal string, secret string) (err error) {
	if error := Init(SIMPLE, clusterName); error != nil {
		return error
	}
	return err
}

var azureSimpleCmd = &cobra.Command{
	Use:   SIMPLE + " --sp service-principal-app-id --secret service-principal-password --gitops-ssh-url manifestr-repo-url-in-ssh-format [--cluster-name name-of-AKS-cluster]",
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
	azureSimpleCmd.Flags().StringVar(&clusterName, "cluster-name", "", "Name of AKS Cluster")
	if error := azureSimpleCmd.MarkFlagRequired("sp"); error != nil {
		return
	}
	if error := azureSimpleCmd.MarkFlagRequired("secret"); error != nil {
		return
	}
	rootCmd.AddCommand(azureSimpleCmd)
}
