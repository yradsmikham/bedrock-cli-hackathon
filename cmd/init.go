package cmd

import (
	"bufio"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/docker/docker/pkg/namesgenerator"
	"github.com/kyokomi/emoji"
	log "github.com/sirupsen/logrus"
)

var randomClusterName string

// Init functions initializes the configuration for the given environment
func Init(environment string, clusterName string) (err error) {
	requiredSystemTools := []string{"git", "helm", "sh", "curl", "terraform", "az"}
	for _, tool := range requiredSystemTools {
		path, err := exec.LookPath(tool)
		if err != nil {
			return err
		}
		log.Info(emoji.Sprintf(":mag: Using %s: %s", tool, path))
	}

	// If cluster name not provided, generate a random cluster name
	if clusterName == "" {
		rand.Seed(time.Now().UnixNano())
		randomClusterName = strings.Replace(namesgenerator.GetRandomName(0), "_", "-", -1)
		clusterName = randomClusterName
	}

	// Check if Bedrock Repo is already cloned
	log.Info(emoji.Sprintf(":open_file_folder: Checking for Bedrock"))
	if output, err := exec.Command("git", "clone", "https://github.com/microsoft/bedrock").CombinedOutput(); output != nil || err != nil {
		log.Info(emoji.Sprintf(":star: Bedrock Repo already cloned"))
	}

	// Copy Terraform Template
	environmentPath := "bedrock/cluster/environments/" + clusterName
	if error := os.MkdirAll(environmentPath, os.ModePerm); error != nil {
		return error
	}

	log.Info(emoji.Sprintf(":flashlight: Creating New Environment %s", environmentPath))
	if output, err := exec.Command("cp", "-r", "bedrock/cluster/environments/"+environment, environmentPath).CombinedOutput(); err != nil {
		log.Error(emoji.Sprintf(":no_entry_sign: %s: %s", err, output))
		return err
	}

	// Generate SSH keys
	fullEnvironmentPath := environmentPath + "/" + environment
	SSHKey := ""
	if environment != COMMON {
		SSHKey, _ = SSH(fullEnvironmentPath, "deploy-key")
	}
	// Create bedrock-config.tfvars
	if err = addConfigTemplate(environment, fullEnvironmentPath, environmentPath, clusterName, SSHKey); err != nil {
		return err
	}
	return err
}

