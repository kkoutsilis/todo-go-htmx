package database

import (
	"database/sql"
	"errors"
	"todo-go-htmx/models"

	"github.com/mattn/go-sqlite3"
)

var (
	ErrDuplicate    = errors.New("record already exists")
	ErrNotExists    = errors.New("row not exists")
	ErrUpdateFailed = errors.New("update failed")
	ErrDeleteFailed = errors.New("delete failed")
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{
		db: db,
	}
}

func (r *SQLiteRepository) Migrate() error {
	query := `
    CREATE TABLE IF NOT EXISTS todos(
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        description TEXT NOT NULL UNIQUE,
        completed BOOLEAN NOT NULL DEFAULT FALSE
    );
    `
	_, err := r.db.Exec(query)
	return err
}

func (r *SQLiteRepository) Create(todo models.Todo) (*models.Todo, error) {
	res, err := r.db.Exec("INSERT INTO todos(description, completed) values(?,?)", todo.Description, todo.Completed)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
				return nil, ErrDuplicate
			}
		}
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	todo.Id = id

	return &todo, nil

}

func (r *SQLiteRepository) All() ([]models.Todo, error) {
	rows, err := r.db.Query("SELECT * FROM todos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var all []models.Todo
	for rows.Next() {
		var todo models.Todo
		if err := rows.Scan(&todo.Id, &todo.Description, &todo.Completed); err != nil {
			return nil, err
		}
		all = append(all, todo)
	}
	return all, nil
}

func (r *SQLiteRepository) GetById(id int64) (*models.Todo, error) {
	row := r.db.QueryRow("SELECT * FROM todos WHERE id = ?", id)

	var todo models.Todo
	if err := row.Scan(&todo.Id, &todo.Description, &todo.Completed); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotExists
		}
		return nil, err
	}
	return &todo, nil
}

func (r *SQLiteRepository) Update(id int64, updated models.Todo) (*models.Todo, error) {
	if id == 0 {
		return nil, errors.New("invalid id")
	}
	res, err := r.db.Exec("UPDATE todos SET description = ?, completed = ? WHERE id = ?", updated.Description, updated.Completed, id)
	if err != nil {
		return nil, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rowsAffected == 0 {
		return nil, ErrUpdateFailed
	}

	return &updated, nil
}

func (r *SQLiteRepository) Delete(id int64) error {
	res, err := r.db.Exec("DELETE FROM todos WHERE id = ?", id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrDeleteFailed
	}

	return err
}
