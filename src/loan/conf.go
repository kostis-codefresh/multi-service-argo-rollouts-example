package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func (loanApp *LoanApplication) readCurrentConfiguration() {
	viper.SetDefault("role", "unknown")
	viper.SetDefault("backendHost", "localhost")
	viper.SetDefault("backendPort", "8080")

	viper.SetConfigName("labels")
	viper.SetConfigType("properties") //Java properties style

	//Development mode
	viper.AddConfigPath(".")

	//This is injected from the Kubernetes downward API that maps
	// all labels as a file in the pod
	// See https://kubernetes.io/docs/concepts/workloads/pods/downward-api/
	viper.AddConfigPath("/etc/podinfo/")

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	//Reload configuration when file changes
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		loanApp.reloadSettings()

	})

	loanApp.reloadSettings()

	viper.WatchConfig()

}

func (loanApp *LoanApplication) reloadSettings() {

	fmt.Printf("Role is set %t\n", viper.IsSet("role"))

	loanApp.CurrentRole = unQuoteIfNeeded(viper.GetString("role"))

	loanApp.BackendHost = unQuoteIfNeeded(viper.GetString("backendHost"))
	loanApp.BackendPort = unQuoteIfNeeded(viper.GetString("backendPort"))

	fmt.Printf("Role is %s\n", loanApp.CurrentRole)
	fmt.Printf("BackendHost is %s\n", loanApp.BackendHost)
	fmt.Printf("BackendPort is %s\n", loanApp.BackendPort)
}

func unQuoteIfNeeded(input string) string {
	result := ""
	if strings.HasPrefix(input, "\"") {
		result, _ = strconv.Unquote(input)
	}
	return result
}
