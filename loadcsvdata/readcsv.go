package loadcsvdata

import (
	"encoding/csv"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
)

func LoadCSVData(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	log.Println("LoadCSVData(+)")
	if r.Method == "POST" {
		lErr := ReadConsStore()
		if lErr != nil {
			http.Error(w, "Error - "+lErr.Error(), http.StatusNotImplemented)
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Data stored successfully ..."))
		}
	} else {
		http.Error(w, "Invalid method. Use POST.", http.StatusMethodNotAllowed)
	}
	log.Println("LoadCSVData(-)")
}

func ReadConsStore() error {
	// To avoid triggering while program in execution
	if GFlag {
		return errors.New(" Data Loading is in running phase ! ")
	}
	// If true, the execution begins
	GFlag = true
	defer func() {
		// If false, the execution completed
		GFlag = false
	}()
	lRecords, lErr := Readfile("./sample.csv")
	if lErr != nil {
		return lErr
	}
	lRec, lErr := ConstructRec(lRecords)
	if lErr != nil {
		return lErr
	}
	lErr = StoreData(lRec)
	if lErr != nil {
		return lErr
	}
	return nil
}

func Readfile(path string) ([][]string, error) { // Open CSV file
	f, lErr := os.Open(path)
	if lErr != nil {
		log.Println("RRF01 :", lErr)
		return nil, lErr
	}
	defer f.Close()

	// Create CSV reader
	reader := csv.NewReader(f)

	// Read all records
	records, lErr := reader.ReadAll()
	if lErr != nil {
		log.Println("RRF02 :", lErr)
		return nil, lErr
	}
	return records, nil
}

func ConstructRec(pRecords [][]string) (InsStruct, error) {

	// Retrieve the datas from the table
	prdMap, CusMap, OrdMap2, lErr := retrievedata()
	if lErr != nil {
		return InsStruct{}, lErr
	}

	var lPrdInsArr, lCusInsArr, lOrdInsArr, lOrdModArr []string

	// loopint each row in a csv
	for idx, v := range pRecords {
		// Omit Headers
		if idx == 0 {
			// for idx2, v2 := range v {
			// 	log.Println(idx2, "header -", v2)
			// }
			continue
		} else {
			lRec := assignstruct(v)

			// checking duplicate product id and name
			if lVal, Exists := prdMap[lRec.Prd_Id]; Exists {
				if lVal != lRec.Prd_name {
					log.Println("Duplicate prd name for prd id - " + lVal)
					return InsStruct{}, errors.New("Duplicate prd name for prd id - " + lVal)
				}
			} else {
				// Construct New Products for bulk insert
				lPrdStr := `('` + lRec.Prd_Id + `','` + lRec.Prd_name + `',` + lRec.Unit_Prc + `,'AUTOBOT',Now())`
				lPrdInsArr = append(lPrdInsArr, lPrdStr)
				prdMap[lRec.Prd_Id] = lRec.Prd_name
			}

			// checking duplicate customer id and name
			if lVal, Exists := CusMap[lRec.Cus_Id]; Exists {
				if lVal != lRec.Cus_name {
					log.Println("Duplicate cus name for cus id")
					return InsStruct{}, errors.New("Duplicate cus name for cus id - " + lVal)
				}
			} else {
				// Construct New Customers for bulk insert
				lCusStr := `('` + lRec.Cus_Id + `','` + lRec.Cus_name + `','` + lRec.Cus_email + `','` + lRec.Cus_addr + `','AUTOBOT',Now())`
				lCusInsArr = append(lCusInsArr, lCusStr)
				CusMap[lRec.Cus_Id] = lRec.Cus_name
			}

			lOrdRec := `(` + lRec.Ord_Id + `,'` + lRec.Prd_Id + `','` + lRec.Cus_Id + `','` + lRec.Category + `','` + lRec.Region + `','` + lRec.SaleDate + `',` + lRec.Qty + `,` + lRec.Discount + `,` + lRec.ShpCost + `,'` + lRec.PayMethod + `','AUTOBOT',Now())`

			// checking duplicate orders
			if val, Exists := OrdMap2[lRec.Ord_Id]; Exists {
				// if val is +ve, order duplicate in file
				if val >= 0 {
					// So overriding with latest order to insert
					lOrdInsArr[val] = lOrdRec
				} else {
					// if val is -ve, order id already exist in table
					// Need to override the existing data. So appending existing order to delete
					lOrdModArr = append(lOrdModArr, `'`+lRec.Ord_Id+`'`)
					// And then appending current order to insert
					lOrdInsArr = append(lOrdInsArr, lOrdRec)
					// And storing the index in the map, to check duplicate order id in a file
					OrdMap2[lRec.Ord_Id] = len(lOrdInsArr) - 1
				}
			} else {
				// Appending current order to insert
				lOrdInsArr = append(lOrdInsArr, lOrdRec)
				// And storing the index in the map, to check duplicate order id in a file
				OrdMap2[lRec.Ord_Id] = len(lOrdInsArr) - 1
			}

		}
	}

	lRec := InsStruct{
		// Converting the Array to string, to perform db execution quickly in a single hit.
		PrdIns: strings.Join(lPrdInsArr, ","),
		CusIns: strings.Join(lCusInsArr, ","),
		OrdMod: strings.Join(lOrdModArr, ","),
		OrdIns: strings.Join(lOrdInsArr, ","),
	}

	return lRec, nil
}

func assignstruct(v []string) ProductDet {
	var lRec ProductDet

	lRec.Ord_Id = v[0]
	lRec.Prd_Id = v[1]
	lRec.Cus_Id = v[2]
	lRec.Prd_name = v[3]
	lRec.Category = v[4]
	lRec.Region = v[5]
	lRec.SaleDate = v[6]
	lRec.Qty = v[7]
	lRec.Unit_Prc = v[8]
	lRec.Discount = v[9]
	lRec.ShpCost = v[10]
	lRec.PayMethod = v[11]
	lRec.Cus_name = v[12]
	lRec.Cus_email = v[13]
	lRec.Cus_addr = v[14]

	lRec.Prd_name = strings.ReplaceAll(lRec.Prd_name, "'", "''")
	lRec.Cus_name = strings.ReplaceAll(lRec.Cus_name, "'", "''")
	lRec.PayMethod = strings.ReplaceAll(lRec.PayMethod, "'", "''")
	lRec.Category = strings.ReplaceAll(lRec.Category, "'", "''")
	lRec.Region = strings.ReplaceAll(lRec.Region, "'", "''")
	lRec.Cus_addr = strings.ReplaceAll(lRec.Cus_addr, "'", "''")
	return lRec
}
