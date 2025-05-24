package loadcsvdata

var GFlag bool

type ProductDet struct {
	Prd_Id    string
	Prd_name  string
	Unit_Prc  string
	Cus_Id    string
	Cus_name  string
	Cus_email string
	Cus_addr  string
	Ord_Id    string
	Category  string
	Region    string
	SaleDate  string
	Qty       string
	Discount  string
	ShpCost   string
	PayMethod string
}

type InsStruct struct {
	PrdIns string
	CusIns string
	OrdIns string
	OrdMod string
}
