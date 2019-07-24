package cmd

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/docker/docker/pkg/namesgenerator"
	"github.com/kyokomi/emoji"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Init functions initializes the configuration for the given environment
func Init(environment string) (err error) {
	requiredSystemTools := []string{"git", "helm", "sh", "curl", "terraform", "az"}
	for _, tool := range requiredSystemTools {
		path, err := exec.LookPath(tool)
		if err != nil {
			return err
		}
		log.Info(emoji.Sprintf(":mag: Using %s: %s", tool, path))
	}

	rand.Seed(time.Now().UnixNano())
	randomName := strings.Replace(namesgenerator.GetRandomName(0), "_", "-", -1)

	// Check if Bedrock Repo is already cloned
	log.Info(emoji.Sprintf(":open_file_folder: Checking for Bedrock"))
	if output, err := exec.Command("git", "clone", "https://github.com/microsoft/bedrock").CombinedOutput(); output != nil || err != nil {
		log.Info(emoji.Sprintf(":star: Bedrock Repo already cloned"))
	}

	// Copy Terraform Template
	environmentPath := "bedrock/cluster/environments/" + randomName
	os.MkdirAll(environmentPath, os.ModePerm)

	log.Info(emoji.Sprintf(":flashlight: Creating New Environment"))
	if output, err := exec.Command("cp", "-r", "bedrock/cluster/environments/azure-simple", environmentPath).CombinedOutput(); err != nil {
		log.Error(emoji.Sprintf(":no_entry_sign: %s: %s", err, output))
		return err
	}

	// To-do: Generate ssh keys
	fullEnvironmentPath := environmentPath + "/azure-simple"
	SSHKey, err := SSH(fullEnvironmentPath, "deploy-key")
	if err == nil {
		// Save bedrock-config.tfvars
		err = addConfigTemplate(environment, fullEnvironmentPath, randomName, SSHKey)

		if err == nil {
			log.Info(emoji.Sprintf(":raised_hands: Cluster " + fullEnvironmentPath + " has been successfully created!"))
			return nil
		}
	}

	return err
}

// Adds a blank bedrock config template
func addConfigTemplate(environment string, environmentPath string, clusterName string, SSHKey string) (err error) {
	SSHKey = strings.TrimSuffix(SSHKey, "\n")

	fmt.Println(environment)
	if environment == "simple" {
		azureSimpleConfig := make(map[string]string)

		azureSimpleConfig["resource_group_name"] = "\"" + clusterName + "-rg\""
		azureSimpleConfig["resource_group_location"] = "\"\""
		azureSimpleConfig["cluster_name"] = "\"" + clusterName + "\""
		azureSimpleConfig["agent_vm_count"] = "\"\""
		azureSimpleConfig["dns_prefix"] = "\"" + clusterName + "\""
		azureSimpleConfig["service_principal_id"] = servicePrincipal
		azureSimpleConfig["service_principal_secret"] = secret
		azureSimpleConfig["ssh_public_key"] = "\"" + SSHKey + "\""
		azureSimpleConfig["gitops_ssh_url"] = "\"\""
		azureSimpleConfig["gitops_ssh_key"] = "\"" + environmentPath + "\""
		azureSimpleConfig["vnet_name"] = "\"" + clusterName + "-vnet\""

		f, err := os.Create(environmentPath + "/bedrock-config.tfvars")
		log.Info(emoji.Sprintf(":raised_hands: Create Bedrock config file " + environmentPath + "/bedrock-config.tfvars"))
		if err != nil {
			return err
		}

		for setting, value := range azureSimpleConfig {
			f.WriteString(setting + " = " + value + "\n")
		}

		f.Close()

		return nil
	}

	if environment == "common-infra" {
		// TO-DO: Need to customize the config for common-infra
		commonInfraConfig := make(map[string]string)

		commonInfraConfig["resource_group_name"] = "\"" + clusterName + "-rg\""
		commonInfraConfig["resource_group_location"] = "\"\""
		commonInfraConfig["cluster_name"] = "\"" + clusterName + "\""
		commonInfraConfig["agent_vm_count"] = "\"\""
		commonInfraConfig["dns_prefix"] = "\"\""
		commonInfraConfig["service_principal_id"] = servicePrincipal
		commonInfraConfig["service_principal_secret"] = secret
		commonInfraConfig["ssh_public_key"] = "\"" + SSHKey + "\""
		commonInfraConfig["gitops_ssh_url"] = "\"\""
		commonInfraConfig["gitops_ssh_key"] = "\"" + environmentPath + "\""
		commonInfraConfig["vnet_name"] = "\"\""

		f, err := os.Create(environmentPath + "/bedrock-config.tfvars")
		log.Info(emoji.Sprintf(":raised_hands: Create Bedrock config file " + environmentPath + "/bedrock-config.tfvars"))
		if err != nil {
			return err
		}

		for setting, value := range commonInfraConfig {
			f.WriteString(setting + " = " + value + "\n")
		}

		f.Close()

		return nil
	}

	if environment == "single-keyvault" {
		// TO-DO: Need to customize the config for single keyvault

		singleKeyvaultConfig := make(map[string]string)

		singleKeyvaultConfig["resource_group_name"] = "\"" + clusterName + "-rg\""
		singleKeyvaultConfig["resource_group_location"] = "\"\""
		singleKeyvaultConfig["cluster_name"] = "\"" + clusterName + "\""
		singleKeyvaultConfig["agent_vm_count"] = "\"\""
		singleKeyvaultConfig["dns_prefix"] = "\"\""
		singleKeyvaultConfig["service_principal_id"] = servicePrincipal
		singleKeyvaultConfig["service_principal_secret"] = secret
		singleKeyvaultConfig["ssh_public_key"] = "\"" + SSHKey + "\""
		singleKeyvaultConfig["gitops_ssh_url"] = "\"\""
		singleKeyvaultConfig["gitops_ssh_key"] = "\"" + environmentPath + "\""
		singleKeyvaultConfig["vnet_name"] = "\"\""

		f, err := os.Create(environmentPath + "/bedrock-config.tfvars")
		log.Info(emoji.Sprintf(":raised_hands: Create Bedrock config file " + environmentPath + "/bedrock-config.tfvars"))
		if err != nil {
			return err
		}

		for setting, value := range singleKeyvaultConfig {
			f.WriteString(setting + " = " + value + "\n")
		}

		f.Close()

		return nil
	}
	return err
}

var initCmd = &cobra.Command{
	Use:   "init <config>",
	Short: "Init an Azure Kubernetes Service (AKS) cluster configuration",
	Long:  `Init an Azure Kubernetes Service (AKS) cluster configuration`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		return Init(environment)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
