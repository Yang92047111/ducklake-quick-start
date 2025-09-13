package storage

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Yang92047111/ducklake-quick-start/internal/loader"
	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(connectionString string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	repo := &PostgresRepository{db: db}
	if err := repo.createTables(); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return repo, nil
}

func (p *PostgresRepository) createTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS exercises (
		id SERIAL PRIMARY KEY,
		name VARCHAR(255) NOT NULL,
		type VARCHAR(100) NOT NULL,
		duration INTEGER NOT NULL,
		calories INTEGER NOT NULL,
		date DATE NOT NULL,
		description TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := p.db.Exec(query)
	return err
}

func (p *PostgresRepository) Insert(exercise loader.Exercise) error {
	query := `
	INSERT INTO exercises (name, type, duration, calories, date, description)
	VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := p.db.Exec(query, exercise.Name, exercise.Type, exercise.Duration,
		exercise.Calories, exercise.Date, exercise.Description)
	return err
}

func (p *PostgresRepository) InsertBatch(exercises []loader.Exercise) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
	INSERT INTO exercises (name, type, duration, calories, date, description)
	VALUES ($1, $2, $3, $4, $5, $6)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, exercise := range exercises {
		_, err := stmt.Exec(exercise.Name, exercise.Type, exercise.Duration,
			exercise.Calories, exercise.Date, exercise.Description)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (p *PostgresRepository) GetByID(id int) (*loader.Exercise, error) {
	query := `SELECT id, name, type, duration, calories, date, description FROM exercises WHERE id = $1`

	var exercise loader.Exercise
	err := p.db.QueryRow(query, id).Scan(&exercise.ID, &exercise.Name, &exercise.Type,
		&exercise.Duration, &exercise.Calories, &exercise.Date, &exercise.Description)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &exercise, nil
}

func (p *PostgresRepository) GetByDateRange(start, end time.Time) ([]loader.Exercise, error) {
	query := `SELECT id, name, type, duration, calories, date, description 
			  FROM exercises WHERE date BETWEEN $1 AND $2 ORDER BY date`

	rows, err := p.db.Query(query, start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return p.scanExercises(rows)
}

func (p *PostgresRepository) GetByType(exerciseType string) ([]loader.Exercise, error) {
	query := `SELECT id, name, type, duration, calories, date, description 
			  FROM exercises WHERE type = $1 ORDER BY date DESC`

	rows, err := p.db.Query(query, exerciseType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return p.scanExercises(rows)
}

func (p *PostgresRepository) GetAll() ([]loader.Exercise, error) {
	query := `SELECT id, name, type, duration, calories, date, description 
			  FROM exercises ORDER BY date DESC`

	rows, err := p.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return p.scanExercises(rows)
}

func (p *PostgresRepository) Update(exercise loader.Exercise) error {
	query := `UPDATE exercises SET name = $2, type = $3, duration = $4, 
			  calories = $5, date = $6, description = $7 WHERE id = $1`

	_, err := p.db.Exec(query, exercise.ID, exercise.Name, exercise.Type,
		exercise.Duration, exercise.Calories, exercise.Date, exercise.Description)
	return err
}

func (p *PostgresRepository) Delete(id int) error {
	query := `DELETE FROM exercises WHERE id = $1`
	_, err := p.db.Exec(query, id)
	return err
}

func (p *PostgresRepository) Close() error {
	return p.db.Close()
}

func (p *PostgresRepository) scanExercises(rows *sql.Rows) ([]loader.Exercise, error) {
	var exercises []loader.Exercise

	for rows.Next() {
		var exercise loader.Exercise
		err := rows.Scan(&exercise.ID, &exercise.Name, &exercise.Type,
			&exercise.Duration, &exercise.Calories, &exercise.Date, &exercise.Description)
		if err != nil {
			return nil, err
		}
		exercises = append(exercises, exercise)
	}

	return exercises, rows.Err()
}
