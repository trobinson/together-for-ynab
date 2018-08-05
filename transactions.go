package main

import (
	"log"

	"github.com/davidsteinsland/ynab-go/ynab"
)

func GetTransactions(budgetMgr *BudgetManager, client *ynab.Client) (*map[string]ynab.HybridTransaction, error) {
	transactionsByID := make(map[string]ynab.HybridTransaction)
	for _, bID := range budgetMgr.IDs() {
		bConfig := budgetMgr.GetByID(bID)
		for _, bCategory := range bConfig.Categories {
			ts, err := client.TransactionsService.GetByCategory(bID, bCategory.ID)
			if err != nil {
				return nil, err
			}

			for _, t := range ts {
				color := bConfig.FlagColor
				t.FlagColor = &color
				transactionsByID[t.Id] = t
			}
		}
	}
	return &transactionsByID, nil
}

func CopyTransactions(srcMgr, dstMgr *BudgetManager, client *ynab.Client) ([]ynab.TransactionDetail, error) {
	var allDetail []ynab.TransactionDetail
	transactionsByID, err := GetTransactions(srcMgr, client)
	if err != nil {
		return allDetail, err
	}
	for _, bID := range dstMgr.IDs() {
		var bulkSubmit []ynab.SaveTransaction
		bConfig := dstMgr.GetByID(bID)
		for tID, t := range *transactionsByID {
			save := ynab.SaveTransaction{
				bConfig.AccountID,
				t.Date,
				t.Amount,
				"",
				t.PayeeName,
				"",
				"",
				"",
				false,
				*t.FlagColor,
				GetMD5Hash(tID),
			}

			bulkSubmit = append(bulkSubmit, save)
		}

		log.Printf("  Submitting %d transactions\n", len(bulkSubmit))
		detail, err := client.TransactionsService.CreateBulk(bID, bulkSubmit)
		allDetail = append(allDetail, detail...)
		if err != nil {
			return allDetail, err
		}
	}

	return allDetail, nil
}
