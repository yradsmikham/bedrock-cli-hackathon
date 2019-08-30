package cmd

import (
	"github.com/kyokomi/emoji"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var commonInfraPath string
var subscription string
var vmSize string
var addressSpace string
var subnetPrefix string
var keyvaultName string
var keyvaultRG string
var commonInfraCreation bool

// Initializes the configuration for the given environment
func azureSingleKeyvault(servicePrincipal string, secret string, commonInfraPath string) (err error) {
	if commonInfraPath != "" {
		log.Info(emoji.Sprintf(":star2: An Azure Common Infra enviroment has been provided"))
		commonInfraCreation = false
	} else {
		commonInfraCreation = true
	}

	if storageAccount != "" && (accessKey == "" || containerName == "") {
		log.Error(emoji.Sprintf(":confounded: Please specify the Storage Account Access Key using --access-key and Storage Container using --container-name"))
		return err
	}

	if containerName != "" && (accessKey == "" || storageAccount == "") {
		log.Error(emoji.Sprintf(":confounded: Please specify a Storage Account Name using '--storage-account' and the Storage Account Access Key --access-key"))
		return err
	}

	if _, _, error := Init(KEYVAULT, clusterName); error != nil {
		return error
	}
	return err
}

var azureSingleKeyvaultCmd = &cobra.Command{
	Use:   KEYVAULT + " --gitops-ssh-url manifest-repo-url-in-ssh-format [--subscription subscription-id] [--sp service-principal-app-id] [--secret service-principal-password] [--tenant serice-principal-tenant-id] [--common-infra-path path-to-azure-common-infra-environment] [--storage-account storage-account-name] [--access-key storage-account-access-key] [--container-name storage-container-name] [--cluster-name name-of-AKS-cluster] [--region region-of-deployment] [--vm-count number-of-nodes-to-deploy-in-cluster] [--vm-size azure-vm-size] [--dns-prefix DNS-prefix] [--poll-interval flux-sync-poll-interval] [--repo-path path-in-repo-to-sync] [--branch repo-branch-to-sync-with] [--keyvault name-of-keyvault] [--keyvault-rg name-of-resource-group-for-keyvault] [--address-space address-space] [--subnet-prefix subnet-prefixes]",
	Short: "Deploys a Bedrock Azure Kubernetes Service (AKS) cluster with an Azure Key Vault",
	Long:  `Deploys a Bedrock Azure Kubernetes Service (AKS) cluster with an Azure Key Vault. Make sure a successful deployment of ` + COMMON + ` is complete before attempting to deploy this one`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		return azureSingleKeyvault(servicePrincipal, secret, commonInfraPath)
	},
}

func init() {
	azureSingleKeyvaultCmd.Flags().StringVar(&resourceGroup, "resource-group", "", "An existing Azure Resource Group")
	azureSingleKeyvaultCmd.Flags().StringVar(&servicePrincipal, "sp", "", "Service Principal App ID")
	azureSingleKeyvaultCmd.Flags().StringVar(&subscription, "subscription", "", "Subscription ID")
	azureSingleKeyvaultCmd.Flags().StringVar(&secret, "secret", "", "Password for the Service Principal")
	azureSingleKeyvaultCmd.Flags().StringVar(&tenant, "tenant", "", "Tenant ID for the Service Principal")
	azureSingleKeyvaultCmd.Flags().StringVar(&gitopsSSHUrl, "gitops-ssh-url", "git@github.com:timfpark/fabrikate-cloud-native-manifests.git", "The git repo that contains the resource manifests that should be deployed in the cluster in ssh format")
	azureSingleKeyvaultCmd.Flags().StringVar(&commonInfraPath, "common-infra-path", "", "Successful deployment of an Azure Common Infra environment")
	azureSingleKeyvaultCmd.Flags().StringVar(&storageAccount, "storage-account", "", "Storage Account Name")
	azureSingleKeyvaultCmd.Flags().StringVar(&accessKey, "access-key", "", "Storage Account Access Key")
	azureSingleKeyvaultCmd.Flags().StringVar(&containerName, "container-name", "", "Storage Container Name")
	azureSingleKeyvaultCmd.Flags().StringVar(&clusterName, "cluster-name", "", "Name of AKS Cluster")
	azureSingleKeyvaultCmd.Flags().StringVar(&region, "region", "westus2", "Region of deployment")
	azureSingleKeyvaultCmd.Flags().StringVar(&vmCount, "vm-count", "3", "Number of nodes to deploy per cluster")
	azureSingleKeyvaultCmd.Flags().StringVar(&vmSize, "vm-size", "Standard_D4s_v3", "Azure VM size")
	azureSingleKeyvaultCmd.Flags().StringVar(&dnsPrefix, "dns-prefix", "", "DNS Prefix")
	azureSingleKeyvaultCmd.Flags().StringVar(&gitopsPollInterval, "poll-interval", "5m", "Period at which to poll git repo for new commits")
	azureSingleKeyvaultCmd.Flags().StringVar(&gitopsPath, "repo-path", "", "Path in repo to sync with")
	azureSingleKeyvaultCmd.Flags().StringVar(&gitopsURLBranch, "branch", "master", "Path in repo to sync with")
	azureSingleKeyvaultCmd.Flags().StringVar(&addressSpace, "address-space", "10.39.0.0/24", "CIDR for cluster address space")
	azureSingleKeyvaultCmd.Flags().StringVar(&subnetPrefix, "subnet-prefix", "10.39.0.0/24", "Subnet prefixes")
	azureSingleKeyvaultCmd.Flags().StringVar(&keyvaultName, "keyvault", "", "Name of Key Vault")
	azureSingleKeyvaultCmd.Flags().StringVar(&keyvaultRG, "keyvault-rg", "", "Resource group of Key Vault")
	if error := azureSingleKeyvaultCmd.MarkFlagRequired("gitops-ssh-url"); error != nil {
		return
	}
	rootCmd.AddCommand(azureSingleKeyvaultCmd)
}
