package main

import (
	"log"

	"github.com/davidsteinsland/ynab-go/ynab"
	"github.com/spf13/viper"
)

func main() {
	/* Configure viper for configuration */
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/together4ynab")

	/* load in config */
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	/* create YNAB client */
	log.Println("Creating YNAB API Client")
	apiKey := viper.GetString("api-key")
	client := ynab.NewDefaultClient(apiKey)
	log.Printf("  Using API Key '%s'\n", apiKey)

	log.Println("Loading your YNAB budgets")
	apiBudgets, err := client.BudgetService.List()
	if err != nil {
		log.Fatal(err)
	}

	/* Get source budget names from config */
	var srcBudgets []Budget
	viper.UnmarshalKey("source-budgets", &srcBudgets)
	srcBudgetMan := NewBudgetManager(srcBudgets, apiBudgets)

	/* Get joint budget names from config */
	var jntBudgets []Budget
	viper.UnmarshalKey("joint-budgets", &jntBudgets)
	jntBudgetMan := NewBudgetManager(jntBudgets, apiBudgets)

	/* Get source budget Category IDs from API */
	log.Println("Retrieving IDs for configured categories")
	err = LoadCategories(&srcBudgetMan, client)
	if err != nil {
		log.Fatal(err)
	}

	/* Get joint budget AccountIDs from API */
	log.Println("Retrieving IDs for configured accounts")
	err = LoadAccounts(&jntBudgetMan, client)
	if err != nil {
		log.Fatal(err)
	}

	/* Process source transactions for adding to joint budget accounts */
	log.Println("Copying transactions")
	details, err := CopyTransactions(&srcBudgetMan, &jntBudgetMan, client)
	if err != nil {
		log.Fatal(err)
	}

	/* Print any details that come back from the API */
	for _, detail := range details {
		log.Println(detail)
	}
}