func copyCommonInfraTemplateToPath(commonInfraPath string, fullEnvironmentPath string, environmentPath string, environment string, config map[string]string) (err error) {
	filename := commonInfraPath + "/bedrock-config.tfvars"
	log.Info(emoji.Sprintf(":hushed: Copying %s variables from %s", COMMON, filename))

	if len(filename) == 0 {
		return err
	}
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if equal := strings.Index(line, "="); equal >= 0 {
			if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
				value := ""
				if len(line) > equal {
					value = strings.TrimSpace(line[equal+1:])
				}
				config[key] = value
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return err
	}

	destinationPath := strings.Replace(fullEnvironmentPath, environment, "", -1)
	originPath := strings.Replace(commonInfraPath, COMMON, "", -1)
	log.Info(emoji.Sprintf(":books: Copying %s template to environment directory", COMMON))
	if output, err := exec.Command("cp", "-r", originPath, destinationPath).CombinedOutput(); err != nil {
		log.Error(emoji.Sprintf(":no_entry_sign: %s: %s", err, output))
		return err
	}

	return err
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

// Adds a blank bedrock config template
func addConfigTemplate(environment string, fullEnvironmentPath string, environmentPath string, clusterName string, SSHKey string) (err error) {
	SSHKey = strings.TrimSuffix(SSHKey, "\n")

	if environment == SIMPLE {
		azureSimpleConfig := make(map[string]string)

		if error := getEnvVariables(clusterName); error != nil {
			return error
		}

		azureSimpleConfig["resource_group_name"] = "\"" + clusterName + "-rg\""
		azureSimpleConfig["resource_group_location"] = "\"" + region + "\""
		azureSimpleConfig["cluster_name"] = "\"" + clusterName + "\""
		azureSimpleConfig["dns_prefix"] = "\"" + dnsPrefix + "\""
		azureSimpleConfig["service_principal_id"] = "\"" + servicePrincipal + "\""
		azureSimpleConfig["service_principal_secret"] = "\"" + secret + "\""
		azureSimpleConfig["ssh_public_key"] = "\"" + SSHKey + "\""
		azureSimpleConfig["gitops_ssh_url"] = "\"" + gitopsSSHUrl + "\""
		azureSimpleConfig["gitops_ssh_key"] = "\"" + "deploy-key" + "\""
		azureSimpleConfig["vnet_name"] = "\"" + vnet + "\""
		azureSimpleConfig["agent_vm_count"] = "\"" + vmCount + "\""
		azureSimpleConfig["gitops_poll_interval"] = "\"" + gitopsPollInterval + "\""
		azureSimpleConfig["gitops_url_branch"] = "\"" + gitopsURLBranch + "\""
		azureSimpleConfig["gitops_path"] = "\"" + gitopsPath + "\""

		// Generate bedrock-config.tfvars
		f, err := os.Create(fullEnvironmentPath + "/bedrock-config.tfvars")
		log.Info(emoji.Sprintf(":page_with_curl: Create Bedrock config file " + fullEnvironmentPath + "/bedrock-config.tfvars"))
		if err != nil {
			return err
		}

		for setting, value := range azureSimpleConfig {
			if _, err := f.WriteString(setting + " = " + value + "\n"); err != nil {
				return err
			}
		}

		f.Close()

		log.Info(emoji.Sprintf(":raised_hands: Azure Simple cluster environment " + fullEnvironmentPath + " has been successfully created!"))
		log.Info(emoji.Sprintf(":white_check_mark: To proceed, run 'bedrock simulate " + environmentPath + "'"))

		return err
	}

	if environment == COMMON {
		commonInfraConfig := make(map[string]string)

		if error := getEnvVariables(clusterName); error != nil {
			return error
		}

		commonInfraConfig["global_resource_group_name"] = "\"" + keyvaultRG + "\""
		commonInfraConfig["global_resource_group_location"] = "\"" + region + "\""
		commonInfraConfig["keyvault_name"] = "\"" + keyvaultName + "\""
		commonInfraConfig["service_principal_id"] = "\"" + servicePrincipal + "\""
		commonInfraConfig["tenant_id"] = "\"" + tenant + "\""
		commonInfraConfig["address_space"] = "\"" + addressSpace + "\""
		commonInfraConfig["subnet_prefix"] = "\"" + subnetPrefix + "\""
		commonInfraConfig["subnet_name"] = "\"" + clusterName + "-subnet\""
		commonInfraConfig["vnet_name"] = "\"" + clusterName + "-vnet\""
		commonInfraConfig["subscription"] = "\"" + subscription + "\""
		commonInfraConfig["secret"] = "\"" + secret + "\""

		// Generate bedrock-config.tfvars
		f, err := os.Create(fullEnvironmentPath + "/bedrock-config.tfvars")
		log.Info(emoji.Sprintf(":page_with_curl: Create Bedrock config file " + fullEnvironmentPath + "/bedrock-config.tfvars"))
		if err != nil {
			return err
		}

		for setting, value := range commonInfraConfig {
			if _, err := f.WriteString(setting + " = " + value + "\n"); err != nil {
				return err
			}
		}

		f.Close()

		configFile, err := os.Create(fullEnvironmentPath + "/bedrock-config.toml")
		log.Info(emoji.Sprintf(":page_with_curl: Create Bedrock config file " + fullEnvironmentPath + "/bedrock-config.toml"))
		if err != nil {
			return err
		}

		for setting, value := range commonInfraConfig {
			if _, err := configFile.WriteString(setting + " = " + value + "\n"); err != nil {
				return err
			}
		}

		configFile.Close()

		commonInfraBackendConfig := make(map[string]string)

		commonInfraBackendConfig["storage_account_name"] = "\"" + storageAccount + "\""
		commonInfraBackendConfig["access_key"] = "\"" + accessKey + "\""
		commonInfraBackendConfig["container_name"] = "\"" + containerName + "\""
		commonInfraBackendConfig["key"] = "\"" + "tfstate-common-infra-" + clusterName + "\""

		// Generate backend-config.tfvars
		backendFile, err := os.Create(fullEnvironmentPath + "/bedrock-backend-config.tfvars")
		log.Info(emoji.Sprintf(":page_with_curl: Create Bedrock backend config file " + fullEnvironmentPath + "/bedrock-backend-config.tfvars"))
		if err != nil {
			return err
		}

		for setting, value := range commonInfraBackendConfig {
			if _, err := backendFile.WriteString(setting + " = " + value + "\n"); err != nil {
				return err
			}
		}

		backendFile.Close()

		commonInfraPath = fullEnvironmentPath

		return err
	}

	if environment == KEYVAULT {

		singleKeyvaultConfig := make(map[string]string)

		// When common infra is not initialized, create one
		if commonInfraPath == "" {
			log.Info(emoji.Sprintf(":two_men_holding_hands: Common Infra path is not set, creating common infra with tenant id %s", tenant))
			if error := Init(COMMON, clusterName); error != nil {
				return error
			}
		}

		if error := getEnvVariables(clusterName); error != nil {
			return error
		}

		log.Info(emoji.Sprintf(":family: Common Infra path is set to %s", commonInfraPath))

		if clusterName == "" {
			if error := copyCommonInfraTemplateToPath(commonInfraPath, fullEnvironmentPath, environmentPath, environment, singleKeyvaultConfig); error != nil {
				return error
			}
		}

		singleKeyvaultConfig["resource_group_name"] = "\"" + clusterName + "-rg\""
		singleKeyvaultConfig["resource_group_location"] = "\"" + region + "\""
		singleKeyvaultConfig["cluster_name"] = "\"" + clusterName + "\""
		singleKeyvaultConfig["agent_vm_size"] = "\"" + vmSize + "\""
		singleKeyvaultConfig["service_principal_id"] = "\"" + servicePrincipal + "\""
		singleKeyvaultConfig["service_principal_secret"] = "\"" + secret + "\""
		singleKeyvaultConfig["ssh_public_key"] = "\"" + SSHKey + "\""
		singleKeyvaultConfig["gitops_ssh_url"] = "\"" + gitopsSSHUrl + "\""
		singleKeyvaultConfig["gitops_ssh_key"] = "\"" + "deploy-key" + "\""
		singleKeyvaultConfig["keyvault_resource_group"] = "\"" + keyvaultRG + "\""
		singleKeyvaultConfig["keyvault_name"] = "\"" + keyvaultName + "\""
		singleKeyvaultConfig["vnet_subnet_id"] = "\"/subscriptions/" + subscription + "/resourceGroups/" + keyvaultRG + "/providers/Microsoft.Network/virtualNetworks/" + clusterName + "-vnet/subnets/" + clusterName + "-subnet" + "\""
		singleKeyvaultConfig["agent_vm_count"] = "\"" + vmCount + "\""
		singleKeyvaultConfig["gitops_poll_interval"] = "\"" + gitopsPollInterval + "\""
		singleKeyvaultConfig["gitops_url_branch"] = "\"" + gitopsURLBranch + "\""
		singleKeyvaultConfig["gitops_path"] = "\"" + gitopsPath + "\""
		singleKeyvaultConfig["dns_prefix"] = "\"" + dnsPrefix + "\""
		singleKeyvaultConfig["address_space"] = "\"" + addressSpace + "\""
		singleKeyvaultConfig["subnet_prefixes"] = "\"" + subnetPrefix + "\""

		f, err := os.Create(fullEnvironmentPath + "/bedrock-config.tfvars")
		log.Info(emoji.Sprintf(":page_with_curl: Create Bedrock config file " + environmentPath + "/bedrock-config.tfvars"))
		if err != nil {
			return err
		}

		for setting, value := range singleKeyvaultConfig {
			if _, error := f.WriteString(setting + " = " + value + "\n"); error != nil {
				return error
			}
		}

		f.Close()

		// Generate backend-config.tfvars
		singleKeyvaultBackendConfig := make(map[string]string)

		singleKeyvaultBackendConfig["storage_account_name"] = "\"" + storageAccount + "\""
		singleKeyvaultBackendConfig["access_key"] = "\"" + accessKey + "\""
		singleKeyvaultBackendConfig["container_name"] = "\"" + containerName + "\""
		singleKeyvaultBackendConfig["key"] = "\"" + "tfstate-single-keyvault-" + clusterName + "\""

		kvBackendFile, err := os.Create(fullEnvironmentPath + "/bedrock-backend-config.tfvars")
		log.Info(emoji.Sprintf(":page_with_curl: Create Bedrock backend config file " + fullEnvironmentPath + "/bedrock-backend-config.tfvars"))
		if err != nil {
			return err
		}

		for setting, value := range singleKeyvaultBackendConfig {
			if _, error := kvBackendFile.WriteString(setting + " = " + value + "\n"); error != nil {
				return error
			}
		}

		kvBackendFile.Close()

		log.Info(emoji.Sprintf(":raised_hands: Azure Single Keyvault cluster environment " + fullEnvironmentPath + " has been successfully created!"))
		log.Info(emoji.Sprintf(":white_check_mark: To proceed, run 'bedrock simulate " + environmentPath + "'"))

		return err
	}

	if environment == MULTIPLE {

		multipleConfig := make(map[string]string)

		// When keyvault is not specified and common infra is not initialized, create one
		if keyvaultName == "" && keyvaultRG == "" && commonInfraPath == "" {
			log.Info(emoji.Sprintf(":two_men_holding_hands: Common Infra path is not set, creating common infra with tenant id %s", tenant))
			if error := Init(COMMON, clusterName); error != nil {
				return error
			}
		}

		if error := getEnvVariables(clusterName); error != nil {
			return error
		}

		log.Info(emoji.Sprintf(":family: Common Infra path is set to %s", commonInfraPath))

		if clusterName == "" {
			if error := copyCommonInfraTemplateToPath(commonInfraPath, fullEnvironmentPath, environmentPath, environment, multipleConfig); error != nil {
				return error
			}
		}

		multipleConfig["agent_vm_count"] = "\"" + "3" + "\""
		multipleConfig["agent_vm_size"] = "\"" + "Standard_D4s_v3" + "\""
		multipleConfig["cluster_name"] = "\"" + clusterName + "\""
		multipleConfig["dns_prefix"] = "\"" + clusterName + "\""
		multipleConfig["keyvault_resource_group"] = "\"" + keyvaultRG + "\""
		multipleConfig["keyvault_name"] = "\"" + keyvaultName + "\""
		multipleConfig["service_principal_id"] = "\"" + servicePrincipal + "\""
		multipleConfig["service_principal_secret"] = "\"" + secret + "\""
		multipleConfig["ssh_public_key"] = "\"" + SSHKey + "\""
		multipleConfig["gitops_ssh_url"] = "\"" + gitopsSSHUrl + "\""
		multipleConfig["gitops_ssh_key"] = "\"" + "deploy-key" + "\""
		multipleConfig["traffic_manager_profile_name"] = "\"" + clusterName + "-tm\""
		multipleConfig["traffic_manager_dns_name"] = "\"" + clusterName + "-tm\""
		multipleConfig["traffic_manager_resource_group_name"] = "\"" + clusterName + "-tm-rg\""
		multipleConfig["traffic_manager_resource_group_location"] = "\"" + "westus2" + "\""
		multipleConfig["west_resource_group_name"] = "\"" + clusterName + "-west-rg\""
		multipleConfig["west_resource_group_location"] = "\"" + "westus2" + "\""
		multipleConfig["gitops_west_path"] = "\"" + gitopsPathWest + "\""
		multipleConfig["east_resource_group_name"] = "\"" + clusterName + "-east-rg\""
		multipleConfig["east_resource_group_location"] = "\"" + "eastus" + "\""
		multipleConfig["gitops_east_path"] = "\"" + gitopsPathEast + "\""
		multipleConfig["central_resource_group_name"] = "\"" + clusterName + "-central-rg\""
		multipleConfig["central_resource_group_location"] = "\"" + "centralus" + "\""
		multipleConfig["gitops_central_path"] = "\"" + gitopsPathCentral + "\""
		multipleConfig["gitops_central_url_branch"] = "\"" + gitopsURLBranchCentral + "\""
		multipleConfig["gitops_east_url_branch"] = "\"" + gitopsURLBranchEast + "\""
		multipleConfig["gitops_west_url_branch"] = "\"" + gitopsURLBranchWest + "\""

		// Generate bedrock-config.tfvars
		f, err := os.Create(fullEnvironmentPath + "/bedrock-config.tfvars")
		log.Info(emoji.Sprintf(":page_with_curl: Create Bedrock config file " + fullEnvironmentPath + "/bedrock-config.tfvars"))
		if err != nil {
			return err
		}

		for setting, value := range multipleConfig {
			if _, err := f.WriteString(setting + " = " + value + "\n"); err != nil {
				return err
			}
		}

		f.Close()

		configFile, err := os.Create(fullEnvironmentPath + "/bedrock-config.toml")
		log.Info(emoji.Sprintf(":page_with_curl: Create Bedrock config file " + fullEnvironmentPath + "/bedrock-config.toml"))
		if err != nil {
			return err
		}

		for setting, value := range multipleConfig {
			if _, err := configFile.WriteString(setting + " = " + value + "\n"); err != nil {
				return err
			}
		}

		configFile.Close()

		log.Info(emoji.Sprintf(":raised_hands: Azure Multiple cluster environment " + fullEnvironmentPath + " has been successfully created!"))
		log.Info(emoji.Sprintf(":white_check_mark: To proceed, run 'bedrock simulate " + environmentPath + "'"))

		return err
	}
	return err
}
