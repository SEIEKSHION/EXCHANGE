package repository

import (
	"github.com/SEIEKSHION/Exchanger/internal/models"
)

// MeasurementRepository содержит методы для работы с таблицей
type MeasurementRepository struct {
	client *DBClient
}

// NewMeasurementRepository создает новый репозиторий
func NewMeasurementRepository(client *DBClient) *MeasurementRepository {
	return &MeasurementRepository{client: client}
}

// Create добавляет новую запись
func (r *MeasurementRepository) Create(m *models.Measurement) error {
	_, err := r.client.DB.Exec(
		`INSERT INTO measurements ("user", muscle, measure, "date") 
		 VALUES ($1, $2, $3, $4)`,
		m.User, m.Muscle, m.Measure, m.Date,
	)
	return err
}

// GetByUser возвращает все записи пользователя
func (r *MeasurementRepository) GetByUser(user string) ([]*models.Measurement, error) {
	rows, err := r.client.DB.Query(
		`SELECT "user", muscle, measure, "date" 
		 FROM measurements WHERE "user" = $1`,
		user,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var measurements []*models.Measurement
	for rows.Next() {
		m := &models.Measurement{}
		if err := rows.Scan(&m.User, &m.Muscle, &m.Measure, &m.Date); err != nil {
			return nil, err
		}
		measurements = append(measurements, m)
	}
	return measurements, nil
}

// Другие методы (Update, Delete) реализуются аналогично
