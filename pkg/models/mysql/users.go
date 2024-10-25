package mysql

import (
	"database/sql"

	"github.com/snippetbox/pkg/models"
)

type UserModel struct{
	DB *sql.DB
}

// Insert func add new record to User table
func (m *UserModel) Insert (name, email, password string) error{
	return nil
}

// verfy whether user exist with provide email and password
func (m *UserModel) Authenticate(email, password string) (int, error){
	return 0, nil
}

// use Get method to fetch details for user based on user ID
func (m *UserModel) Get(id int) (*models.User, error){
	return nil, nil
}

