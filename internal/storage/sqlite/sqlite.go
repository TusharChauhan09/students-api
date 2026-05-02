package sqlite

import (
	"database/sql"

	"github.com/TusharChauhan09/students-api/internal/config"
	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	Db *sql.DB
}

func New (cfg *config.Config) (*Sqlite, error){
	db , err := sql.Open("sqlite3",cfg.StoragePath)
	if err != nil{
		return nil, err
	}

	_,err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
	id INT PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	age INT,
	email TEXT
	)`)

	if err!=nil{
		return nil,err
	}

	return &Sqlite{Db: db},nil
}

func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {
	stmt, err := s.Db.Prepare(
		"INSERT INTO students (name, email, age) VALUES (?, ?, ?)",
	)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastID, nil
}