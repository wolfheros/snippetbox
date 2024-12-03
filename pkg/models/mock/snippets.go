package mock

import (
	"time"

	"github.com/snippetbox/pkg/models"
)

// create a mock database record.
var mockSnippet = &models.Snippet{
	ID:      1,
	Title:   "An old silent pond",
	Content: "An old silent pond...",
	Created: time.Now(),
	Expires: time.Now(),
}

type SnippetModel struct{}

// mock Insert on database
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	return 2, nil
}

// mock Get a record from database
func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	switch id {
	case 1:
		return mockSnippet, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	return []*models.Snippet{mockSnippet}, nil
}
