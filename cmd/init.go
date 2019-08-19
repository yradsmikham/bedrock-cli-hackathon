package cmd

import (
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/docker/docker/pkg/namesgenerator"
	"github.com/kyokomi/emoji"
	log "github.com/sirupsen/logrus"
)

// SSHKey is the public key
var SSHKey string
var randomClusterName string

// Generates random cluster name if "--cluster-name" not specified
func nameGenerator() (name string) {
	rand.Seed(time.Now().UnixNano())
	randomClusterName = strings.Replace(namesgenerator.GetRandomName(0), "_", "-", -1)
	if len(randomClusterName) >= 30 {
		randomClusterName = nameGenerator()
	}
	return randomClusterName
}

// Init function initializes the configuration for a given environment
func Init(environment string, clusterName string) (cluster string, err error) {
	requiredSystemTools := []string{"git", "helm", "sh", "curl", "terraform", "az"}
	for _, tool := range requiredSystemTools {
		path, err := exec.LookPath(tool)
		if err != nil {
			return "", err
		}
		log.Info(emoji.Sprintf(":mag: Using %s: %s", tool, path))
	}

	// Check if Bedrock Repo is already cloned
	log.Info(emoji.Sprintf(":open_file_folder: Checking for Bedrock"))
	if output, err := exec.Command("git", "clone", "https://github.com/microsoft/bedrock").CombinedOutput(); output != nil || err != nil {
		log.Info(emoji.Sprintf(":star: Bedrock Repo already cloned"))
	}

	// If cluster name not provided, generate a random cluster name
	if clusterName == "" {
		randomClusterName := nameGenerator()
		log.Info(emoji.Sprintf(":space_invader: Bedrock Cluster Name: %s", randomClusterName))
		clusterName = randomClusterName
	}

	// Copy Terraform Template
	environmentPath := "bedrock/cluster/environments/" + clusterName
	if error := os.MkdirAll(environmentPath, os.ModePerm); error != nil {
		return "", error
	}

	log.Info(emoji.Sprintf(":flashlight: Creating New Environment %s", environmentPath))
	if output, err := exec.Command("cp", "-r", "bedrock/cluster/environments/"+environment, environmentPath).CombinedOutput(); err != nil {
		log.Error(emoji.Sprintf(":no_entry_sign: %s: %s", err, output))
		return "", err
	}

	// Generate SSH keys
	fullEnvironmentPath := environmentPath + "/" + environment
	if environment != COMMON {
		SSHKey, _ = SSH(fullEnvironmentPath, "deploy-key")
	}

	// Create bedrock-config.tfvars
	if err := addConfigTemplate(environment, fullEnvironmentPath, environmentPath, clusterName, SSHKey); err != nil {
		return "", err
	}
	return clusterName, err
}

// Get environment variables
func getEnvVariables(clusterName string) (err error) {
	if storageAccount == "" {
		storageAccount = os.Getenv("AZURE_STORAGE_ACCOUNT")
	}
	if accessKey == "" {
		accessKey = os.Getenv("AZURE_STORAGE_KEY")
	}
	if containerName == "" {
		containerName = os.Getenv("AZURE_CONTAINER")
	}
	if subscription == "" {
		subscription = os.Getenv("ARM_SUBSCRIPTION_ID")
	}
	if servicePrincipal == "" {
		servicePrincipal = os.Getenv("ARM_CLIENT_ID")
	}
	if tenant == "" {
		tenant = os.Getenv("ARM_TENANT_ID")
	}
	if secret == "" {
		secret = os.Getenv("ARM_CLIENT_SECRET")
	}
	if vnet == "" {
		vnet = clusterName + "-vnet"
	}
	if dnsPrefix == "" {
		dnsPrefix = clusterName
	}
	if keyvaultName == "" {
		keyvaultName = clusterName + "-kv"
	}
	if keyvaultRG == "" {
		keyvaultRG = clusterName + "-kv-rg"
	}
	return err
}

// Generate bedrock-config.tfvars (and bedrock-config.toml) and bedrock-backend-config.tfvars (if appropriate)
func generateTfvars(envPath string, envType string, clusterName string, sshKey string) (err error) {

	configMap := make(map[string]string)
	backendConfigMap := make(map[string]string)

	backendTfvarsFile, _ := os.Create(envPath + "/bedrock-backend-config.tfvars")
	tfvarsFile, _ := os.Create(envPath + "/bedrock-config.tfvars")
	configTomlFile, _ := os.Create(envPath + "/bedrock-config.toml")
	log.Info(emoji.Sprintf(":page_with_curl: Create Bedrock config file " + envPath + "/bedrock-config.tfvars"))

	// Supported environments
	if envType == SIMPLE {
		azureSimpleTemplate(configMap, clusterName, sshKey)
	}
	if envType == COMMON {
		backendTemplate(backendConfigMap, clusterName, COMMON)
		azureCommonInfraTemplate(configMap, clusterName, sshKey)
	}
	if envType == KEYVAULT {
		backendTemplate(backendConfigMap, clusterName, KEYVAULT)
		azureSingleKVTemplate(configMap, clusterName, sshKey)
	}
	if envType == MULTIPLE {
		azureMultipleTemplate(configMap, clusterName, sshKey)
	}

	// Iterate through backend config
	for setting, value := range backendConfigMap {
		if _, err := backendTfvarsFile.WriteString(setting + " = " + value + "\n"); err != nil {
			return err
		}
	}
	backendTfvarsFile.Close()

	// Iterate through normal config (terraform.tfvars)
	for setting, value := range configMap {
		if _, err := tfvarsFile.WriteString(setting + " = " + value + "\n"); err != nil {
			return err
		}
	}
	tfvarsFile.Close()

	// Generate the toml file will be used to extract environment variables via "viper"
	for setting, value := range configMap {
		if _, err := configTomlFile.WriteString(setting + " = " + value + "\n"); err != nil {
			return err
		}
	}
	configTomlFile.Close()

	return err
}

