package main

import (
	"log"

	"github.com/davidsteinsland/ynab-go/ynab"
)

type Category struct {
	ID    string
	Group string
	Name  string
}

type Budget struct {
	Name       string
	Account    string
	AccountID  string
	Categories []Category
	FlagColor  string `mapstrucutre:"flag-color"`
}

type BudgetManager struct {
	budgetsByName map[string]Budget
	budgetsByID   map[string]Budget
}

func NewBudgetManager(configBudgets []Budget, apiBudgets []ynab.BudgetSummary) BudgetManager {
	budgetsByName := make(map[string]Budget)
	budgetsByID := make(map[string]Budget)

	for _, budget := range configBudgets {
		budgetsByName[budget.Name] = budget
	}

	for _, budget := range apiBudgets {
		_, ok := budgetsByName[budget.Name]
		if ok {
			log.Printf("  Budget '%s' has id '%s'\n", budget.Name, budget.Id)
			budgetsByID[budget.Id] = budgetsByName[budget.Name]
		}
	}

	return BudgetManager{budgetsByName, budgetsByID}
}

func (b BudgetManager) IDs() []string {
	ids := make([]string, 0, len(b.budgetsByID))
	for id := range b.budgetsByID {
		ids = append(ids, id)
	}
	return ids
}

func (b BudgetManager) GetByID(id string) Budget {
	return b.budgetsByID[id]
}

func (b BudgetManager) GetByName(name string) Budget {
	return b.budgetsByName[name]
}

func (b BudgetManager) SetAccountID(budgetID, accountID string) {
	oldBudget, ok := b.budgetsByID[budgetID]
	if ok {
		b.budgetsByID[budgetID] = Budget{
			oldBudget.Name,
			oldBudget.Account,
			accountID,
			oldBudget.Categories,
			oldBudget.FlagColor,
		}
	}
}
