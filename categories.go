package main

import (
	"log"

	"github.com/davidsteinsland/ynab-go/ynab"
)

func LoadCategories(budgetMgr *BudgetManager, client *ynab.Client) error {
	for _, bID := range budgetMgr.IDs() {
		categoryGroups, err := client.CategoriesService.List(bID)
		if err != nil {
			return err
		}

		bConfig := budgetMgr.GetByID(bID)
		for _, categoryGroup := range categoryGroups {
			for _, category := range categoryGroup.Categories {
				for i, confCategory := range bConfig.Categories {
					if categoryGroup.Name == confCategory.Group && category.Name == confCategory.Name {
						log.Printf("  Category '%s/%s' has id '%s'\n", categoryGroup.Name, category.Name, category.Id)
						bConfig.Categories[i].ID = category.Id
						break
					}
				}
			}
		}
	}

	return nil
}
