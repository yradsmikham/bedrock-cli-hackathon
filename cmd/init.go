package cmd

import "C"
import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/docker/docker/pkg/namesgenerator"
	"github.com/kyokomi/emoji"
	log "github.com/sirupsen/logrus"
	util "github.com/yradsmikham/bedrock-cli/util"
)

// SSHKey is the public key
var SSHKey string
var randomClusterName string
var subnet string
var resources []string

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
func Init(environment string, clusterName string) (cluster string, resourceList []string, err error) {
	//resources := []string{}
	requiredSystemTools := []string{"git", "helm", "sh", "curl", "terraform", "az"}
	for _, tool := range requiredSystemTools {
		path, err := exec.LookPath(tool)
		if err != nil {
			return "", nil, err
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

	// Set Environment Variables
	if error := VerifyEnvVariables(clusterName, environment); error != nil {
		return "", nil, error
	}

	// Check if resource group exists, if it doesn't create it
	if resourceGroup == "" {
		if environment == COMMON {
			// Create the resource group
			log.Info(emoji.Sprintf(":construction: Creating new resource group: %s", clusterName+"-kv-rg"))
			_, err := exec.Command("az", "group", "create", "--name", clusterName+"-kv-rg", "--location", region).CombinedOutput()
			if err != nil {
				log.Error(emoji.Sprintf(":no_entry_sign: There was an error with creating the resource group!"))
				panic(fmt.Errorf("Please try again"))
			}
			resources = append(resources, clusterName+"-kv-rg")
		} else if environment == MULTIPLE {
			if resourceGroupWest == "" || resourceGroupCentral == "" || resourceGroupEast == "" {
				// Create resource groups for every region
				log.Info(emoji.Sprintf(":construction: Creating new resource group: %s", clusterName+"-west-rg, "+clusterName+"-central-rg, "+clusterName+"-east-rg"))
				_, westRgCreationErr := exec.Command("az", "group", "create", "--name", clusterName+"-west-rg", "--location", regionWest).CombinedOutput()
				if westRgCreationErr != nil {
					log.Error(emoji.Sprintf(":no_entry_sign: There was an error with creating the resource group!"))
					panic(fmt.Errorf("Please try again"))
				}
				resourceGroupWest = clusterName + "-west-rg"
				_, centralRgCreationErr := exec.Command("az", "group", "create", "--name", clusterName+"-central-rg", "--location", regionCentral).CombinedOutput()
				if centralRgCreationErr != nil {
					log.Error(emoji.Sprintf(":no_entry_sign: There was an error with creating the resource group!"))
					panic(fmt.Errorf("Please try again"))
				}
				resourceGroupEast = clusterName + "-east-rg"
				_, eastRgCreationErr := exec.Command("az", "group", "create", "--name", clusterName+"-east-rg", "--location", regionEast).CombinedOutput()
				if eastRgCreationErr != nil {
					log.Error(emoji.Sprintf(":no_entry_sign: There was an error with creating the resource group!"))
					panic(fmt.Errorf("Please try again"))
				}
				resourceGroupCentral = clusterName + "-central-rg"

				if resourceGroupTm == "" {
					_, trafficManagerErr := exec.Command("az", "group", "create", "--name", clusterName+"-tm-rg", "--location", regionEast).CombinedOutput()
					if trafficManagerErr != nil {
						log.Error(emoji.Sprintf(":no_entry_sign: There was an error with creating the resource group!"))
						panic(fmt.Errorf("Please try again"))
					}
					resourceGroupTm = clusterName + "-tm-rg"
				}
			}
			resources = append(resources, resourceGroupWest, resourceGroupEast, resourceGroupCentral, resourceGroupTm)
		} else {
			// Create the resource group
			log.Info(emoji.Sprintf(":construction: Creating new resource group: %s", clusterName+"-rg"))
			_, err := exec.Command("az", "group", "create", "--name", clusterName+"-rg", "--location", region).CombinedOutput()
			if err != nil {
				log.Error(emoji.Sprintf(":no_entry_sign: There was an error with creating the resource group!"))
				panic(fmt.Errorf("Please try again"))
			}
			//resourceGroup = clusterName + "-rg"
			resources = append(resources, clusterName+"-rg")
		}
	} else {
		log.Info(emoji.Sprintf(":mag_right: Verifying Resource Group..."))
		output, _ := exec.Command("az", "group", "show", "--name", resourceGroup).CombinedOutput()
		if strings.Contains(string(output), "could not be found") {
			log.Error(emoji.Sprintf(":question: The resource group specified does not exist!"))
			panic(fmt.Errorf("Please specify an existing resource group, or do not use the '--resource-group' to auto-generate one"))
		}
	}

	// Copy Terraform Template
	environmentPath := "bedrock/cluster/environments/" + clusterName
	if error := os.MkdirAll(environmentPath, os.ModePerm); error != nil {
		return "", nil, error
	}

	log.Info(emoji.Sprintf(":flashlight: Creating New Environment %s", environmentPath))
	if output, err := exec.Command("cp", "-r", "bedrock/cluster/environments/"+environment, environmentPath).CombinedOutput(); err != nil {
		log.Error(emoji.Sprintf(":no_entry_sign: %s: %s", err, output))
		return "", nil, err
	}

	// Generate SSH keys
	fullEnvironmentPath := environmentPath + "/" + environment
	if environment != COMMON {
		SSHKey, _ = SSH(fullEnvironmentPath, "deploy-key")
	}

	// Create bedrock-config.tfvars
	if err := addConfigTemplate(environment, fullEnvironmentPath, environmentPath, clusterName, SSHKey); err != nil {
		return "", nil, err
	}
	return clusterName, resources, err
}

// VerifyEnvVariables function verifies that SP is set
func VerifyEnvVariables(clusterName string, envType string) (err error) {

	_, subscriptionExists := os.LookupEnv("ARM_SUBSCRIPTION_ID")
	if subscriptionExists {
		log.Info(emoji.Sprintf(":globe_with_meridians: A Subscription ID was found in the environment variables."))
		subscription = os.Getenv("ARM_SUBSCRIPTION_ID")
	} else {
		if subscription == "" {
			log.Error(emoji.Sprintf(":confounded: A Subscription environment variable was not found. Please specify the ARM_SUBSCRIPTION_ID environment variable, or use the --subscription argument when creating the environment."))
			panic(fmt.Errorf("A Subscription ID needs to be specified"))
		} else {
			os.Setenv("ARM_SUBSCRIPTION_ID", subscription)
		}
	}
	_, spExists := os.LookupEnv("ARM_CLIENT_ID")
	if spExists {
		log.Info(emoji.Sprintf(":globe_with_meridians: A Service Principal was found in the environment variables."))
		servicePrincipal = os.Getenv("ARM_CLIENT_ID")
	} else {
		if servicePrincipal == "" {
			log.Error(emoji.Sprintf(":confounded: A Service Principal environment variable was not found. Please specify the ARM_CLIENT_ID environment variable, or use the --sp argument when creating the environment."))
			panic(fmt.Errorf("A Service Principal needs to be specified"))
		} else {
			os.Setenv("ARM_CLIENT_ID", servicePrincipal)
		}
	}
	_, secretExists := os.LookupEnv("ARM_CLIENT_SECRET")
	if secretExists {
		log.Info(emoji.Sprintf(":globe_with_meridians: A Service Principal Secret was found in the environment variables."))
		secret = os.Getenv("ARM_CLIENT_SECRET")
	} else {
		if secret == "" {
			log.Error(emoji.Sprintf(":confounded: A Service Principal Secret environment variable was not found. Please specify the ARM_CLIENT_SECRET environment variable, or use the --secret argument when creating the environment."))
			panic(fmt.Errorf("A Service Principal Password needs to be specified"))
		} else {
			os.Setenv("ARM_CLIENT_SECRET", secret)
		}
	}
	_, tenantExists := os.LookupEnv("ARM_TENANT_ID")
	if tenantExists {
		log.Info(emoji.Sprintf(":globe_with_meridians: A Service Principal Tenant ID was found in the environment variables."))
		tenant = os.Getenv("ARM_TENANT_ID")
	} else {
		if tenant == "" {
			log.Error(emoji.Sprintf(":confounded: A Service Principal Tenant ID environment variable was not found. Please specify the ARM_TENANT_ID environment variable, or use the --tenant argument when creating the environment."))
			panic(fmt.Errorf("A Service Principal Tenant ID needs to be specified"))
		} else {
			os.Setenv("ARM_TENANT_ID", tenant)
		}
	}
	return err
}

// GetEnvVariables function retrieves values from environment variables or sets them
func GetEnvVariables(clusterName string, envType string) (err error) {
	revisedClusterName := strings.Replace(clusterName, "-", "", -1)

	if envType == COMMON || envType == KEYVAULT || envType == MULTIPLE {
		if storageAccount == "" {
			_, exists := os.LookupEnv("AZURE_STORAGE_ACCOUNT")

			if exists {
				storageAccount = os.Getenv("AZURE_STORAGE_ACCOUNT")
			} else {
				error := util.CreateStorageAccount(revisedClusterName, clusterName+"-storage-rg", "centralus")
				resources = append(resources, clusterName+"-storage-rg")
				if error != nil {
					return error
				}
				storageAccount = revisedClusterName
			}
		}
		if accessKey == "" {
			_, exists := os.LookupEnv("AZURE_STORAGE_KEY")

			if exists {
				accessKey = os.Getenv("AZURE_STORAGE_KEY")
			} else {
				key, error := util.GetAccessKeys(revisedClusterName, clusterName+"-storage-rg")
				if error != nil {
					return error
				}
				accessKey = key
			}
		}
		if containerName == "" {
			_, exists := os.LookupEnv("AZURE_CONTAINER")

			if exists {
				containerName = os.Getenv("AZURE_CONTAINER")
			} else {
				if error := util.CreateStorageContainer(clusterName+"-container", revisedClusterName, accessKey); error != nil {
					return error
				}
				containerName = clusterName + "-container"
			}
		}
		if keyvaultName == "" {
			keyvaultName = clusterName + "-kv"
		}
		if keyvaultRG == "" {
			keyvaultRG = clusterName + "-kv-rg"
		}
	}
	if vnet == "" {
		vnet = clusterName + "-vnet"
	}
	if subnet == "" {
		subnet = clusterName + "-subnet"
	}
	if dnsPrefix == "" {
		dnsPrefix = clusterName
	}
	return err
}

// ReadTfvarsFile function will parse .tfvars file
func ReadTfvarsFile(filename string) (config map[string]string, err error) {
	tfvarsConfig := make(map[string]string)

	if len(filename) == 0 {
		return nil, err
	}
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
		return nil, err
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
				tfvarsConfig[key] = value
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return tfvarsConfig, err
}

// CopyFile is a function that copies a file to another destination
func CopyFile(source string, dest string) (err error) {
	sourcefile, err := os.Open(source)
	if err != nil {
		return err
	}

	defer sourcefile.Close()

	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer destfile.Close()

	_, err = io.Copy(destfile, sourcefile)
	if err == nil {
		sourceinfo, err := os.Stat(source)
		if err != nil {
			if err = os.Chmod(dest, sourceinfo.Mode()); err != nil {
				return err
			}
		}

	}

	return
}

// CopyDir is a function that copies an entire directory to another directory
func CopyDir(source string, dest string) (err error) {

	// Get properties of source dir
	sourceinfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	// Create destination dir
	err = os.MkdirAll(dest, sourceinfo.Mode())
	if err != nil {
		return err
	}

	directory, _ := os.Open(source)
	objects, err := directory.Readdir(-1)

	for _, obj := range objects {

		sourcefilepointer := source + "/" + obj.Name()
		destinationfilepointer := dest + "/" + obj.Name()

		if obj.IsDir() {
			// Create sub-directories - recursively
			err = CopyDir(sourcefilepointer, destinationfilepointer)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			// Perform copy
			err = CopyFile(sourcefilepointer, destinationfilepointer)
			if err != nil {
				fmt.Println(err)
			}
		}

	}
	return
}

// Generate bedrock-config.tfvars (and bedrock-config.toml) and bedrock-backend-config.tfvars (if appropriate)
func generateTfvars(envPath string, envType string, clusterName string, sshKey string) (err error) {

	configMap := make(map[string]string)
	backendConfigMap := make(map[string]string)
	spConfigMap := make(map[string]string)

	backendTfvarsFile, _ := os.Create(envPath + "/bedrock-backend-config.tfvars")
	tfvarsFile, _ := os.Create(envPath + "/bedrock-config.tfvars")
	spTomlFile, _ := os.Create(envPath + "/bedrock-sp-config.toml")
	log.Info(emoji.Sprintf(":page_with_curl: Create Bedrock config file " + envPath + "/bedrock-config.tfvars"))

	// Supported environments
	if envType == SIMPLE {
		azureSimpleTemplate(configMap, clusterName, sshKey)
		servicePrincipalTemplate(spConfigMap)
	}
	if envType == COMMON {
		backendTemplate(backendConfigMap, clusterName, COMMON)
		azureCommonInfraTemplate(configMap, clusterName, sshKey)
		servicePrincipalTemplate(spConfigMap)
	}
	if envType == KEYVAULT {
		backendTemplate(backendConfigMap, clusterName, KEYVAULT)
		azureSingleKVTemplate(configMap, clusterName, sshKey)
		servicePrincipalTemplate(spConfigMap)
	}
	if envType == MULTIPLE {
		azureMultipleTemplate(configMap, clusterName, sshKey)
		servicePrincipalTemplate(spConfigMap)
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
	for setting, value := range spConfigMap {
		if _, err := spTomlFile.WriteString(setting + " = " + value + "\n"); err != nil {
			return err
		}
	}
	spTomlFile.Close()

	return err
}

func servicePrincipalTemplate(config map[string]string) {
	config["subscription"] = "\"" + subscription + "\""
	config["service_principal"] = "\"" + servicePrincipal + "\""
	config["secret"] = "\"" + secret + "\""
	config["tenant_id"] = "\"" + tenant + "\""
}

func backendTemplate(config map[string]string, clusterName string, env string) {
	accessKey = strings.TrimSuffix(accessKey, "\n")
	config["storage_account_name"] = "\"" + storageAccount + "\""
	config["access_key"] = "\"" + accessKey + "\""
	config["container_name"] = "\"" + containerName + "\""
	config["key"] = "\"" + "tfstate-" + env + "-" + clusterName + "\""
}

func azureSimpleTemplate(config map[string]string, clusterName string, sshKey string) {
	config["resource_group_name"] = "\"" + clusterName + "-rg\""
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
	config["keyvault_name"] = "\"" + keyvaultName + "\""
	config["service_principal_id"] = "\"" + servicePrincipal + "\""
	config["address_space"] = "\"" + addressSpace + "\""
	config["subnet_prefix"] = "\"" + subnetPrefix + "\""
	config["subnet_name"] = "\"" + subnet + "\""
	config["vnet_name"] = "\"" + vnet + "\""
}

func azureSingleKVTemplate(config map[string]string, clusterName string, sshKey string) {
	config["resource_group_name"] = "\"" + clusterName + "-rg\""
	config["cluster_name"] = "\"" + clusterName + "\""
	config["agent_vm_size"] = "\"" + vmSize + "\""
	config["service_principal_id"] = "\"" + servicePrincipal + "\""
	config["service_principal_secret"] = "\"" + secret + "\""
	config["ssh_public_key"] = "\"" + sshKey + "\""
	config["gitops_ssh_url"] = "\"" + gitopsSSHUrl + "\""
	config["gitops_ssh_key"] = "\"" + "deploy-key" + "\""
	config["keyvault_resource_group"] = "\"" + keyvaultRG + "\""
	config["keyvault_name"] = "\"" + keyvaultName + "\""
	config["subnet_name"] = "\"" + subnet + "\""
	config["vnet_name"] = "\"" + vnet + "\""
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
	config["traffic_manager_resource_group_name"] = "\"" + resourceGroupTm + "\""
	//config["traffic_manager_resource_group_location"] = "\"" + regionWest + "\""
	config["west_resource_group_name"] = "\"" + resourceGroupWest + "\""
	//config["west_resource_group_location"] = "\"" + "westus2" + "\""
	config["gitops_west_path"] = "\"" + gitopsPathWest + "\""
	config["east_resource_group_name"] = "\"" + resourceGroupEast + "\""
	//config["east_resource_group_location"] = "\"" + regionEast + "\""
	config["gitops_east_path"] = "\"" + gitopsPathEast + "\""
	config["central_resource_group_name"] = "\"" + resourceGroupCentral + "\""
	//config["central_resource_group_location"] = "\"" + regionCentral + "\""
	config["gitops_central_path"] = "\"" + gitopsPathCentral + "\""
	config["gitops_central_url_branch"] = "\"" + gitopsURLBranchCentral + "\""
	config["gitops_east_url_branch"] = "\"" + gitopsURLBranchEast + "\""
	config["gitops_west_url_branch"] = "\"" + gitopsURLBranchWest + "\""
}

// Adds a blank bedrock config template
func addConfigTemplate(environment string, fullEnvironmentPath string, environmentPath string, clusterName string, sshKey string) (err error) {
	sshKey = strings.TrimSuffix(sshKey, "\n")

	if environment == SIMPLE {

		if error := GetEnvVariables(clusterName, SIMPLE); error != nil {
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

		if error := GetEnvVariables(clusterName, COMMON); error != nil {
			return error
		}
		if error := generateTfvars(fullEnvironmentPath, COMMON, clusterName, sshKey); error != nil {
			return error
		}

		commonInfraPath = environmentPath

		log.Info(emoji.Sprintf(":raised_hands: Azure Common Infra environment " + fullEnvironmentPath + " has been successfully created!"))

		return err
	}

	if environment == KEYVAULT {

		// When common infra is a dependency but does not exist, create one
		if commonInfraPath == "" {
			log.Info(emoji.Sprintf(":two_men_holding_hands: Common Infra path is not set, creating one now..."))
			_, _, error := Init(COMMON, clusterName)

			if error != nil {
				return error
			}
		} else {
			log.Info(emoji.Sprintf(":two_men_holding_hands: Contents of Azure Common Infra are being copied..."))
			if error := CopyDir(commonInfraPath, environmentPath); error != nil {
				return error
			}

			if _, err := os.Stat(environmentPath + "/" + COMMON + "/" + ".terraform"); err == nil {
				chmodCmd := exec.Command("chmod", "-R", "777", ".terraform")
				chmodCmd.Dir = string(environmentPath) + "/" + COMMON
				if output, err := chmodCmd.CombinedOutput(); err != nil {
					log.Error(emoji.Sprintf(":no_entry_sign: %s: %s", err, output))
					return err
				}
			} else {
				log.Info(emoji.Sprintf(":two_men_holding_hands: Terraform Init has not occurred for Azure Common Infra"))
			}

			configOutput, error := ReadTfvarsFile(environmentPath + "/" + COMMON + "/" + "bedrock-config.tfvars")
			if error != nil {
				log.Error(emoji.Sprintf(":no_entry_sign: %s", err))
				return err
			}
			subnet = configOutput["subnet_name"][1 : len(configOutput["subnet_name"])-1]
			vnet = configOutput["vnet_name"][1 : len(configOutput["vnet_name"])-1]
			keyvaultName = configOutput["keyvault_name"][1 : len(configOutput["keyvault_name"])-1]
			keyvaultRG = configOutput["global_resource_group_name"][1 : len(configOutput["global_resource_group_name"])-1]
		}

		log.Info(emoji.Sprintf(":family: Common Infra path is set to %s", commonInfraPath))

		if error := GetEnvVariables(clusterName, KEYVAULT); error != nil {
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

		if commonInfraPath != "" {
			log.Info(emoji.Sprintf(":two_men_holding_hands: Contents of Azure Common Infra are being copied..."))
			if error := CopyDir(commonInfraPath, environmentPath); error != nil {
				return error
			}

			if _, err := os.Stat(environmentPath + "/" + COMMON + "/" + ".terraform"); err == nil {
				chmodCmd := exec.Command("chmod", "-R", "777", ".terraform")
				chmodCmd.Dir = string(environmentPath) + "/" + COMMON
				if output, err := chmodCmd.CombinedOutput(); err != nil {
					log.Error(emoji.Sprintf(":no_entry_sign: %s: %s", err, output))
					return err
				}
			} else {
				log.Info(emoji.Sprintf(":two_men_holding_hands: Terraform Init has not occurred for Azure Common Infra"))
			}

			configOutput, error := ReadTfvarsFile(environmentPath + "/" + COMMON + "/" + "bedrock-config.tfvars")
			if error != nil {
				log.Error(emoji.Sprintf(":no_entry_sign: %s", err))
				return err
			}
			subnet = configOutput["subnet_name"][1 : len(configOutput["subnet_name"])-1]
			vnet = configOutput["vnet_name"][1 : len(configOutput["vnet_name"])-1]
			keyvaultName = configOutput["keyvault_name"][1 : len(configOutput["keyvault_name"])-1]
			keyvaultRG = configOutput["global_resource_group_name"][1 : len(configOutput["global_resource_group_name"])-1]
		}

		// When keyvault is not specified and common infra does not exist, create one
		if keyvaultName == "" && keyvaultRG == "" && commonInfraPath == "" {
			log.Info(emoji.Sprintf(":two_men_holding_hands: Common Infra path is not set, creating new Azure Common Infra environment"))
			if _, _, error := Init(COMMON, clusterName); error != nil {
				return error
			}
		}

		log.Info(emoji.Sprintf(":family: Common Infra path is set to %s", commonInfraPath))

		if error := GetEnvVariables(clusterName, MULTIPLE); error != nil {
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
