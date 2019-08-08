package cmd

import (
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/kyokomi/emoji"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// SSH function will generate a new SSH Key using `ssh-keygen`
func SSH(path string, name string) (key string, err error) {
	log.Info(emoji.Sprintf(":lock_with_ink_pen: Creating new SSH with name " + name))

	keyPath := path + "/" + name
	log.Info(emoji.Sprintf(":gift_heart: Current wd is %s", path))

	if _, err := os.Stat(keyPath); err == nil {
		log.Info(emoji.Sprintf(":palm_tree: Path to %s exists, removing...", keyPath))
		os.Remove(keyPath)
		os.Remove(keyPath + ".pub")
	}

	// Create SSH Keys
	log.Info(emoji.Sprintf(":closed_lock_with_key: Creating New SSH Keys"))
	if output, err := exec.Command("ssh-keygen", "-t", "rsa", "-b", "4096", "-f", keyPath, "-P", "''").CombinedOutput(); err != nil {
		log.Error(emoji.Sprintf(":no_entry_sign: %s: %s", err, output))
		return "", err
	}

	log.Info(emoji.Sprintf(":raised_hands: SSH key " + name + " has been created!"))
	log.Info(emoji.Sprintf(":pray: Add the following SSH key to 'Deploy Keys' in your Manifest repository"))
	file, err := ioutil.ReadFile(keyPath + ".pub")
	if err != nil {
		log.Error(emoji.Sprintf(":no_entry_sign: %s: %s", err, file))
		return "", err
	}
	log.Info(string(file))
	return string(file), nil

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

		currentPath, err := os.Getwd()

		if err != nil {
			panic(err)
		}

		_, err = SSH(currentPath, name)

		return err
	},
}

func init() {
	rootCmd.AddCommand(sshCmd)
}
