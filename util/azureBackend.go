package util

import (
	"os"
	"os/exec"

	"github.com/kyokomi/emoji"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// CreateStorageAccount function will create an Azure Storage Account if not provided
func CreateStorageAccount(storageAccount string, resourceGroup string, region string) (err error) {
	log.Info(emoji.Sprintf(":computer: Creating a Storage Account"))

	// Create Resource Group for Storage Account resources

	rgCmd := exec.Command("az", "group", "create", "--location", "centralus", "--name", resourceGroup)
	if output, err := rgCmd.CombinedOutput(); err != nil {
		log.Error(emoji.Sprintf(":no_entry_sign: %s: %s", err, output))
		return err
	}

	storageCmd := exec.Command("az", "storage", "account", "create", "--name", storageAccount, "--resource-group", resourceGroup, "--location", region, "--sku", "Standard_LRS", "--encryption", "blob")
	if output, err := storageCmd.CombinedOutput(); err != nil {
		log.Error(emoji.Sprintf(":no_entry_sign: %s: %s", err, output))
		return err
	}

	// Export as environment variable
	os.Setenv("AZURE_STORAGE_ACCOUNT", viper.GetString(storageAccount))

	log.Info(emoji.Sprintf(":raised_hands: Storage Account created!"))
	return err

}

// CreateStorageContainer function will create a blob storage in the Azure Storage Account
func CreateStorageContainer(storageContainer string, storageAccount string, accessKey string) (err error) {
	log.Info(emoji.Sprintf(":package: Creating a Storage Container"))

	cmd := exec.Command("az", "storage", "container", "create", "--name", storageContainer, "--account-name", storageAccount, "--account-key", accessKey)

	if output, err := cmd.CombinedOutput(); err != nil {
		log.Error(emoji.Sprintf(":no_entry_sign: %s: %s", err, output))
		return err
	}

	// Export as environment variable
	os.Setenv("AZURE_CONTAINER", viper.GetString(storageContainer))

	log.Info(emoji.Sprintf(":raised_hands: Storage Container created!"))
	return err
}

// GetAccessKeys function will retrieve the Azure Storage Account Access Keys
func GetAccessKeys(storageAccount string, resourceGroup string) (key string, err error) {
	log.Info(emoji.Sprintf(":key: Retreiving Storage Account Access Keys"))

	output, err := exec.Command("/bin/sh", "-c", "echo $(az storage account keys list --account-name "+storageAccount+" --resource-group "+resourceGroup+" --output table | awk 'FNR == 3 {print $3}')").Output()
	//output, err := exec.Command("az", "storage", "account", "keys", "list", "--account-name", storageAccount, "--resource-group", resourceGroup, "|", "awk", "'FNR == 3 {print $3}'")
	if err != nil {
		log.Error(emoji.Sprintf(":no_entry_sign: %s: %s", err, output))
		return "", err
	}

	// Export as environment variable
	os.Setenv("AZURE_STORAGE_KEY", viper.GetString(string(output)))

	return string(output), err
}
