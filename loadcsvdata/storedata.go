package loadcsvdata

import (
	"log"
	"sales/db"
)

func StoreData(pRec InsStruct) error {
	lPrdIns := `INSERT INTO product_details (
  prd_id, 
  name, 
  unitprice, 
  createdby, 
  createddate
) VALUES`

	lCusIns := `INSERT INTO customer_details (
  cus_id, 
  name, 
  email, 
  address, 
  createdby, 
  createddate
) VALUES`

	lOrdDel := `DELETE FROM sales_details
  WHERE order_id in `

	lOrdIns := `INSERT INTO sales_details (
  order_id, 
  product_id, 
  customer_id, 
  category, 
  region, 
  date_of_sale, 
  quantity_sold, 
  discount, 
  shipping_cost, 
  payment_method,
  createdby,
  createddate
) VALUES`

	lTrnsDB, lErr := db.GDBCon.Begin()
	if lErr != nil {
		return lErr
	}

	if len(pRec.PrdIns) > 0 {
		_, lErr = lTrnsDB.Exec(lPrdIns + pRec.PrdIns)
		if lErr != nil {
			lTrnsDB.Rollback()
			log.Println("RSD01 -", lErr)
			return lErr
		}
	}
	if len(pRec.CusIns) > 0 {
		_, lErr = lTrnsDB.Exec(lCusIns + pRec.CusIns)
		if lErr != nil {
			lTrnsDB.Rollback()
			log.Println("RSD02 -", lErr)
			return lErr
		}
	}

	if len(pRec.OrdMod) > 0 {
		_, lErr = lTrnsDB.Exec(lOrdDel + "(" + pRec.OrdMod + ")")
		if lErr != nil {
			lTrnsDB.Rollback()
			log.Println("RSD03 -", lErr)
			return lErr
		}
	}

	if len(pRec.OrdIns) > 0 {
		_, lErr = lTrnsDB.Exec(lOrdIns + pRec.OrdIns)
		if lErr != nil {
			lTrnsDB.Rollback()
			log.Println("RSD04 -", lErr)
			return lErr
		}
	}

	lTrnsDB.Commit()

	return nil
}
