package domain

import (
	"github.com/google/uuid"
)

func GenerateUUID() uuid.UUID {
	return uuid.New()
}

// dateToNumber конвертирует MM-YYYY в числовой формат YYYYMM
func dateToNumber(date string) (int, error) {
	year, month, err := ParseDate(date)
	if err != nil {
		return 0, err
	}
	return year*100 + month, nil
}
