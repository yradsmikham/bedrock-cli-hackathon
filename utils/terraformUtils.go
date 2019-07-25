package utils

import (
	"bufio"
	"os"
	"os/exec"
	"github.com/kyokomi/emoji"
	log "github.com/sirupsen/logrus"
)

func runCommandWithOutput(cmd *exec.Cmd) (err error) {

	// create a pipe for the output of the script
	cmdReader, err := cmd.StdoutPipe()
    if err != nil {
        log.Error(os.Stderr, "Error creating StdoutPipe for cmd", err)
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
        log.Error(os.Stderr, "Error starting Cmd", err)
        return
    }

    err = cmd.Wait()
    if err != nil {
        log.Error(os.Stderr, "Error waiting for Cmd", err)
        return
	}

	return err
}

func TerraformInitWithOutput(directory string) (err error) {

    cmd := exec.Command("terraform", "init")
	cmd.Dir = directory

	return runCommandWithOutput(cmd)

    // // create a pipe for the output of the script
	// cmdReader, err := cmd.StdoutPipe()
	

    // scanner := bufio.NewScanner(cmdReader)
    // go func() {
    //     for scanner.Scan() {
    //         log.Info(scanner.Text())
    //     }
    // }()

    // err = cmd.Start()
    // if err != nil {
    //     log.Error(os.Stderr, "Error starting Cmd", err)
    //     return
    // }

    // err = cmd.Wait()
    // if err != nil {
    //     log.Error(os.Stderr, "Error waiting for Cmd", err)
    //     return
    // }

	// log.Info(emoji.Sprintf(":thumbsup: Terraform Init Complete!"))
	// return err
}

// Performs terraform init in the given directory
func TerraformInit(directory string) (err error) {
	// Terraform Initialization (terraform init -backend-config=./bedrock-backend.tfvars)
	log.Info(emoji.Sprintf(":package: Terraform Init Starting."))

	// TODO: If there is a bedrock-cli-backend.tfvars, then use that. OR default to backend.tfvars, OR just assume we're not using a backend deployment?
	// cmd := exec.Command("terraform", "init", "-backend-config=./bedrock-backend.tfvars")

	initCmd := exec.Command("terraform", "init")
	initCmd.Dir = directory
	if output, err := initCmd.CombinedOutput(); err != nil {
		log.Error(emoji.Sprintf(":no_entry_sign: %s: %s", err, output))
		return err
	}

	log.Info(emoji.Sprintf(":thumbsup: Terraform Init Complete!"))
	return err
}

func TerraformPlan(directory string) (err error) {
	log.Info(emoji.Sprintf(":hammer: Terraform Plan Starting."))

	// Terraform Plan (terraform plan -var-file=./bedrock-terraform.tfvars)
	// TODO: Check that ./bedrock-terraform.tfvars exists. Throw error if it doesn't, (or default to terraform.tfvars?)
	planCmd := exec.Command("terraform", "plan", "-var-file=./bedrock-terraform.tfvars")
	planCmd.Dir = directory
	if output, err := planCmd.CombinedOutput(); err != nil {
		log.Error(emoji.Sprintf(":no_entry_sign: %s: %s", err, output))
		return err
	}

	log.Info(emoji.Sprintf(":thumbsup: Terraform Plan Complete!"))
	return err
}

func TerraformApply(directory string) (err error) {
	log.Info(emoji.Sprintf(":hammer: Terraform Apply Starting."))

	log.Info(emoji.Sprintf(":bangbang: WARNING command is attempting to deploy real resources. :bangbang:"))

	applyCmd := exec.Command("terraform", "apply", "-var-file=./bedrock-terraform.tfvars", "-auto-approve")
	applyCmd.Dir = directory
	if output, err := applyCmd.CombinedOutput(); err != nil {
		log.Error(emoji.Sprintf(":no_entry_sign: %s: %s", err, output))
		return err
	}

	log.Info(emoji.Sprintf(":thumbsup: Terraform Apply Complete!"))
	return err
}