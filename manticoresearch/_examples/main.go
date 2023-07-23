package main

import (
	"fmt"
	"mybakbibu/internal/manticoresearch"
)

func main() {
	// manticore configration
	manticoreConf := map[string]string{
		"url":       "http://127.0.0.1:9308",
		"readonly":  "false",
		"debugmode": "true",
	}

	mcdb := manticoresearch.NewManticoreClient(
		manticoresearch.RegisterMCHttpClient("", 0, false, manticoreConf["debugmode"] == "true"),
		manticoresearch.RegisterMCApiSettings(manticoreConf["url"], manticoreConf["readonly"] == "true"),
	)

	// Search Manager
	// searchManager := mcmanager.NewMcSearchManager(mcdb)

	// Start to Search

	// generate new query builder
	queryBuilder := manticoresearch.NewMCSearchQueryBuilder("models_mc")
	// queryBuilder.SetProfile(true)
	queryBuilder.SetOffset(2)
	queryBuilder.SetLimit(5)

	// Equal Query
	queryOptions := manticoresearch.NewMcQueryOptions()
	queryOptions.AddEqualsAll("maker_id", 3)
	queryBuilder.SetQuery(queryOptions)

	// Sort Query
	sortOptions := manticoresearch.NewMcSortOptions()
	sortOptions = sortOptions.SingleField("score", manticoresearch.MCSortOrderDESC)
	sortOptions = sortOptions.SingleField("updated_at", manticoresearch.MCSortOrderASC)
	queryBuilder.SetSort(sortOptions)

	resp, err := mcdb.Search(queryBuilder)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Resp: %#v\n", resp)
}
