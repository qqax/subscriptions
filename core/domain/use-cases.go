package domain

// calculateOverlapMonths вычисляет количество пересекающихся месяцев
func calculateOverlapMonths(subStart, subEnd, periodStart, periodEnd string) (int, error) {
	// Конвертируем даты в числовой формат (YYYYMM)
	start1, err := dateToNumber(subStart)
	if err != nil {
		return 0, NewValidationError("subscription_start", "invalid date: "+err.Error())
	}

	end1, err := dateToNumber(subEnd)
	if err != nil {
		return 0, NewValidationError("subscription_end", "invalid date: "+err.Error())
	}

	start2, err := dateToNumber(periodStart)
	if err != nil {
		return 0, NewValidationError("period_start", "invalid date: "+err.Error())
	}

	end2, err := dateToNumber(periodEnd)
	if err != nil {
		return 0, NewValidationError("period_end", "invalid date: "+err.Error())
	}

	// Проверяем есть ли пересечение
	if end1 < start2 || start1 > end2 {
		return 0, nil // Нет пересечения
	}

	// Находим начало и конец пересечения
	overlapStart := max(start1, start2)
	overlapEnd := min(end1, end2)

	// Вычисляем количество месяцев
	startYear := overlapStart / 100
	startMonth := overlapStart % 100
	endYear := overlapEnd / 100
	endMonth := overlapEnd % 100

	months := (endYear-startYear)*12 + (endMonth - startMonth) + 1
	return months, nil
}
