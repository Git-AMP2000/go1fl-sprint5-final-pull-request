// Package spentenergy рассчитывает дистанцию, среднюю скорость и расход энергии.
package spentenergy

import (
	"errors"
	"math"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе.
)

// Distance возвращает дистанцию в километрах по количеству шагов и росту в метрах.
// Для нулевых и отрицательных входных значений возвращается 0.
func Distance(steps int, height float64) float64 {
	if steps <= 0 || height <= 0 || math.IsNaN(height) || math.IsInf(height, 0) {
		return 0
	}

	return float64(steps) * height * stepLengthCoefficient / mInKm
}

// MeanSpeed возвращает среднюю скорость в км/ч.
// При неположительной продолжительности или отрицательном количестве шагов возвращает 0.
func MeanSpeed(steps int, height float64, duration time.Duration) float64 {
	if steps < 0 || duration <= 0 {
		return 0
	}

	return Distance(steps, height) / duration.Hours()
}

// RunningSpentCalories возвращает расход энергии при беге в ккал.
// Количество шагов, вес, рост и продолжительность должны быть положительными.
func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if err := validateParameters(steps, weight, height, duration); err != nil {
		return 0, err
	}

	calories := weight * MeanSpeed(steps, height, duration) * duration.Minutes() / minInH
	if math.IsNaN(calories) || math.IsInf(calories, 0) {
		return 0, errors.New("результат расчёта калорий не является конечным числом")
	}

	return calories, nil
}

// WalkingSpentCalories возвращает расход энергии при ходьбе в ккал.
// Используется формула бега с корректирующим коэффициентом walkingCaloriesCoefficient.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	calories, err := RunningSpentCalories(steps, weight, height, duration)
	if err != nil {
		return 0, err
	}

	return calories * walkingCaloriesCoefficient, nil
}

func validateParameters(steps int, weight, height float64, duration time.Duration) error {
	switch {
	case steps <= 0:
		return errors.New("количество шагов должно быть положительным")
	case math.IsNaN(weight) || math.IsInf(weight, 0) || weight <= 0:
		return errors.New("вес должен быть положительным конечным числом")
	case math.IsNaN(height) || math.IsInf(height, 0) || height <= 0:
		return errors.New("рост должен быть положительным конечным числом")
	case duration <= 0:
		return errors.New("продолжительность должна быть положительной")
	default:
		return nil
	}
}
