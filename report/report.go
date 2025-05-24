package report

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sales/db"
	"strings"
)

func GetRevenueDetails(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Methods", "PUT, OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	log.Println("GetRevenueDetails(+)")
	if r.Method == "PUT" {
		var lReq ReqData
		var lDebug Debug

		var lInterface interface{}

		lDebug.Status = "S"
		lDebug.ErrMsg = ""

		lBody, lErr := ioutil.ReadAll(r.Body)
		if lErr != nil {
			lDebug.Status = "E"
			lDebug.ErrMsg = "RGRD01 : " + lErr.Error()
			lInterface = lDebug
			goto Marshal
		}
		lErr = json.Unmarshal(lBody, &lReq)
		if lErr != nil {
			lDebug.Status = "E"
			lDebug.ErrMsg = "RGRD02 : " + lErr.Error()
			lInterface = lDebug
			goto Marshal
		}
		if lReq.Type == "" {
			var lRevResp TtlRevStruct
			lRevResp.Revenue, lErr = genTtlRevenue(lReq)
			if lErr != nil {
				lDebug.Status = "E"
				lDebug.ErrMsg = "RGRD03 : " + lErr.Error()
				lInterface = lDebug
				goto Marshal
			}
			lRevResp.Debug = lDebug
			lInterface = lRevResp
		} else {
			var lRevResp FieldRevStruct
			lRevResp.Revenue, lErr = genFieldRevenue(lReq)
			if lErr != nil {
				lDebug.Status = "E"
				lDebug.ErrMsg = "Error while fetching the data - " + lErr.Error()
				lInterface = lDebug
				goto Marshal
			}
			lRevResp.Debug = lDebug
			lInterface = lRevResp
		}
	Marshal:
		// log.Println("lInterface-", lInterface)
		data, err := json.Marshal(lInterface)
		if err != nil {
			fmt.Fprintf(w, "Error taking data"+err.Error())
		} else {
			fmt.Fprintf(w, string(data))
		}
		log.Println("GetRevenueDetails(-)")
	}

}

func genTtlRevenue(lReq ReqData) (float64, error) {
	var revenue float64

	query := `
				SELECT IFNULL(ROUND(SUM(((pd.unitprice * sd.quantity_sold) * (1 - sd.discount)) + sd.shipping_cost), 2),0) AS revenue
				FROM sales_details sd
				JOIN product_details pd ON sd.product_id = pd.prd_id
				WHERE sd.date_of_sale BETWEEN ? AND ? ;
			`
	err := db.GDBCon.QueryRow(query, lReq.FromDate, lReq.ToDate).Scan(&revenue)
	if err != nil {
		// log.Fatal("Query failed:", err)
		return 0, err
	}

	// fmt.Printf("Total Revenue: %.2f\n", revenue)
	return revenue, nil

}

func genFieldRevenue(lReq ReqData) ([]Searchdata, error) {
	var lArr []Searchdata

	lField := checktype(lReq.Type)
	if lField == "" {
		return nil, errors.New("Invalid type in request - " + lReq.Type)
	}

	query := `
        SELECT 
            ` + lField + `,
            ROUND(SUM(
                ((pd.unitprice * sd.quantity_sold) * (1 - sd.discount)) + sd.shipping_cost
            ), 2) AS revenue
        FROM sales_details sd
        JOIN product_details pd ON sd.product_id = pd.prd_id
        WHERE sd.date_of_sale BETWEEN ? AND ?
        GROUP BY ` + lField + `;
    `
	rows, err := db.GDBCon.Query(query, lReq.FromDate, lReq.ToDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var lRec Searchdata

		err := rows.Scan(&lRec.Name, &lRec.Revenue)
		if err != nil {
			return nil, err
		}
		lArr = append(lArr, lRec)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// log.Println(lReq.Type, "lArr - ", lArr)
	return lArr, nil
}

func checktype(pType string) string {
	switch strings.ToLower(pType) {
	case "product":
		return "pd.name"
	case "category":
		return "sd.category"
	case "region":
		return "sd.region"
	default:
		return ""
	}
}
