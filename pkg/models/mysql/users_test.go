package mysql

import (
	"reflect"
	"testing"
	"time"

	"github.com/snippetbox/pkg/models"
)

func TestUserModelGet(t *testing.T) {
	//

	if testing.Short() {
		t.Skip("mysql:skipping integration test")
	}

	// Set up a suite of table-driven tests and expected results
	// for different results.
	//
	tests := []struct {
		name      string
		userID    int
		wantUser  *models.User
		wantError error
	}{
		{ // Return valid result
			name:   "Valid ID",
			userID: 1,
			wantUser: &models.User{
				ID:      1,
				Name:    "Alice Jones",
				Email:   "alice@example.com",
				Created: time.Date(2018, 12, 23, 17, 25, 22, 0, time.UTC),
				Active:  true,
			},
			wantError: nil,
		},
		{ // Zero ID get NoRecord result
			name:      "Zero ID",
			userID:    0,
			wantUser:  nil,
			wantError: models.ErrNoRecord,
		},
		{ // any noexist ID get NoRecord result
			name:      "Non-existent ID",
			userID:    2,
			wantUser:  nil,
			wantError: models.ErrNoRecord,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initial a connection pool to our test database.
			// And defer a teardown function
			db, teardown := newTestDB(t)
			defer teardown()

			// Create a new instance of the UserModel
			m := UserModel{DB: db}

			//Start check userID exist or not by calling UserModel.Get()
			user, err := m.Get(tt.userID)
			if err != tt.wantError {
				t.Errorf("want %v; got %s", tt.wantError, err)
			}

			if !reflect.DeepEqual(user, tt.wantUser) {
				t.Errorf("want %v; got %v", tt.wantUser, user)
			}

		})
	}
}
