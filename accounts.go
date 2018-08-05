package main

import (
	"log"

	"github.com/davidsteinsland/ynab-go/ynab"
)

func LoadAccounts(budgetMgr *BudgetManager, client *ynab.Client) error {
	for _, bID := range budgetMgr.IDs() {
		accounts, err := client.AccountsService.List(bID)
		if err != nil {
			return err
		}
		bConfig := budgetMgr.GetByID(bID)
		for _, account := range accounts {
			if bConfig.Account == account.Name {
				log.Printf("  Account '%s' has id '%s'\n", account.Name, account.Id)
				budgetMgr.SetAccountID(bID, account.Id)
				break
			}
		}
	}
	return nil
}
