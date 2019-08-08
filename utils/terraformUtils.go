package utils

import (
	"bufio"
	"os/exec"

	"github.com/kyokomi/emoji"
	log "github.com/sirupsen/logrus"
)

func runCommandWithOutput(cmd *exec.Cmd) (err error) {

	// create a pipe for the output of the script
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		log.Error("Error creating StdoutPipe for cmd", err)
		return
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			log.Info(scanner.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		log.Error("Error starting Cmd", err)
		return
	}

	err = cmd.Wait()
	if err != nil {
		log.Error("Error waiting for Cmd", err)
		return
	}

	return err
}

// TerraformInit will run `terraform init` in the given directory
func TerraformInit(directory string) (err error) {
	// Terraform Initialization (terraform init -backend-config=./bedrock-backend.tfvars)
	log.Info(emoji.Sprintf(":package: Terraform Init Starting."))

	// TODO: If there is a bedrock-cli-backend.tfvars, then use that. OR default to backend.tfvars, OR just assume we're not using a backend deployment?
	// cmd := exec.Command("terraform", "init", "-backend-config=./bedrock-backend.tfvars")

	cmd := exec.Command("terraform", "init")
	cmd.Dir = directory

	runErr := runCommandWithOutput(cmd)

	log.Info(emoji.Sprintf(":thumbsup: Terraform Init Complete!"))
	return runErr
}

// TerraformInitBackend will run `terraform init` with a backend in the given directory
func TerraformInitBackend(directory string) (err error) {
	// Terraform Initialization (terraform init -backend-config=./bedrock-backend.tfvars)
	log.Info(emoji.Sprintf(":package: Terraform Init Starting..."))

	// TODO: If there is a bedrock-cli-backend.tfvars, then use that. OR default to backend.tfvars, OR just assume we're not using a backend deployment?
	// cmd := exec.Command("terraform", "init", "-backend-config=./bedrock-backend.tfvars")

	cmd := exec.Command("terraform", "init", "-backend-config=./bedrock-backend-config.tfvars")
	cmd.Dir = directory

	runErr := runCommandWithOutput(cmd)

	log.Info(emoji.Sprintf(":thumbsup: Terraform Init Complete!"))
	return runErr
}

// TerraformPlan will run `terraform plan` in given directory
func TerraformPlan(directory string) (err error) {
	log.Info(emoji.Sprintf(":hammer: Terraform Plan Starting..."))

	// Terraform Plan (terraform plan -var-file=./bedrock-terraform.tfvars)
	// TODO: Check that ./bedrock-terraform.tfvars exists. Throw error if it doesn't, (or default to terraform.tfvars?)

	cmd := exec.Command("terraform", "plan", "-var-file=./bedrock-config.tfvars")
	cmd.Dir = directory

	runErr := runCommandWithOutput(cmd)

	log.Info(emoji.Sprintf(":thumbsup: Terraform Plan Complete!"))
	return runErr
}

// TerraformApply will run `terraform apply` in given directory
func TerraformApply(directory string) (err error) {
	log.Info(emoji.Sprintf(":hammer: Terraform Apply Starting..."))
	// TODO: Add confirmation input from user, suggest running Plan prior to Apply.
	log.Info(emoji.Sprintf(":bangbang: WARNING: COMMAND IS ATTEMPTING TO DEPLOY RESOURCES :bangbang:"))
	log.Info(emoji.Sprintf(":bangbang: IF YOU WOULD LIKE FOR THIS TO STOP, PRESS CRTL + C :bangbang:"))
	cmd := exec.Command("terraform", "apply", "-var-file=./bedrock-config.tfvars", "-auto-approve")
	cmd.Dir = directory

	runErr := runCommandWithOutput(cmd)

	log.Info(emoji.Sprintf(":thumbsup: Terraform Apply Complete!"))
	return runErr
}