func backendTemplate(config map[string]string, clusterName string, env string) {
	config["storage_account_name"] = "\"" + storageAccount + "\""
	config["access_key"] = "\"" + accessKey + "\""
	config["container_name"] = "\"" + containerName + "\""
	config["key"] = "\"" + "tfstate-" + env + clusterName + "\""
}

func azureSimpleTemplate(config map[string]string, clusterName string, sshKey string) {
	config["resource_group_name"] = "\"" + clusterName + "-rg\""
	config["resource_group_location"] = "\"" + region + "\""
	config["cluster_name"] = "\"" + clusterName + "\""
	config["dns_prefix"] = "\"" + dnsPrefix + "\""
	config["service_principal_id"] = "\"" + servicePrincipal + "\""
	config["service_principal_secret"] = "\"" + secret + "\""
	config["ssh_public_key"] = "\"" + sshKey + "\""
	config["gitops_ssh_url"] = "\"" + gitopsSSHUrl + "\""
	config["gitops_ssh_key"] = "\"" + "deploy-key" + "\""
	config["vnet_name"] = "\"" + vnet + "\""
	config["agent_vm_count"] = "\"" + vmCount + "\""
	config["gitops_poll_interval"] = "\"" + gitopsPollInterval + "\""
	config["gitops_url_branch"] = "\"" + gitopsURLBranch + "\""
	config["gitops_path"] = "\"" + gitopsPath + "\""
}

func azureCommonInfraTemplate(config map[string]string, clusterName string, sshKey string) {
	config["global_resource_group_name"] = "\"" + keyvaultRG + "\""
	config["global_resource_group_location"] = "\"" + region + "\""
	config["keyvault_name"] = "\"" + keyvaultName + "\""
	config["service_principal_id"] = "\"" + servicePrincipal + "\""
	config["tenant_id"] = "\"" + tenant + "\""
	config["address_space"] = "\"" + addressSpace + "\""
	config["subnet_prefix"] = "\"" + subnetPrefix + "\""
	config["subnet_name"] = "\"" + clusterName + "-subnet\""
	config["vnet_name"] = "\"" + clusterName + "-vnet\""
	config["subscription"] = "\"" + subscription + "\""
	config["secret"] = "\"" + secret + "\""
}

func azureSingleKVTemplate(config map[string]string, clusterName string, sshKey string) {
	config["resource_group_name"] = "\"" + clusterName + "-rg\""
	config["resource_group_location"] = "\"" + region + "\""
	config["cluster_name"] = "\"" + clusterName + "\""
	config["agent_vm_size"] = "\"" + vmSize + "\""
	config["service_principal_id"] = "\"" + servicePrincipal + "\""
	config["service_principal_secret"] = "\"" + secret + "\""
	config["ssh_public_key"] = "\"" + sshKey + "\""
	config["gitops_ssh_url"] = "\"" + gitopsSSHUrl + "\""
	config["gitops_ssh_key"] = "\"" + "deploy-key" + "\""
	config["keyvault_resource_group"] = "\"" + keyvaultRG + "\""
	config["keyvault_name"] = "\"" + keyvaultName + "\""
	config["vnet_subnet_id"] = "\"/subscriptions/" + subscription + "/resourceGroups/" + keyvaultRG + "/providers/Microsoft.Network/virtualNetworks/" + clusterName + "-vnet/subnets/" + clusterName + "-subnet" + "\""
	config["agent_vm_count"] = "\"" + vmCount + "\""
	config["gitops_poll_interval"] = "\"" + gitopsPollInterval + "\""
	config["gitops_url_branch"] = "\"" + gitopsURLBranch + "\""
	config["gitops_path"] = "\"" + gitopsPath + "\""
	config["dns_prefix"] = "\"" + dnsPrefix + "\""
	config["address_space"] = "\"" + addressSpace + "\""
	config["subnet_prefixes"] = "\"" + subnetPrefix + "\""
}

