package util

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/kyokomi/emoji"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// createStorageAccount function will create an Azure Storage Account if not provided
func createStorageAccount(storageAccount string, resourceGroup string, region string) (err error) {
	log.Info(emoji.Sprintf(":computer: Creating a Storage Account"))

	cmd := exec.Command("az", "storage", "account", "create", "--name", storageAccount, "--resource-group", resourceGroup, "--location", region, "--sku", "Standard_LRS", "--encryption", "blob")

	if output, err := cmd.CombinedOutput(); err != nil {
		log.Error(emoji.Sprintf(":no_entry_sign: %s: %s", err, output))
		return err
	}

	// Export as environment variable
	os.Setenv("AZURE_STORAGE_ACCOUNT", viper.GetString(storageAccount))

	log.Info(emoji.Sprintf(":raised_hands: Storage Account created!"))
	return err

}

// createStorageContainer function will create a blob storage in the Azure Storage Account
func createStorageContainer(storageContainer string) (err error) {
	log.Info(emoji.Sprintf(":package: Creating a Storage Container"))

	cmd := exec.Command("az", "storage", "container", "--name", storageContainer)

	if output, err := cmd.CombinedOutput(); err != nil {
		log.Error(emoji.Sprintf(":no_entry_sign: %s: %s", err, output))
		return err
	}

	// Export as environment variable
	os.Setenv("AZURE_CONTAINER", viper.GetString(storageContainer))

	log.Info(emoji.Sprintf(":raised_hands: Storage Container created!"))
	return err
}

// getAccessKeys function will retrieve the Azure Storage Account Access Keys
func getAccessKeys(storageContainer string, resourceGroup string) (err error) {
	log.Info(emoji.Sprintf(":key: Retreiving Storage Account Access Keys"))

	cmd := exec.Command("/bin/sh", "-c", "az storage account keys list --account-name "+storageContainer+" --resource-group "+resourceGroup+" | awk 'FNR == 3 {print $3}'")

	output, err := cmd.Output()
	if err != nil {
		log.Error(emoji.Sprintf(":no_entry_sign: %s: %s", err, output))
		return err
	}
	fmt.Printf("The Storage Account Access Key: %s", output)

	// Export as environment variable
	os.Setenv("AZURE_STORAGE_KEY", viper.GetString(string(output)))

	return err
}
