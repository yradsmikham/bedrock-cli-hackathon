package cmd

import (
	"errors"
	"os/exec"

	"github.com/docker/docker/pkg/namesgenerator"
	"github.com/kyokomi/emoji"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Create a cluster environment (azure simple, multi-cluster, keyvault, etc.)
func Demo(environment string) (err error) {
	// Make sure host system contains all utils needed by Fabrikate
	requiredSystemTools := []string{"git", "helm", "sh", "curl", "terraform", "az"}
	for _, tool := range requiredSystemTools {
		path, err := exec.LookPath(tool)
		if err != nil {
			return err
		}
		log.Info(emoji.Sprintf(":mag: Using %s: %s", tool, path))
	}

	randomName := namesgenerator.GetRandomName(0)

	// Clone Microsoft Bedrock Repo
	log.Info(emoji.Sprintf(":crystal_ball: Cloning Bedrock Repo"))
	if output, err := exec.Command("git", "clone", "https://github.com/microsoft/bedrock").CombinedOutput(); err != nil {
		log.Error(emoji.Sprintf(":no_entry_sign: %s: %s", err, output))
		return err
	}

	// Create an Azure SP
	/* 	log.Info(emoji.Sprintf(":cop: Creating New Service Principal"))
	   	if output, err := exec.Command("az", "ad", "sp", "create-for-rbac", "--subscription", "<subscription-id>").CombinedOutput(); err != nil {
	   		log.Error(emoji.Sprintf(":no_entry_sign: %s: %s", err, output))
	   		return err
	   	} */

	// Copy Terraform Template
	log.Info(emoji.Sprintf(":flashlight: Creating New Environment"))
	if output, err := exec.Command("cp", "-r", "bedrock/cluster/environments/azure-simple", "bedrock/cluster/environments/bedrock-"+randomName).CombinedOutput(); err != nil {
		log.Error(emoji.Sprintf(":no_entry_sign: %s: %s", err, output))
		return err
	}

	if err == nil {
		log.Info(emoji.Sprintf(":raised_hands: Cluster has been successfully created!"))
	}

	return err
}

var environment string

var demoCmd = &cobra.Command{
	Use:   "demo <config>",
	Short: "Demo an Azure Kubernetes Service (AKS) cluster using Terraform",
	Long:  `Demo an Azure Kubernetes Service (AKS) cluster using Terraform`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		if !((args[0] == "simple") || (args[0] == "multi")) {
			return errors.New("the environment you specified is not of the following: simple, multi, keyvault")
		}
		return Demo(environment)
	},
}

func init() {
	rootCmd.AddCommand(demoCmd)
}
