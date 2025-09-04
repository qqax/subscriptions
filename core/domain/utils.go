package domain

import "github.com/google/uuid"

// dateToNumber конвертирует MM-YYYY в числовой формат YYYYMM
func dateToNumber(date string) (int, error) {
	year, month, err := parseDate(date)
	if err != nil {
		return 0, err
	}
	return year*100 + month, nil
}

func generateUUID() uuid.UUID {
	return uuid.New()
}
