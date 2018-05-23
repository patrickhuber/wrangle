package main

import (
	"fmt"
	"os"

	"github.com/patrickhuber/cli-mgr/config"
	credhub "github.com/patrickhuber/cli-mgr/store/credhub"
	file "github.com/patrickhuber/cli-mgr/store/file"
)

func main() {
	configStoreManager := createConfigStoreManager()
	validateConfigStoreManager(configStoreManager)

	println("Success!")
}

func createConfigStoreManager() *config.ConfigStoreManager {
	manager := config.NewConfigStoreManager()
	manager.Register(&credhub.CredHubConfigStoreProvider{})
	manager.Register(&file.FileConfigStoreProvider{})
	return manager
}

func validateConfigStoreManager(manager *config.ConfigStoreManager) {
	if manager == nil {
		fmt.Printf("unable to create config store manager\n")
		os.Exit(1)
	}
}
