package report

type Searchdata struct {
	Name    string  `json:"name"`
	Revenue float64 `json:"revn"`
}

type ReqData struct {
	Type     string `json:"type"`
	FromDate string `json:"fdt"`
	ToDate   string `json:"tdt"`
}

type TtlRevStruct struct {
	Revenue float64 `json:"revn"`
	Debug
}

type FieldRevStruct struct {
	Revenue []Searchdata `json:"revndata"`
	Debug
}

type Debug struct {
	Status string `json:"sts"`
	ErrMsg string `json:"emsg"`
}
