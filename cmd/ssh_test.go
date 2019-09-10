// +build !testing

package cmd

import (
	"fmt"
	"os"
	"testing"
)

func TestSSH(t *testing.T) {

	if _, error := SSH(".", "testSSHKey"); error != nil {
		return
	}

	if fileExists("testSSHKey") {
		fmt.Println("Private SSH Key created successfully.")
	} else {
		t.Error("Private SSH Key was not generated.")
	}
	if fileExists("testSSHKey.pub") {
		fmt.Println("Public SSH Key created successfully.")
	} else {
		t.Error("Public SSH Key was not generated.")
	}

	os.Remove("testSSHKey")
	os.Remove("testSSHKey.pub")
}

// fileExists checks if a file exists
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
