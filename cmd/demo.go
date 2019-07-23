package cmd

import (
	"math/rand"
	"os/exec"
	"time"

	"github.com/docker/docker/pkg/namesgenerator"
	"github.com/kyokomi/emoji"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Create a cluster environment (azure simple, multi-cluster, keyvault, etc.)
func Demo() (err error) {
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
	randomName := namesgenerator.GetRandomName(0)

	// Check if Bedrock Repo is already cloned
	log.Info(emoji.Sprintf(":open_file_folder: Checking for Bedrock"))
	if output, err := exec.Command("git", "clone", "https://github.com/microsoft/bedrock").CombinedOutput(); output != nil || err != nil {
		log.Info(emoji.Sprintf(":star: Bedrock Repo already cloned"))
	}

	// Create an Azure SP
	/* 	log.Info(emoji.Sprintf(":cop: Creating New Service Principal"))
	   	if output, err := exec.Command("az", "ad", "sp", "create-for-rbac", "--subscription", "<subscription-id>").CombinedOutput(); err != nil {
	   		log.Error(emoji.Sprintf(":no_entry_sign: %s: %s", err, output))
	   		return err
	   	} */

	// Copy Terraform Template
	log.Info(emoji.Sprintf(":flashlight: Creating New Environment"))
	if output, err := exec.Command("cp", "-r", "bedrock/cluster/environments/azure-simple", "bedrock/cluster/environments/bedrock_"+randomName+"_cluster").CombinedOutput(); err != nil {
		log.Error(emoji.Sprintf(":no_entry_sign: %s: %s", err, output))
		return err
	}

	// Terraform Initialization
	/* 	log.Info(emoji.Sprintf(":checkered_flag: Terraform Init"))
	   	cmd := exec.Command("terraform", "init")
	   	cmd.Dir = "bedrock/cluster/environments/bedrock_" + randomName + "_cluster"
	   	if output, err := cmd.CombinedOutput(); err != nil {
	   		log.Error(emoji.Sprintf(":no_entry_sign: %s: %s", err, output))
	   		return err
	   	}
	*/
	if err == nil {
		log.Info(emoji.Sprintf(":raised_hands: Cluster has been successfully created!"))
	}

	return err
}

var demoCmd = &cobra.Command{
	Use:   "demo",
	Short: "Demo an Azure Kubernetes Service (AKS) cluster using Terraform",
	Long:  `Demo an Azure Kubernetes Service (AKS) cluster using Terraform`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return Demo()
	},
}

func init() {
	rootCmd.AddCommand(demoCmd)
}
