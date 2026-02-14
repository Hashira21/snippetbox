package models

import (
	"database/sql"
	"errors"
	"time"
)

// Объевим тип Snippet для хранения данных по каждому снипету
// Поля структу совпадают с столбцами из БД
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// Объявим тип SnippetModel, который содержит пул подключений sql.DB
type SnippetModel struct {
	DB *sql.DB
}

// This will insert a new snippet into the database
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil

}

// This will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() AND id = ?`

	row := m.DB.QueryRow(stmt, id)

	var s Snippet

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord // из models/errors.go
		} else {
			return Snippet{}, err
		}
	}

	return s, nil
}

// This will return the 10 most recently created snippets.
func (m *SnippetModel) Latest() ([]Snippet, error) {

	stmt := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	// возвращает набор результатов sql.Rows
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	// Мы откладываем вызов rows.Close(), чтобы убедиться,
	// что набор результатов sql.Rows корректно закрывается
	// перед возвратом из метода Latest().
	defer rows.Close()

	// Слайс для хранения сниппетов
	var snippets []Snippet

	// Проходимся по строкам из набор результатов sql.Rows.
	// Если итерация по всем строкам завершена, то
	// набор результатов автоматически закрывается и освобождает базовое соединение с базой данных
	for rows.Next() {
		var s Snippet
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
