package cmd

import (
	"os/exec"
    "io/ioutil"
	"github.com/kyokomi/emoji"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

// Create a new SSH key
func SSH(name string) (err error) {
	log.Info(emoji.Sprintf(":lock_with_ink_pen: Creating new SSH with name " + name))
	currentPath, err := os.Getwd()
    if err != nil {
        panic(err)
    }
	keyPath := currentPath + "/" + name
	log.Info(emoji.Sprintf(":gift_heart: Current wd is %s", currentPath))

	if _, err := os.Stat(keyPath); err == nil {
		log.Info(emoji.Sprintf(":palm_tree: Path to %s exists, removing...", keyPath))
		os.Remove(keyPath)
		os.Remove(keyPath + ".pub")
	}
	
	// Make sure host system contains all utils needed by this module
	requiredSystemTools := []string{"git", "helm", "sh", "curl", "terraform", "az"}
	for _, tool := range requiredSystemTools {
		path, err := exec.LookPath(tool)
		if err != nil {
			return err
		}
		log.Info(emoji.Sprintf(":mag: Using %s: %s", tool, path))
	}

	// Create SSH Keys
	log.Info(emoji.Sprintf(":closed_lock_with_key: Creating New SSH Keys"))
	if output, err := exec.Command("ssh-keygen",  "-t", "rsa", "-b", "4096", "-f", name, "-P", "''").CombinedOutput(); err != nil {
		log.Error(emoji.Sprintf(":no_entry_sign: %s: %s", err, output))
		return err
	}

	if err == nil {
		log.Info(emoji.Sprintf(":raised_hands: SSH key " + name + " has been created!" ))
		log.Info(emoji.Sprintf(":pray: Add the following SSH key to 'Deploy Keys' in your Manifest repository"))
		file, err := ioutil.ReadFile(keyPath + ".pub") 
		if err != nil {
			log.Error(emoji.Sprintf(":no_entry_sign: %s: %s", err, file))
			return err
		}
		log.Info(string(file))
	}

	return err
}

var sshCmd = &cobra.Command{
	Use:   "ssh [ssh_key_name]",
	Short: "Create an SSH key",
	Long:  `Create an SSH key`,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		var name = "id_rsa"
		if len(args) > 0 {
			name = args[0]
		} 
		return SSH(name)
	},
}

func init() {
	rootCmd.AddCommand(sshCmd)
}
