package main

import (
	"log"
	"net/http"
	"sales/db"
	"sales/loadcsvdata"
	"sales/refreshdata"
	"sales/report"
)

func main() {
	// Open DB Connection
	lErr := db.OpenConn()
	if lErr != nil {
		log.Println("Error while connecting db -", lErr)
	} else {
		// Close DB Connection
		defer db.CloseConn()
		log.Print("Server Started")

		// periodic referesh based on the configured values
		go refreshdata.PeriodicRefresh()

		// API to load the csv data into database by managing duplicates.
		http.HandleFunc("/loadcsv", loadcsvdata.LoadCSVData)
		// API to fetch the revenue details based on the request type.
		http.HandleFunc("/getrevenue", report.GetRevenueDetails)

		http.ListenAndServe(":8080", nil)
	}
}
