package psql

import (
	"database/sql"

	"github.com/ninomaj/snippetbox/pkg/models"
)

// SnippetModel type which wraps a sql.DB connection pool.
type SnippetModel struct {
	DB *sql.DB
}

// Insert will insert a new snippet into the database.
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
		VALUES (?, ?, NOW(),  NOW() + 365 * INTERVAL '1 day'
		RETURNING id`

	var snippetID int
	err := m.DB.QueryRow(stmt, title, content, expires).Scan(&snippetID)
	if err != nil {
		return 0, nil
	}

	return snippetID, nil
}

// Get will return a specific snippet based on its id.
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	return nil, nil
}
