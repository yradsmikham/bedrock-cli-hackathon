package cmd

import (
	"os/exec"
	// "io/ioutil"

	"github.com/kyokomi/emoji"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// Apply and deploy a bedrock environment (azure simple, multi-cluster, keyvault, etc.)
func Testcmd(name string) (err error) {

	// KUBECONFIG=./output/bedrock_kube_config:~/.kube/config kubectl config view --flatten > merged-config && mv merged-config ~/.kube/config
	// TODO: Check if a bedrock kubeconfig output exists, then add that file to the local kubeconfig.
	log.Info(emoji.Sprintf(":mailbox_with_mail: Found Kubeconfig output. Merging into local kubeconfig."))
	mergeConfigCmd := exec.Command("/bin/sh", "-c", "KUBECONFIG=./output/bedrock_kube_config:~/.kube/config kubectl config view --flatten > merged-config && mv merged-config ~/.kube/config")
	mergeConfigCmd.Dir = name
	if output, err := mergeConfigCmd.CombinedOutput(); err != nil {
		log.Error(emoji.Sprintf(":no_entry_sign: %s: %s", err, output))
		return err
	}
	output, err := mergeConfigCmd.CombinedOutput()

	log.Info(string(output))

	if err == nil {
		log.Info(emoji.Sprintf(":raised_hands: Completed Terraform environment deployment!"))
	}
	
	return err
}

var testcmd = &cobra.Command{
	Use:   "testcmd <environment-name>",
	Short: "Deploy the bedrock environment using Terraform",
	Long:  `Deploy the bedrock environment deployment using terraform init and apply and adds the cluster credentials to the local kubeconfig.`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		log.Info("hello!!!!!")
		var name = "unique-environment-name"

		if len(args) > 0 {
			name = args[0]
		} 
		return Testcmd(name)
	},
}

func init() {
	rootCmd.AddCommand(testcmd)
}
