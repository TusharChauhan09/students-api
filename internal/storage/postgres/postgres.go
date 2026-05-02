package postgres

import (
	"database/sql"
	"fmt"

	"github.com/TusharChauhan09/students-api/internal/config"
	"github.com/TusharChauhan09/students-api/internal/types"
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

// implements the interface : CreateStudent(name string, email string, age int) (int64 , error)
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

func (p *Postgres) GetStudentById(id int64) (types.Student,error) {
	var student types.Student
	err := p.Db.QueryRow(
		`SELECT id, name, email, age FROM students WHERE id = $1`,
		id,
	).Scan(
		&student.Id,
		&student.Name,
		&student.Email,
		&student.Age,
	)

	if err != nil {
		if err == sql.ErrNoRows{
			return types.Student{}, fmt.Errorf("no student found with id %d", id)
		}
		return types.Student{}, fmt.Errorf("query error: %w",err)
	}

	return student, nil

}

func (p *Postgres) GetStudents() ([]types.Student,error){
	
	rows, err := p.Db.Query(
		`SELECT id, name, email, age FROM students`,
	)
	if err != nil {
		return nil, fmt.Errorf("error while fetching students: %w", err)
	}
	defer rows.Close()
	
	var students []types.Student
	
	for rows.Next(){
		var student types.Student

		err := rows.Scan(
			&student.Id,
			&student.Name,
			&student.Email,
			&student.Age,
		)
		if err != nil {
			return nil , fmt.Errorf("scan error: %w",err)
		}
		students = append(students,student)
	}

	return students, nil
}

// prepare query : postgres

// func (p *Postgres) CreateStudent(name, email string, age int) (int64, error) {
// 	stmt, err := p.Db.Prepare(
// 		"INSERT INTO students(name,email,age) VALUES($1,$2,$3) RETURNING id",
// 	)
// 	if err != nil {
// 		return 0, err
// 	}
// 	defer stmt.Close()

// 	var id int64

// 	err = stmt.QueryRow(name, email, age).Scan(&id)
// 	if err != nil {
// 		return 0, err
// 	}

// 	return id, nil
// }

// func (p *Postgres) GetStudentByID(id int64) (*types.Student, error) {
// 	stmt, err := p.Db.Prepare(
// 		"SELECT id, name, email, age FROM students WHERE id = $1",
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer stmt.Close()

// 	var student types.Student

// 	err = stmt.QueryRow(id).Scan(
// 		&student.Id,
// 		&student.Name,
// 		&student.Email,
// 		&student.Age,
// 	)
// 	if err != nil {
// 		if errors.Is(err, sql.ErrNoRows) {
// 			return nil, fmt.Errorf("student %d not found", id)
// 		}
// 		return nil, err
// 	}

// 	return &student, nil
// }

// func (p *Postgres) GetStudents() ([]types.Student, error) {
// 	stmt, err := p.Db.Prepare(
// 		`SELECT id, name, email, age FROM students`,
// 	)
// 	if err != nil {
// 		return nil, fmt.Errorf("prepare error: %w", err)
// 	}
// 	defer stmt.Close()

// 	rows, err := stmt.Query()
// 	if err != nil {
// 		return nil, fmt.Errorf("query error: %w", err)
// 	}
// 	defer rows.Close()

// 	var students []types.Student

// 	for rows.Next() {
// 		var student types.Student

// 		err := rows.Scan(
// 			&student.Id,
// 			&student.Name,
// 			&student.Email,
// 			&student.Age,
// 		)
// 		if err != nil {
// 			return nil, fmt.Errorf("scan error: %w", err)
// 		}

// 		students = append(students, student)
// 	}

// 	return students, nil
// }