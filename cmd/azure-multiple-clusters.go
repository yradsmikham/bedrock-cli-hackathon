package cmd

import (
	"github.com/kyokomi/emoji"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var gitopsPathWest string
var gitopsPathEast string
var gitopsPathCentral string
var gitopsURLBranchWest string
var gitopsURLBranchEast string
var gitopsURLBranchCentral string

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

	if _, error := Init(MULTIPLE, clusterName); error != nil {
		return error
	}
	return err
}

var azureMultiClusterCmd = &cobra.Command{
	Use:   MULTIPLE + " --subscription subscription-id --sp service-principal-app-id --secret service-principal-password --tenant tenant-id --gitops-ssh-url manifest-repo-url-in-ssh-format [--cluster-name name-of-AKS-cluster] [--vm-count number-of-nodes-to-deploy-in-cluster] [--dns-prefix DNS-prefix] [--poll-interval flux-sync-poll-interval] [--west-repo-path path-in-repo-to-sync-for-west-cluster] [--central-repo-path path-in-repo-to-sync-for-central-cluster] [--east-repo-path path-in-repo-to-sync-for-east-cluster] [--west-branch repo-branch-to-sync-with-for-west-cluster] [--central-branch repo-branch-to-sync-with-for-central-cluster] [--east-branch repo-branch-to-sync-with-for-east-cluster] [--keyvault name-of-keyvault] [--keyvault-rg name-of-resource-group-for-keyvault]",
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
	azureMultiClusterCmd.Flags().StringVar(&region, "region", "westus2", "Region of deployment")
	azureMultiClusterCmd.Flags().StringVar(&vmCount, "vm-count", "3", "Number of nodes to deploy per cluster")
	azureMultiClusterCmd.Flags().StringVar(&vmSize, "vm-size", "Standard_D4s_v3", "Azure VM size")
	azureMultiClusterCmd.Flags().StringVar(&dnsPrefix, "dns-prefix", "", "DNS Prefix")
	azureMultiClusterCmd.Flags().StringVar(&gitopsPollInterval, "poll-interval", "5m", "Period at which to poll git repo for new commits")
	azureMultiClusterCmd.Flags().StringVar(&keyvaultName, "keyvault", "", "Name of Key Vault")
	azureMultiClusterCmd.Flags().StringVar(&keyvaultRG, "keyvault-rg", "", "Resource group of Key Vault")
	azureMultiClusterCmd.Flags().StringVar(&gitopsPathWest, "west-repo-path", "", "Path in repo to sync with")
	azureMultiClusterCmd.Flags().StringVar(&gitopsPathEast, "east-repo-path", "", "Path in repo to sync with")
	azureMultiClusterCmd.Flags().StringVar(&gitopsPathCentral, "central-repo-path", "", "Path in repo to sync with")
	azureMultiClusterCmd.Flags().StringVar(&gitopsURLBranchWest, "west-branch", "master", "Path in repo to sync with")
	azureMultiClusterCmd.Flags().StringVar(&gitopsURLBranchEast, "east-branch", "master", "Path in repo to sync with")
	azureMultiClusterCmd.Flags().StringVar(&gitopsURLBranchCentral, "central-branch", "master", "Path in repo to sync with")
	if error := azureMultiClusterCmd.MarkFlagRequired("sp"); error != nil {
		return
	}
	if error := azureMultiClusterCmd.MarkFlagRequired("secret"); error != nil {
		return
	}
	rootCmd.AddCommand(azureMultiClusterCmd)
}
