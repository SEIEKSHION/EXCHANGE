package models

import (
	"errors"
	"time"
)

var (
	EnteredEmptyNameError      = errors.New("You should enter not an empty name")
	EnteredInvalidMeasureError = errors.New("You should enter a point number more or equal 0")
	MuscleNotFoundError        = errors.New("Muscle not found")
)

type Muscle struct {
	Name          string    `json:"name"`      // <-- Большая буква
	Measure       float64   `json:"measure"`   // <-- Большая буква
	DateOfMeasure time.Time `json:"timestamp"` // <-- Большая буква + правильный тег
}

func NewMuscle(name string, measure float64, dateofmeasure time.Time) (Muscle, error) {
	if name == "" {
		return Muscle{}, EnteredEmptyNameError
	}
	if measure < 0 {
		return Muscle{}, EnteredInvalidMeasureError
	}
	return Muscle{
		Name:          name, // <-- Обновил имена полей
		Measure:       measure,
		DateOfMeasure: dateofmeasure.UTC(),
	}, nil
}

func (m *Muscle) Rename(newname string) error {
	if newname == "" {
		return EnteredEmptyNameError
	}
	m.Name = newname // <--
	return nil
}

func (m *Muscle) UpdateMeasure(newmeasure float64) error {
	if newmeasure <= 0 {
		return EnteredInvalidMeasureError
	}
	m.Measure = newmeasure // <--
	return nil
}

func FindByName(muscles []Muscle, name string) (Muscle, error) {
	if name == "" {
		return Muscle{}, EnteredEmptyNameError
	}
	for _, m := range muscles {
		if m.Name == name { // <--
			return m, nil
		}
	}
	return Muscle{}, MuscleNotFoundError
}
