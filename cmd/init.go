package cmd

import (
	"errors"
	"time"
	"math/rand"
	"os"
	"os/exec"

	"github.com/docker/docker/pkg/namesgenerator"
	"github.com/kyokomi/emoji"
	"github.com/spf13/cobra"
	log "github.com/sirupsen/logrus"
)

// Initializes the configuration for the given environment
func Init(environment string) (err error) {

	if (environment != "azure-simple"){
		return errors.New("Environment " + environment + " not supported.")
	}

	requiredSystemTools := []string{"git", "helm", "sh", "curl", "terraform", "az"}
	for _, tool := range requiredSystemTools {
		path, err := exec.LookPath(tool)
		if err != nil {
			return err
		}
		log.Info(emoji.Sprintf(":mag: Using %s: %s", tool, path))
	}

	rand.Seed(time.Now().UnixNano())
	randomName := namesgenerator.GetRandomName(0)

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
	fullEnvironmentPath := environmentPath + "/" + environment
	SSH(fullEnvironmentPath, "my-key")

	// Save bedrock-config.tfvars
	err = addConfigTemplate(fullEnvironmentPath, environment, randomName)

	if err == nil {
		log.Info(emoji.Sprintf(":raised_hands: Cluster " + fullEnvironmentPath +" has been successfully created!"))
		return nil
	}

	return err
}

// Adds a blank bedrock config template
func addConfigTemplate(path string, environment string, clusterName string)(error){
   
	if (environment == "azure-simple") {
 
	 azureSimpleConfig := make(map[string]string)
 
	 azureSimpleConfig["resource_group_name"] = "\"" + clusterName + "_rg\""
	 azureSimpleConfig["resource_group_location"] = "\"\""
	 azureSimpleConfig["cluster_name"] = "\"" + clusterName +"\""
	 azureSimpleConfig["agent_vm_count"]  = "\"\""
	 azureSimpleConfig["dns_prefix"] = "\"\""
	 azureSimpleConfig["service_principal_id"] = servicePrincipalID
	 azureSimpleConfig["service_principal_secret"] = servicePrincipalSecret
	 azureSimpleConfig["ssh_public_key"]  = "\"\"" // To-do: read the ssh public key file
	 azureSimpleConfig["gitops_ssh_url"] = gitopsSSHURL
	 azureSimpleConfig["gitops_ssh_key"] = path + "/keys/repo_ssh_key" // To-do get correct name
	 azureSimpleConfig["vnet_name"] = "\"\""

	 f, err := os.Create(path + "/bedrock-config.tfvars")
	 log.Info(emoji.Sprintf(":raised_hands: Create Bedrock config file " + path + "/bedrock-config.tfvars"))
	 if err != nil {
		 return err
	 }

	 for setting, value := range azureSimpleConfig {
		f.WriteString(setting + " = " + value + "\n")
	 }

	 f.Close()
 
	 return nil
	}
 
	return errors.New("Environment " + environment + " not supported.")
 }
 

var environment string
var servicePrincipalID string
var servicePrincipalSecret string
var gitopsSSHURL string

var initCmd = &cobra.Command{
	Use:   "init <config>",
	Short: "Init an Azure Kubernetes Service (AKS) cluster configuration",
	Long:  `Init an Azure Kubernetes Service (AKS) cluster configuration`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		return Init(environment)
	},
}


func init() {
	initCmd.Flags().StringVar(&environment, "environment", "azure-simple", "The cluster environment to be generated")
	initCmd.Flags().StringVar(&servicePrincipalID, "service-principal-id", "", "The Azure service principal ID")
	initCmd.Flags().StringVar(&servicePrincipalSecret, "service-principal-secret", "", "The Azure service principal secret")
	initCmd.Flags().StringVar(&gitopsSSHURL, "gitops-ssh-url", "", "The url of the GitOps repository")
	initCmd.MarkFlagRequired("service-principal-id")
	initCmd.MarkFlagRequired("service-principal-secret")
	initCmd.MarkFlagRequired("gitops-ssh-url")

	rootCmd.AddCommand(initCmd)
}