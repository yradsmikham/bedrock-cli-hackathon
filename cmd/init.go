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
func Init(environment string, clusterName string) (configPath string, err error) {
	requiredSystemTools := []string{"git", "helm", "sh", "curl", "terraform", "az"}
	for _, tool := range requiredSystemTools {
		path, err := exec.LookPath(tool)
		if err != nil {
			return "", err
		}
		log.Info(emoji.Sprintf(":mag: Using %s: %s", tool, path))
	}

	if clusterName == "" {
		rand.Seed(time.Now().UnixNano())
		randomClusterName = strings.Replace(namesgenerator.GetRandomName(0), "_", "-", -1)
	} else {
		randomClusterName = clusterName
	}

	// Check if Bedrock Repo is already cloned
	log.Info(emoji.Sprintf(":open_file_folder: Checking for Bedrock"))
	if output, err := exec.Command("git", "clone", "https://github.com/microsoft/bedrock").CombinedOutput(); output != nil || err != nil {
		log.Info(emoji.Sprintf(":star: Bedrock Repo already cloned"))
	}

	// Copy Terraform Template
	environmentPath := "bedrock/cluster/environments/" + randomClusterName
	os.MkdirAll(environmentPath, os.ModePerm)

	log.Info(emoji.Sprintf(":flashlight: Creating New Environment %s", environmentPath))
	if output, err := exec.Command("cp", "-r", "bedrock/cluster/environments/"+environment, environmentPath).CombinedOutput(); err != nil {
		log.Error(emoji.Sprintf(":no_entry_sign: %s: %s", err, output))
		return "", err
	}

	// Generate ssh keys
	fullEnvironmentPath := environmentPath + "/" + environment
	SSHKey := ""
	if environment != COMMON {
		SSHKey, _ = SSH(fullEnvironmentPath, "deploy-key")
	}
	if err == nil {
		// Save bedrock-config.tfvars
		err = addConfigTemplate(environment, fullEnvironmentPath, environmentPath, randomClusterName, SSHKey)

		if err == nil {
			return "", nil
		}
	}
	configPath = fullEnvironmentPath + "/bedrock-config.tfvars"
	return
}

