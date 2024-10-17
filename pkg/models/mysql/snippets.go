package mysql

import (
	"database/sql"
	"github.com/snippetbox/pkg/models" // internal package - domainname/projectname/filepath/packagename
)

// wraps sql.DB conn
type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
		VALUES(?,?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
	// receive the result from DB exc method.
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

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	// build SQL statment
	stmt := `SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() AND id = ?`
	// query the result
	row := m.DB.QueryRow(stmt, id)
	// initialize a pointer to a new zeroed Snippet struct
	s := &models.Snippet{}

	// covert sql result to Snippet s
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

	// check, if there is err with no record.
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, models.ErrNoRecord
		default:
			return nil, err
		}

	}

	return s, nil
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets 
		WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)

	if err != nil {
		return nil, err
	}

	// because rows's resultset as return value, it force the database connection open,
	// close it.
	defer rows.Close()

	// prepare a empty Snippet slice pointer
	snippets := []*models.Snippet{}

	for rows.Next() {
		s := &models.Snippet{}
		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		// append() return sometime is a new slice
		snippets = append(snippets, s)
	}
	// in case err, during the iterator process.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
