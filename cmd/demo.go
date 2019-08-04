package cmd

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os/exec"
	"strings"
	"time"

	"github.com/docker/docker/pkg/namesgenerator"
	"github.com/kyokomi/emoji"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Demo is a function that will automate all the steps to creating an Azure Simple Cluster
func Demo(servicePrincipal string, secret string) (err error) {
	// Make sure host system contains all utils needed by Fabrikate
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
		log.Info(emoji.Sprintf(":star: Bedrock repo already cloned"))
	}

	// Copy Terraform Template
	log.Info(emoji.Sprintf(":flashlight: Creating New Environment: bedrock-" + randomName + "-cluster"))
	if output, err := exec.Command("cp", "-r", "bedrock/cluster/environments/"+SIMPLE, "bedrock/cluster/environments/bedrock-"+randomName+"-cluster").CombinedOutput(); err != nil {
		log.Error(emoji.Sprintf(":no_entry_sign: %s: %s", err, output))
		return err
	}

	// Generate SSH Keys
	SSH("bedrock/cluster/environments/bedrock-"+randomName+"-cluster", "deploy_key")
	sshPublic, err := ioutil.ReadFile("bedrock/cluster/environments/bedrock-" + randomName + "-cluster/deploy_key.pub")
	if err != nil {
		fmt.Print(err)
	}
	sshKey := string(sshPublic)

	// Terraform Init
	log.Info(emoji.Sprintf(":checkered_flag: Terraform Init"))
	tfInitCmd := exec.Command("terraform", "init")
	tfInitCmd.Dir = "bedrock/cluster/environments/bedrock-" + randomName + "-cluster"
	if output, err := tfInitCmd.CombinedOutput(); err != nil {
		log.Error(emoji.Sprintf(":no_entry_sign: %s: %s", err, output))
		return err
	}

	// Terraform Apply
	log.Info(emoji.Sprintf(":rocket: Terraform Apply"))
	tfApplyCmd := exec.Command("terraform", "apply", "-auto-approve", "--var", "resource_group_name=bedrock-"+randomName+"-rg", "--var", "cluster_name=bedrock-"+randomName+"-cluster", "--var", "dns_prefix=bedrock-"+randomName, "--var", "service_principal_id="+servicePrincipal, "--var", "service_principal_secret="+secret, "--var", "ssh_public_key="+sshKey, "--var", "gitops_ssh_key=deploy_key", "--var", "vnet_name=bedrock-"+randomName+"-vnet")
	tfApplyCmd.Dir = "bedrock/cluster/environments/bedrock-" + randomName + "-cluster"
	if output, err := tfApplyCmd.CombinedOutput(); err != nil {
		log.Error(emoji.Sprintf(":no_entry_sign: %s: %s", err, output))
		return err
	}

	// Add to KUBECONFIG
	log.Info(emoji.Sprintf(":heavy_plus_sign: Add Cluster to KubeConfig"))
	kubeConfigCmd := exec.Command("/bin/sh", "-c", "KUBECONFIG=./output/bedrock_kube_config:~/.kube/config kubectl config view --flatten > merged-config && mv merged-config ~/.kube/config")
	kubeConfigCmd.Dir = "bedrock/cluster/environments/bedrock-" + randomName + "-cluster"
	if output, err := kubeConfigCmd.CombinedOutput(); err != nil {
		log.Error(emoji.Sprintf(":no_entry_sign: %s: %s", err, output))
		return err
	}

	if err == nil {
		log.Info(emoji.Sprintf(":raised_hands: Cluster has been successfully created!"))
	}

	return err
}

var clusterName string
var servicePrincipal string
var secret string

var demoCmd = &cobra.Command{
	Use:   "demo [--sp service-principal-app-id] [--secret service-principal-password]",
	Short: "Demo an Azure Kubernetes Service (AKS) cluster using Terraform",
	Long:  `Demo an Azure Kubernetes Service (AKS) cluster using Terraform`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return Demo(servicePrincipal, secret)
	},
}

func init() {
	demoCmd.Flags().StringVar(&servicePrincipal, "sp", "", "Service Principal App Id")
	demoCmd.Flags().StringVar(&secret, "secret", "", "Password for the Service Principal")
	demoCmd.MarkFlagRequired("sp")
	demoCmd.MarkFlagRequired("secret")
	rootCmd.AddCommand(demoCmd)
}