func copyCommonInfraTemplateToPath(commonInfraPath string, fullEnvironmentPath string, environmentPath string, environment string, config map[string]string) (err error) {
	filename := commonInfraPath + "/bedrock-config.tfvars"
	log.Info(emoji.Sprintf(":hushed: Copying %s variables from %s", COMMON, filename))

	if len(filename) == 0 {
		return nil
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

// Adds a blank bedrock config template
func addConfigTemplate(environment string, fullEnvironmentPath string, environmentPath string, clusterName string, SSHKey string) (err error) {
	SSHKey = strings.TrimSuffix(SSHKey, "\n")

	if environment == SIMPLE {
		azureSimpleConfig := make(map[string]string)

		azureSimpleConfig["resource_group_name"] = "\"" + clusterName + "-rg\""
		azureSimpleConfig["resource_group_location"] = "\"\""
		azureSimpleConfig["cluster_name"] = "\"" + clusterName + "\""
		azureSimpleConfig["agent_vm_count"] = "\"\""
		azureSimpleConfig["dns_prefix"] = "\"" + clusterName + "\""
		azureSimpleConfig["service_principal_id"] = "\"" + servicePrincipal + "\""
		azureSimpleConfig["service_principal_secret"] = "\"" + secret + "\""
		azureSimpleConfig["ssh_public_key"] = "\"" + SSHKey + "\""
		azureSimpleConfig["gitops_ssh_url"] = "\"" + gitopsSSHUrl + "\""
		azureSimpleConfig["gitops_ssh_key"] = "\"" + "deploy-key" + "\""
		azureSimpleConfig["vnet_name"] = "\"" + clusterName + "-vnet\""
		azureSimpleConfig["agent_vm_count"] = "\"" + "3" + "\""
		azureSimpleConfig["resource_group_location"] = "\"" + "westus2" + "\""

		f, err := os.Create(fullEnvironmentPath + "/bedrock-config.tfvars")
		log.Info(emoji.Sprintf(":page_with_curl: Create Bedrock config file " + fullEnvironmentPath + "/bedrock-config.tfvars"))
		if err != nil {
			return err
		}

		for setting, value := range azureSimpleConfig {
			f.WriteString(setting + " = " + value + "\n")
		}

		f.Close()

		log.Info(emoji.Sprintf(":raised_hands: Azure Simple cluster environment " + fullEnvironmentPath + " has been successfully created!"))
		log.Info(emoji.Sprintf(":white_check_mark: To proceed, run 'bedrock simulate " + environmentPath + "'"))

		return nil
	}

	if environment == COMMON {
		commonInfraConfig := make(map[string]string)

		commonInfraConfig["global_resource_group_name"] = "\"" + clusterName + "-rg\""
		commonInfraConfig["global_resource_group_location"] = "\"" + "westus2" + "\""
		commonInfraConfig["keyvault_name"] = "\"" + clusterName + "-kv\""
		commonInfraConfig["service_principal_id"] = "\"" + servicePrincipal + "\""
		commonInfraConfig["tenant_id"] = "\"" + tenant + "\""
		commonInfraConfig["address_space"] = "\"" + "10.39.0.0/16" + "\""
		commonInfraConfig["subnet_prefix"] = "\"" + "10.39.0.0/24" + "\""
		commonInfraConfig["subnet_name"] = "\"" + clusterName + "-subnet\""
		commonInfraConfig["vnet_name"] = "\"" + clusterName + "-vnet\""
		commonInfraConfig["subscription"] = "\"" + subscription + "\""
		commonInfraConfig["secret"] = "\"" + secret + "\""

		f, err := os.Create(fullEnvironmentPath + "/bedrock-config.tfvars")
		log.Info(emoji.Sprintf(":page_with_curl: Create Bedrock config file " + fullEnvironmentPath + "/bedrock-config.tfvars"))
		if err != nil {
			return err
		}

		for setting, value := range commonInfraConfig {
			f.WriteString(setting + " = " + value + "\n")
		}

		f.Close()

		configFile, err := os.Create(fullEnvironmentPath + "/bedrock-config.toml")
		log.Info(emoji.Sprintf(":page_with_curl: Create Bedrock config file " + fullEnvironmentPath + "/bedrock-config.toml"))
		if err != nil {
			return err
		}

		for setting, value := range commonInfraConfig {
			configFile.WriteString(setting + " = " + value + "\n")
		}

		configFile.Close()

		commonInfraBackendConfig := make(map[string]string)

		commonInfraBackendConfig["storage_account_name"] = "\"" + storageAccount + "\""
		commonInfraBackendConfig["access_key"] = "\"" + accessKey + "\""
		commonInfraBackendConfig["container_name"] = "\"" + containerName + "\""
		commonInfraBackendConfig["key"] = "\"" + "tfstate-common-infra-" + clusterName + "\""

		backendFile, err := os.Create(fullEnvironmentPath + "/bedrock-backend-config.tfvars")
		log.Info(emoji.Sprintf(":page_with_curl: Create Bedrock backend config file " + fullEnvironmentPath + "/bedrock-backend-config.tfvars"))
		if err != nil {
			return err
		}

		for setting, value := range commonInfraBackendConfig {
			backendFile.WriteString(setting + " = " + value + "\n")
		}

		backendFile.Close()

		commonInfraPath = fullEnvironmentPath

		return nil
	}

	if environment == KEYVAULT {

		singleKeyvaultConfig := make(map[string]string)

		// When common infra is not initialized, create one
		if commonInfraPath == "" {
			log.Info(emoji.Sprintf(":two_men_holding_hands: Common Infra path is not set, creating common infra with tenant id %s", tenant))
			Init(COMMON, clusterName)
		}

		log.Info(emoji.Sprintf(":family: Common Infra path is set to %s", commonInfraPath))

		if clusterName == "" {
			copyCommonInfraTemplateToPath(commonInfraPath, fullEnvironmentPath, environmentPath, environment, singleKeyvaultConfig)
		}

		singleKeyvaultConfig["resource_group_name"] = "\"" + clusterName + "-rg\""
		singleKeyvaultConfig["resource_group_location"] = "\"" + "westus2" + "\""
		singleKeyvaultConfig["cluster_name"] = "\"" + clusterName + "\""
		singleKeyvaultConfig["agent_vm_count"] = "\"" + "3" + "\""
		singleKeyvaultConfig["agent_vm_size"] = "\"Standard_D4s_v3\""
		singleKeyvaultConfig["service_principal_id"] = "\"" + servicePrincipal + "\""
		singleKeyvaultConfig["service_principal_secret"] = "\"" + secret + "\""
		singleKeyvaultConfig["ssh_public_key"] = "\"" + SSHKey + "\""
		singleKeyvaultConfig["gitops_ssh_url"] = "\"" + gitopsSSHUrl + "\""
		singleKeyvaultConfig["gitops_ssh_key"] = "\"" + "deploy-key" + "\""
		singleKeyvaultConfig["keyvault_resource_group"] = singleKeyvaultConfig["global_resource_group_name"]
		singleKeyvaultConfig["subnet_prefixes"] = singleKeyvaultConfig["subnet_prefix"]
		singleKeyvaultConfig["vnet_subnet_id"] = "\"/subscriptions/" + subscription + "/resourceGroups/" + strings.Replace(singleKeyvaultConfig["global_resource_group_name"], "\"", "", -1) + "/providers/Microsoft.Network/virtualNetworks/" + strings.Replace(singleKeyvaultConfig["vnet_name"], "\"", "", -1) + "/subnets/" + strings.Replace(singleKeyvaultConfig["subnet_name"], "\"", "", -1) + "\""

		f, err := os.Create(fullEnvironmentPath + "/bedrock-config.tfvars")
		log.Info(emoji.Sprintf(":page_with_curl: Create Bedrock config file " + environmentPath + "/bedrock-config.tfvars"))
		if err != nil {
			return err
		}

		for setting, value := range singleKeyvaultConfig {
			f.WriteString(setting + " = " + value + "\n")
		}

		f.Close()

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
			kvBackendFile.WriteString(setting + " = " + value + "\n")
		}

		kvBackendFile.Close()

		log.Info(emoji.Sprintf(":raised_hands: Azure Single Keyvault cluster environment " + fullEnvironmentPath + " has been successfully created!"))
		log.Info(emoji.Sprintf(":white_check_mark: To proceed, run 'bedrock simulate " + environmentPath + "'"))

		return nil
	}

	if environment == MULTIPLE {

		multipleConfig := make(map[string]string)

		// When common infra is not initialized, create one
		/* if commonInfraPath == "" {
			log.Info(emoji.Sprintf(":two_men_holding_hands: Common Infra path is not set, creating common infra with tenant id %s", tenant))
			Init(COMMON)
		}

		log.Info(emoji.Sprintf(":family: Common Infra path is set to %s", commonInfraPath))
		copyCommonInfraTemplateToPath(commonInfraPath, fullEnvironmentPath, environmentPath, environment, multipleConfig)
		*/

		multipleConfig["agent_vm_count"] = "\"" + "3" + "\""
		multipleConfig["agent_vm_size"] = "\"" + "Standard_D4s_v3" + "\""
		multipleConfig["cluster_name"] = "\"" + clusterName + "\""
		multipleConfig["dns_prefix"] = "\"" + clusterName + "\""
		multipleConfig["keyvault_resource_group"] = "\"" + clusterName + "-kv-rg\""
		multipleConfig["keyvault_name"] = "\"" + clusterName + "-kv\""
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
		multipleConfig["gitops_west_path"] = "\"\""
		multipleConfig["east_resource_group_name"] = "\"" + clusterName + "-east-rg\""
		multipleConfig["east_resource_group_location"] = "\"" + "eastus" + "\""
		multipleConfig["gitops_east_path"] = "\"\""
		multipleConfig["central_resource_group_name"] = "\"" + clusterName + "-central-rg\""
		multipleConfig["central_resource_group_location"] = "\"" + "centralus" + "\""
		multipleConfig["gitops_central_path"] = "\"\""
		multipleConfig["subscription"] = "\"" + subscription + "\""

		f, err := os.Create(fullEnvironmentPath + "/bedrock-config.tfvars")
		log.Info(emoji.Sprintf(":page_with_curl: Create Bedrock config file " + fullEnvironmentPath + "/bedrock-config.tfvars"))
		if err != nil {
			return err
		}

		for setting, value := range multipleConfig {
			f.WriteString(setting + " = " + value + "\n")
		}

		f.Close()

		configFile, err := os.Create(fullEnvironmentPath + "/bedrock-config.toml")
		log.Info(emoji.Sprintf(":page_with_curl: Create Bedrock config file " + fullEnvironmentPath + "/bedrock-config.toml"))
		if err != nil {
			return err
		}

		for setting, value := range multipleConfig {
			configFile.WriteString(setting + " = " + value + "\n")
		}

		configFile.Close()

		log.Info(emoji.Sprintf(":raised_hands: Azure Multiple cluster environment " + fullEnvironmentPath + " has been successfully created!"))
		log.Info(emoji.Sprintf(":white_check_mark: To proceed, run 'bedrock simulate " + environmentPath + "'"))

		return nil
	}
	return err
}
