package util

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
		return err
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
		return err
	}

	if err = cmd.Wait(); err != nil {
		err = err.(*exec.ExitError)
		log.Error(err.Error())
		return err
	}

	return err
}

// TerraformInit will run `terraform init` in the given directory
func TerraformInit(directory string) (err error) {
	log.Info(emoji.Sprintf(":package: Terraform Init Starting."))

	tfInitCmd := exec.Command("terraform", "init")
	tfInitCmd.Dir = directory
	if err := runCommandWithOutput(tfInitCmd); err != nil {
		log.Error(emoji.Sprintf(":no_entry_sign: %s", err))
		return err
	}

	log.Info(emoji.Sprintf(":thumbsup: Terraform Init Complete!"))
	return err
}

// TerraformInitBackend will run `terraform init` with a backend in the given directory
func TerraformInitBackend(directory string) (err error) {
	log.Info(emoji.Sprintf(":package: Terraform Init Starting..."))

	tfInitBackendCmd := exec.Command("terraform", "init", "-backend-config=./bedrock-backend-config.tfvars")
	tfInitBackendCmd.Dir = directory
	if err := runCommandWithOutput(tfInitBackendCmd); err != nil {
		log.Error(emoji.Sprintf(":no_entry_sign: %s", err))
		return err
	}

	log.Info(emoji.Sprintf(":thumbsup: Terraform Init Complete!"))
	return err
}

// TerraformPlan will run `terraform plan` in given directory
func TerraformPlan(directory string) (err error) {
	log.Info(emoji.Sprintf(":hammer: Terraform Plan Starting..."))

	tfPlanCmd := exec.Command("terraform", "plan", "-var-file=bedrock-config.tfvars")
	tfPlanCmd.Dir = directory

	// Displays Terraform errors
	if output, err := tfPlanCmd.CombinedOutput(); err != nil {
		log.Error(emoji.Sprintf(":no_entry_sign: %s: %s", err, output))
		return err
	}

	//log.Info(emoji.Sprintf("DIRECTORY: %s", directory))
	/* 	if err := runCommandWithOutput(tfPlanCmd); err != nil {
	   		log.Error(emoji.Sprintf(":no_entry_sign: %s", err))
	   		return err
	   	}
	*/
	log.Info(emoji.Sprintf(":thumbsup: Terraform Plan Complete!"))
	return err
}

// TerraformApply will run `terraform apply` in given directory
func TerraformApply(directory string) (err error) {
	log.Info(emoji.Sprintf(":hammer: Terraform Apply Starting..."))
	log.Info(emoji.Sprintf(":bangbang: WARNING: COMMAND IS ATTEMPTING TO DEPLOY RESOURCES :bangbang:"))
	log.Info(emoji.Sprintf(":bangbang: IF YOU WOULD LIKE FOR THIS TO STOP, PRESS CRTL + C :bangbang:"))
	cmd := exec.Command("terraform", "apply", "-var-file=./bedrock-config.tfvars", "-auto-approve")
	cmd.Dir = directory

	runErr := runCommandWithOutput(cmd)

	log.Info(emoji.Sprintf(":thumbsup: Terraform Apply Complete!"))
	return runErr
}
