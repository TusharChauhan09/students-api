package postgres

import (
	"database/sql"

	"github.com/TusharChauhan09/students-api/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Postgres struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Postgres, error) {
	db, err := sql.Open("pgx", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS students (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			age INT NOT NULL
		)
	`)
	if err != nil {
		return nil, err
	}

	return &Postgres{
		Db: db,
	}, nil
}

func (p *Postgres) CreateStudent(name string, email string, age int) (int64 , error){

	var lastID int64

	err := p.Db.QueryRow(
		"INSERT INTO students (name, email, age) VALUES ($1, $2, $3) RETURNING id",
		name,
		email,
		age,
	).Scan(&lastID)

	if err != nil {
		return 0, err
	}

	return lastID, nil
	
}