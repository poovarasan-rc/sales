package loadcsvdata

import (
	"sales/db"
	"strings"
)

func retrievedata() (map[string]string, map[string]string, map[string]int, error) {

	lPrdMap, lErr := getStoredData(`select prd_id ,name from product_details`)
	if lErr != nil {
		return nil, nil, nil, lErr
	}
	lCUsMap, lErr := getStoredData(`select cus_id ,name from customer_details`)
	if lErr != nil {
		return nil, nil, nil, lErr
	}
	OrdMap2, lErr := getStoreOrder(`select order_id from sales_details`)
	if lErr != nil {
		return nil, nil, nil, lErr
	}
	return lPrdMap, lCUsMap, OrdMap2, lErr

}

func getStoredData(query string) (map[string]string, error) {

	lMap := make(map[string]string)

	rows, err := db.GDBCon.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Loop through results
	for rows.Next() {
		var lId, lName string

		err := rows.Scan(&lId, &lName)
		if err != nil {
			return nil, err
		}
		lMap[lId] = strings.ReplaceAll(lName, "'", "''")
	}
	// Check for row iteration error
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return lMap, nil
}

func getStoreOrder(query string) (map[string]int, error) {

	lMap := make(map[string]int)

	rows, err := db.GDBCon.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Loop through results
	for rows.Next() {
		var lId string

		err := rows.Scan(&lId)
		if err != nil {
			return nil, err
		}
		lMap[lId] = -1
	}
	// Check for row iteration error
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return lMap, nil
}
