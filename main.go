package main

import (
	"fmt"
	"math"
	"time"
)

const (
	LenStep    = 0.65  // Средняя длина шага.
	MInKm      = 1000  // Количество метров в километре.
	MinsInHour = 60    // Количество минут в часе.
	KmHInMsec  = 0.278 // Коэффициент для преобразования км/ч в м/с.
	CmInM      = 100   // Количество сантиметров в метре.

	// Литералы сообщений
	messageTemplate = "Тип тренировки: %s\nДлительность: %v мин\nДистанция: %.2f км.\nСр. скорость: %.2f км/ч\nПотрачено ккал: %.2f\n"

	CaloriesMeanSpeedMultiplier = 18   // Множитель средней скорости бега
	CaloriesMeanSpeedShift      = 1.79 // Коэффициент изменения средней скорости

	CaloriesWeightMultiplier      = 0.035 // Коэффициент для веса
	CaloriesSpeedHeightMultiplier = 0.029 // Коэффициент для роста

	SwimmingLenStep                  = 1.38 // Длина одного гребка
	SwimmingCaloriesMeanSpeedShift   = 1.1  // Коэффициент изменения средней скорости
	SwimmingCaloriesWeightMultiplier = 2    // Множитель веса пользователя

)

// Training общая структура для всех тренировок
type Training struct {
	TrainingType string        // Тип тренировки
	Action       int           // Количество повторов(шаги, гребки при плавании)
	LenStep      float64       // Длина одного шага или гребка в м
	Duration     time.Duration // Продолжительность тренировки
	Weight       float64       // Вес пользователя в кг
}

// distance возвращает дистанцию, которую преодолел пользователь.
func (t Training) distance() float64 {
	return float64(t.Action) * t.LenStep / MInKm
}

// meanSpeed возвращает среднюю скорость бега или ходьбы.
func (t Training) meanSpeed() float64 {
	if t.Duration == 0 {
		return 0.0
	}
	return t.distance() / t.Duration.Hours()
}

// Calories возвращает количество потраченных килокалорий на тренировке.
func (t Training) Calories() float64 {
	return 0
}

// InfoMessage содержит информацию о проведенной тренировке.
type InfoMessage struct {
	TrainingType string        // Тип тренировки
	Duration     time.Duration // Длительность тренировки
	Distance     float64       // Расстояние, которое преодолел пользователь
	Speed        float64       // Средняя скорость, с которой двигался пользователь
	Calories     float64       // Количество потраченных килокалорий на тренировке
}

// TrainingInfo возвращает структуру InfoMessage, в которой хранится вся информация о проведенной тренировке.
func (t Training) TrainingInfo() InfoMessage {
	return InfoMessage{
		TrainingType: t.TrainingType,
		Duration:     t.Duration,
		Distance:     t.distance(),
		Speed:        t.meanSpeed(),
		Calories:     t.Calories(),
	}
}

// String возвращает строку с информацией о проведенной тренировке.
func (i InfoMessage) String() string {
	return fmt.Sprintf(messageTemplate,
		i.TrainingType,
		i.Duration.Minutes(),
		i.Distance,
		i.Speed,
		i.Calories,
	)
}

// CaloriesCalculator интерфейс для структур: Running, Walking и Swimming.
type CaloriesCalculator interface {
	Calories() float64
	TrainingInfo() InfoMessage
}

// Running структура, описывающая тренировку Бег.
type Running struct {
	Training
}

// Calories возвращает количество потраченных килокалория при беге.
func (r Running) Calories() float64 {
	averageSpeed := r.meanSpeed()
	if r.Duration == 0 {
		return 0.0
	}
	return (CaloriesMeanSpeedMultiplier*averageSpeed + CaloriesMeanSpeedShift) *
		r.Weight / MInKm * r.Duration.Hours() * MinsInHour
}

// TrainingInfo возвращает структуру InfoMessage с информацией о проведенной тренировке.
func (r Running) TrainingInfo() InfoMessage {
	return r.Training.TrainingInfo()
}

// Walking структура описывающая тренировку Ходьба
type Walking struct {
	Training
	Height float64 // Рост пользователя
}

// Calories возвращает количество потраченных килокалорий при ходьбе.
func (w Walking) Calories() float64 {
	// вставьте ваш код ниже
	averageSpeed := w.meanSpeed() * KmHInMsec
	if w.Height == 0 {
		return 0.0
	}
	return ((CaloriesWeightMultiplier*w.Weight + (math.Pow(averageSpeed, 2))/
		w.Height/CmInM) * CaloriesSpeedHeightMultiplier * w.Weight) *
		w.Duration.Hours() * MinsInHour
}

// TrainingInfo возвращает структуру InfoMessage с информацией о проведенной тренировке.
// Это переопределенный метод TrainingInfo() из Training.
func (w Walking) TrainingInfo() InfoMessage {
	return w.Training.TrainingInfo()
}

// Swimming структура, описывающая тренировку Плавание
type Swimming struct {
	Training
	LengthPool int // Длина бассейна
	CountPool  int // Количество пересечений бассейна
}

// distance возвращает дистанцию, которую проплыл пользователь.
func (s Swimming) distance() float64 {
	return float64(s.LengthPool*s.CountPool) / MInKm
}

// meanSpeed возвращает среднюю скорость при плавании.
func (s Swimming) meanSpeed() float64 {
	if s.Duration == 0 {
		return 0.0
	}
	return s.distance() / s.Duration.Hours()
}

// Calories возвращает количество калорий, потраченных при плавании.
func (s Swimming) Calories() float64 {
	averageSpeed := s.meanSpeed()
	return (averageSpeed + SwimmingCaloriesMeanSpeedShift) *
		SwimmingCaloriesWeightMultiplier * s.Weight * s.Duration.Hours()
}

// TrainingInfo returns info about swimming training.
func (s Swimming) TrainingInfo() InfoMessage {
	// вставьте ваш код ниже
	return InfoMessage{
		TrainingType: s.TrainingType,
		Duration:     s.Duration,
		Distance:     s.distance(),
		Speed:        s.meanSpeed(),
		Calories:     s.Calories(),
	}
}

// ReadData возвращает информацию о проведенной тренировке.
func ReadData(training CaloriesCalculator) string {
	info := training.TrainingInfo()
	info.Calories = training.Calories()
	return fmt.Sprint(info)
}

func main() {
	swimming := Swimming{
		Training: Training{
			TrainingType: "Плавание",
			Action:       2000,
			LenStep:      SwimmingLenStep,
			Duration:     90 * time.Minute,
			Weight:       85,
		},
		LengthPool: 50,
		CountPool:  5,
	}

	fmt.Println(ReadData(swimming))

	walking := Walking{
		Training: Training{
			TrainingType: "Ходьба",
			Action:       20000,
			LenStep:      LenStep,
			Duration:     3*time.Hour + 45*time.Minute,
			Weight:       85,
		},
		Height: 185,
	}

	fmt.Println(ReadData(walking))

	running := Running{
		Training: Training{
			TrainingType: "Бег",
			Action:       5000,
			LenStep:      LenStep,
			Duration:     30 * time.Minute,
			Weight:       85,
		},
	}

	fmt.Println(ReadData(running))

}
