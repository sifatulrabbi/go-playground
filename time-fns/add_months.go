package timefns

import (
	"fmt"
	"time"
)

func AddingMonthsToATime() {
	startDate := time.Date(2025, time.June, 1, 0, 0, 0, 0, time.Local)
	fmt.Println("start date:", startDate.Format(time.RFC1123))
	endDate := startDate.AddDate(0, 5, 0)
	fmt.Println("start date:", endDate.Format(time.RFC1123))
}