func azureMultipleTemplate(config map[string]string, clusterName string, sshKey string) {
	config["agent_vm_count"] = "\"" + "3" + "\""
	config["agent_vm_size"] = "\"" + "Standard_D4s_v3" + "\""
	config["cluster_name"] = "\"" + clusterName + "\""
	config["dns_prefix"] = "\"" + clusterName + "\""
	config["keyvault_resource_group"] = "\"" + keyvaultRG + "\""
	config["keyvault_name"] = "\"" + keyvaultName + "\""
	config["service_principal_id"] = "\"" + servicePrincipal + "\""
	config["service_principal_secret"] = "\"" + secret + "\""
	config["ssh_public_key"] = "\"" + sshKey + "\""
	config["gitops_ssh_url"] = "\"" + gitopsSSHUrl + "\""
	config["gitops_ssh_key"] = "\"" + "deploy-key" + "\""
	config["traffic_manager_profile_name"] = "\"" + clusterName + "-tm\""
	config["traffic_manager_dns_name"] = "\"" + clusterName + "-tm\""
	config["traffic_manager_resource_group_name"] = "\"" + clusterName + "-tm-rg\""
	config["traffic_manager_resource_group_location"] = "\"" + "westus2" + "\""
	config["west_resource_group_name"] = "\"" + clusterName + "-west-rg\""
	config["west_resource_group_location"] = "\"" + "westus2" + "\""
	config["gitops_west_path"] = "\"" + gitopsPathWest + "\""
	config["east_resource_group_name"] = "\"" + clusterName + "-east-rg\""
	config["east_resource_group_location"] = "\"" + "eastus" + "\""
	config["gitops_east_path"] = "\"" + gitopsPathEast + "\""
	config["central_resource_group_name"] = "\"" + clusterName + "-central-rg\""
	config["central_resource_group_location"] = "\"" + "centralus" + "\""
	config["gitops_central_path"] = "\"" + gitopsPathCentral + "\""
	config["gitops_central_url_branch"] = "\"" + gitopsURLBranchCentral + "\""
	config["gitops_east_url_branch"] = "\"" + gitopsURLBranchEast + "\""
	config["gitops_west_url_branch"] = "\"" + gitopsURLBranchWest + "\""
}

// Adds a blank bedrock config template
func addConfigTemplate(environment string, fullEnvironmentPath string, environmentPath string, clusterName string, sshKey string) (err error) {
	sshKey = strings.TrimSuffix(sshKey, "\n")

	if environment == SIMPLE {

		if error := getEnvVariables(clusterName); error != nil {
			return error
		}
		if error := generateTfvars(fullEnvironmentPath, SIMPLE, clusterName, sshKey); error != nil {
			return error
		}

		log.Info(emoji.Sprintf(":raised_hands: Azure Simple cluster environment " + fullEnvironmentPath + " has been successfully created!"))
		log.Info(emoji.Sprintf(":white_check_mark: To proceed, run 'bedrock simulate " + environmentPath + "'"))

		return err
	}

	if environment == COMMON {

		if error := getEnvVariables(clusterName); error != nil {
			return error
		}
		if error := generateTfvars(fullEnvironmentPath, COMMON, clusterName, sshKey); error != nil {
			return error
		}

		commonInfraPath = fullEnvironmentPath

		log.Info(emoji.Sprintf(":raised_hands: Azure Common Infra environment " + fullEnvironmentPath + " has been successfully created!"))

		return err
	}

	if environment == KEYVAULT {

		// When common infra is a dependency but does not exist, create one
		if commonInfraPath == "" {
			log.Info(emoji.Sprintf(":two_men_holding_hands: Common Infra path is not set, creating common infra with tenant id %s", tenant))
			if _, error := Init(COMMON, clusterName); error != nil {
				return error
			}
		}
		log.Info(emoji.Sprintf(":family: Common Infra path is set to %s", commonInfraPath))

		if error := getEnvVariables(clusterName); error != nil {
			return error
		}
		if error := generateTfvars(fullEnvironmentPath, KEYVAULT, clusterName, sshKey); error != nil {
			return error
		}

		log.Info(emoji.Sprintf(":raised_hands: Azure Single Keyvault cluster environment " + fullEnvironmentPath + " has been successfully created!"))
		log.Info(emoji.Sprintf(":white_check_mark: To proceed, run 'bedrock simulate " + environmentPath + "'"))

		return err
	}

	if environment == MULTIPLE {

		// When keyvault is not specified and common infra does not exist, create one
		if keyvaultName == "" && keyvaultRG == "" && commonInfraPath == "" {
			log.Info(emoji.Sprintf(":two_men_holding_hands: Common Infra path is not set, creating common infra with tenant id %s", tenant))
			if _, error := Init(COMMON, clusterName); error != nil {
				return error
			}
		}
		log.Info(emoji.Sprintf(":family: Common Infra path is set to %s", commonInfraPath))

		if error := getEnvVariables(clusterName); error != nil {
			return error
		}
		if error := generateTfvars(fullEnvironmentPath, MULTIPLE, clusterName, sshKey); error != nil {
			return error
		}

		log.Info(emoji.Sprintf(":raised_hands: Azure Multiple cluster environment " + fullEnvironmentPath + " has been successfully created!"))
		log.Info(emoji.Sprintf(":white_check_mark: To proceed, run 'bedrock simulate " + environmentPath + "'"))

		return err
	}
	return err
}
