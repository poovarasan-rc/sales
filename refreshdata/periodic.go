package refreshdata

import (
	"fmt"
	"sales/db"
	"sales/loadcsvdata"
	"time"
)

func PeriodicRefresh() {
	dbconfig := db.ReadTomlConfig("./toml/config.toml")

	RefreshDuration := fmt.Sprintf("%v", dbconfig.(map[string]interface{})["RefreshDuration"])

	lDuration, err := time.ParseDuration(RefreshDuration)
	if err != nil {
		lDuration = time.Duration(2 * time.Hour)
	}

	fmt.Println("Periodic refresh will run in", lDuration)

	time.AfterFunc(lDuration, func() {

		lErr := loadcsvdata.ReadConsStore()
		if lErr != nil {
			fmt.Println("Error while periodic refresh -", lErr)
		}

		PeriodicRefresh()

	})
}
