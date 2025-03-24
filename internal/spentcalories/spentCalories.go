package spentcalories

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep = 0.65 // средняя длина шага.
	mInKm   = 1000 // количество метров в километре.
	minInH  = 60   // количество минут в часе.
)

func parseTraining(data string) (int, string, time.Duration, error) {
	splitString := strings.Split(data, ",")

	if len(splitString) != 3 {
		return 0, "", time.Duration(0), errors.New("неверный формат для парсинга строки")
	}

	steps, err := strconv.Atoi(splitString[0])
	if err != nil {
		return 0, "", time.Duration(0), err
	}

	activity := splitString[1]

	timeDuration, err := time.ParseDuration(splitString[2])
	if err != nil {
		return 0, "", time.Duration(0), err
	}

	return steps, activity, timeDuration, nil
}

// distance возвращает дистанцию(в километрах), которую преодолел пользователь за время тренировки.
//
// Параметры:
//
// steps int — количество совершенных действий (число шагов при ходьбе и беге).
func distance(steps int) float64 {
	return float64(steps) * lenStep / mInKm
}

// meanSpeed возвращает значение средней скорости движения во время тренировки.
//
// Параметры:
//
// steps int — количество совершенных действий(число шагов при ходьбе и беге).
// duration time.Duration — длительность тренировки.
func meanSpeed(steps int, duration time.Duration) float64 {
	if duration <= time.Duration(0) {
		return 0
	}

	calculatedDistance := distance(steps)

	return calculatedDistance / duration.Hours()
}

// ShowTrainingInfo возвращает строку с информацией о тренировке.
//
// Параметры:
//
// data string - строка с данными.
// weight, height float64 — вес и рост пользователя.
func TrainingInfo(data string, weight, height float64) string {
	steps, activity, timeDuration, err := parseTraining(data)
	if err != nil {
		return err.Error()
	}

	calories := 0.0

	switch activity {
	case "Ходьба":
		calories = WalkingSpentCalories(steps, weight, height, timeDuration)
		break
	case "Бег":
		calories = RunningSpentCalories(steps, weight, timeDuration)
		break
	default:
		return "неизвестный тип тренировки"
	}

	dist := distance(steps)
	speed := meanSpeed(steps, timeDuration)

	return fmt.Sprintf(`
		Тип тренировки: %s
		Длительность: %.2f ч.
		Дистанция: %.2f км.
		Скорость: %.2f км/ч
		Сожгли калорий: %.2f
	`, activity, timeDuration.Hours(), dist, speed, calories)
}

// Константы для расчета калорий, расходуемых при беге.
const (
	runningCaloriesMeanSpeedMultiplier = 18.0 // множитель средней скорости.
	runningCaloriesMeanSpeedShift      = 20.0 // среднее количество сжигаемых калорий при беге.
)

// RunningSpentCalories возвращает количество потраченных колорий при беге.
//
// Параметры:
//
// steps int - количество шагов.
// weight float64 — вес пользователя.
// duration time.Duration — длительность тренировки.
func RunningSpentCalories(steps int, weight float64, duration time.Duration) float64 {
	meanS := meanSpeed(steps, duration)
	return ((runningCaloriesMeanSpeedMultiplier * meanS) - runningCaloriesMeanSpeedShift) * weight
}

// Константы для расчета калорий, расходуемых при ходьбе.
const (
	walkingCaloriesWeightMultiplier = 0.035 // множитель массы тела.
	walkingSpeedHeightMultiplier    = 0.029 // множитель роста.
)

// WalkingSpentCalories возвращает количество потраченных калорий при ходьбе.
//
// Параметры:
//
// steps int - количество шагов.
// duration time.Duration — длительность тренировки.
// weight float64 — вес пользователя.
// height float64 — рост пользователя.
func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) float64 {
	meanS := meanSpeed(steps, duration)
	return ((walkingCaloriesWeightMultiplier * weight) + (meanS*meanS/height)*walkingSpeedHeightMultiplier) * duration.Hours() * minInH
}
