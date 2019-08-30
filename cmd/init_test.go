package cmd

import (
	"fmt"
	"os/exec"
	"testing"

	"github.com/kyokomi/emoji"
	log "github.com/sirupsen/logrus"
)

func TestInit(t *testing.T) {
	resources := []string{}
	env := map[string]string{
		"azuresimple":   SIMPLE,
		"azuresinglekv": KEYVAULT,
		"azuremultiple": MULTIPLE,
	}

	for k, v := range env {
		fmt.Println("Testing Init function for environment", k)
		_, resourceList, error := Init(v, "test"+k)
		if error != nil {
			t.Error("There was an error running the Init function.")
			return
		}
		for _, rg := range resourceList {
			//fmt.Print(rg)
			resources = append(resources, string(rg))
		}
		//fmt.Print(resourceList)
		if fileExists("bedrock/cluster/environments/test" + k + "/" + v + "/bedrock-config.tfvars") {
			log.Info(emoji.Sprintf(":trophy: Configuration file successfully built for environment %s", k))
		} else {
			t.Error("There was an error with creating bedrock-config.tfvars")
		}
	}
	uniqueResources := unique(resources)
	fmt.Println("-------------------RESOURCES:-------------------")
	fmt.Println(uniqueResources)
	if len(uniqueResources) >= 7 {
		t.Error("There was an error with creating resource groups for environments.")
	}

	// Clean up test environment
	log.Info(emoji.Sprintf(":boom: Deleting resources..."))
	for _, j := range uniqueResources {
		cmd := exec.Command("az", "group", "delete", "--name", j, "--yes")
		if output, err := cmd.CombinedOutput(); err != nil {
			log.Error(emoji.Sprintf(":no_entry_sign: %s: %s", err, output))
			return
		}
	}
	//os.RemoveAll("bedrock")
	log.Info(emoji.Sprintf(":white_check_mark: Simulation Test Complete!"))
}

func unique(intSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
