package utils

import (
	"os/exec"
	"github.com/kyokomi/emoji"
	log "github.com/sirupsen/logrus"
)

// Performs terraform init in the given directory
func TerraformInit(directory string) (err error) {
// Terraform Initialization (terraform init -backend-config=./bedrock-backend.tfvars)
	log.Info(emoji.Sprintf(":package: Terraform Init"))

	// TODO: If there is a bedrock-cli-backend.tfvars, then use that. OR default to backend.tfvars, OR just assume we're not using a backend deployment?
	// cmd := exec.Command("terraform", "init", "-backend-config=./bedrock-backend.tfvars")

	initCmd := exec.Command("terraform", "init")
	initCmd.Dir = directory
	if output, err := initCmd.CombinedOutput(); err != nil {
		log.Error(emoji.Sprintf(":no_entry_sign: %s: %s", err, output))
		return err
	}
	return err
}