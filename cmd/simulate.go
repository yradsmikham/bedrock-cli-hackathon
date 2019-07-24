package cmd

import (
	"os/exec"
	// "io/ioutil"

	"github.com/kyokomi/emoji"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Simulate or dry-run a bedrock environment creation (azure simple, multi-cluster, keyvault, etc.)
func Simulate(name string) (err error) {
	
	// TODO: For each subdirectory inside the named environment directory, run terraform init and plan?
	// Alternatively, just look for *common*, and run that directory first?

	// Terraform Initialization (terraform init -backend-config=./bedrock-backend.tfvars)
	log.Info(emoji.Sprintf(":package: Terraform Init"))

	// TODO: If there is a bedrock-cli-backend.tfvars, then use that. OR default to backend.tfvars, OR just assume we're not using a backend deployment?
	// cmd := exec.Command("terraform", "init", "-backend-config=./bedrock-backend.tfvars")

	initCmd := exec.Command("terraform", "init")
	initCmd.Dir = name
	if output, err := initCmd.CombinedOutput(); err != nil {
		log.Error(emoji.Sprintf(":no_entry_sign: %s: %s", err, output))
		return err
	}

	// Terraform Plan (terraform plan -var-file=./bedrock-terraform.tfvars)
	// TODO: Check that ./bedrock-terraform.tfvars exists. Throw error if it doesn't, (or default to terraform.tfvars?)
	log.Info(emoji.Sprintf(":hammer: Terraform Plan"))
	planCmd := exec.Command("terraform", "plan", "-var-file=./bedrock-terraform.tfvars")
	planCmd.Dir = name
	if output, err := planCmd.CombinedOutput(); err != nil {
		log.Error(emoji.Sprintf(":no_entry_sign: %s: %s", err, output))
		return err
	}

	if err == nil {
		log.Info(emoji.Sprintf(":raised_hands: Completed simulated dry-run of environment deployment!"))
	}

	return err
}

var simulateCmd = &cobra.Command{
	Use:   "simulate <environment-name>",
	Short: "Simulate the environment deployment using Terraform",
	Long:  `Simulate the environment deployment using terraform init and plan`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {

		var name = "unique-environment-name"
		
		if len(args) > 0 {
			name = args[0]
		} 
		return Simulate(name)
	},
}

func init() {
	rootCmd.AddCommand(simulateCmd)
}
