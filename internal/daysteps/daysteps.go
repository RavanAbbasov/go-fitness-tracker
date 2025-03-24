package daysteps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/RavanAbbasov/go-fitness-tracker/internal/spentcalories"
)

var (
	StepLength = 0.65 // длина шага в метрах
)

/*
*
Принимает строку формата 678,0h50m
*/
func parsePackage(data string) (int, time.Duration, error) {
	splitStr := strings.Split(data, ",")

	if len(splitStr) != 2 {
		return 0, 0, errors.New("неверный формат для парсинга строки")
	}

	steps, err := strconv.Atoi(splitStr[0])
	if err != nil {
		return 0, 0, err
	}

	if steps <= 0 {
		return 0, 0, errors.New("количество шагов должно быть больше нуля")
	}

	timeDuration, err := time.ParseDuration(splitStr[1])
	if err != nil {
		return 0, 0, err
	}

	return steps, timeDuration, nil
}

// DayActionInfo обрабатывает входящий пакет, который передаётся в
// виде строки в параметре data. Параметр storage содержит пакеты за текущий день.
// Если время пакета относится к новым суткам, storage предварительно
// очищается.
// Если пакет валидный, он добавляется в слайс storage, который возвращает
// функция. Если пакет невалидный, storage возвращается без изменений.
func DayActionInfo(data string, weight, height float64) string {
	steps, timeDuration, err := parsePackage(data)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	distance := StepLength * float64(steps)
	distanceInKm := distance / 1000
	calories := spentcalories.WalkingSpentCalories(steps, weight, height, timeDuration)

	return fmt.Sprintf(`
		Количество шагов: %d.
		Дистанция составила %.2f км.
		Вы сожгли %.2f ккал.
	`, steps, distanceInKm, calories)
}
