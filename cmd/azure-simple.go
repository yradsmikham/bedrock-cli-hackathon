package cmd

import (
	"github.com/spf13/cobra"
)

var gitopsSSHUrl string
var region string
var vmCount string
var vnet string
var gitopsPollInterval string
var gitopsPath string
var gitopsURLBranch string
var dnsPrefix string
var resourceGroup string

// Initializes the configuration for the given environment
func azureSimple(servicePrincipal string, secret string) (err error) {
	if _, error := Init(SIMPLE, clusterName); error != nil {
		return error
	}
	return err
}

var azureSimpleCmd = &cobra.Command{
	Use:   SIMPLE + " --gitops-ssh-url manifest-repo-url-in-ssh-format [--subscription subscription-id] [--sp service-principal-app-id] [--secret service-principal-password] [--tenant serice-principal-tenant-id] [--cluster-name name-of-AKS-cluster] [--region region-of-deployment] [--vm-count number-of-nodes-to-deploy-in-cluster] [--vnet name-of-vnet] [--dns-prefix DNS-prefix] [--poll-interval flux-sync-poll-interval] [--repo-path path-in-repo-to-sync] [--branch repo-branch-to-sync-with]",
	Short: "Deploys a Bedrock Simple Azure Kubernetes Service (AKS) cluster configuration",
	Long:  `Deploys a Bedrock Simple Azure Kubernetes Service (AKS) cluster configuration`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		return azureSimple(servicePrincipal, secret)
	},
}

func init() {
	azureSimpleCmd.Flags().StringVar(&servicePrincipal, "sp", "", "Service Principal App Id")
	azureSimpleCmd.Flags().StringVar(&secret, "secret", "", "Password for the Service Principal")
	azureSimpleCmd.Flags().StringVar(&subscription, "subscription", "", "Azure Subscription ID")
	azureSimpleCmd.Flags().StringVar(&tenant, "tenant", "", "Tenant ID for Service Principal")
	azureSimpleCmd.Flags().StringVar(&resourceGroup, "resource-group", "", "An existing Azure Resource Group")
	azureSimpleCmd.Flags().StringVar(&gitopsSSHUrl, "gitops-ssh-url", "git@github.com:timfpark/fabrikate-cloud-native-manifests.git", "The git repo that contains the resource manifests that should be deployed in the cluster in ssh format")
	azureSimpleCmd.Flags().StringVar(&clusterName, "cluster-name", "", "Name of AKS Cluster")
	azureSimpleCmd.Flags().StringVar(&region, "region", "westus2", "Region of deployment")
	azureSimpleCmd.Flags().StringVar(&vmCount, "vm-count", "3", "Number of nodes to deploy per cluster")
	azureSimpleCmd.Flags().StringVar(&vnet, "vnet", "", "Name of vnet resource")
	azureSimpleCmd.Flags().StringVar(&dnsPrefix, "dns-prefix", "", "DNS Prefix")
	azureSimpleCmd.Flags().StringVar(&gitopsPollInterval, "poll-interval", "5m", "Period at which to poll git repo for new commits")
	azureSimpleCmd.Flags().StringVar(&gitopsPath, "repo-path", "", "Path in repo to sync with")
	azureSimpleCmd.Flags().StringVar(&gitopsURLBranch, "branch", "master", "Path in repo to sync with")
	if error := azureSimpleCmd.MarkFlagRequired("gitops-ssh-url"); error != nil {
		return
	}
	rootCmd.AddCommand(azureSimpleCmd)
}
