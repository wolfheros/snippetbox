package mysql

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/snippetbox/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	DB *sql.DB
}

// Insert func add new record to User table
func (m *UserModel) Insert(name, email, password string) error {

	// Generate hash code from password, 12 means cost to decode it, 2^12 = 4096, recommend.
	// it will return a 60 long hash characters.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	stmt := `INSERT INTO users (name, email, hashed_password, created) VALUES(?, ?, ?, UTC_TIMESTAMP())`

	// insert data to database
	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		// any error happened, make sure whether it is the mysql error and it is duple-email error.
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return models.ErrDuplicateEmail
			}
		}
		// if its not the err we want, then return
		return err
	}
	return nil
}

// verfy whether user exist with provide email and password
func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte

	stmt := "SELECT id, hashed_password FROM users WHERE email = ? AND active= TRUE"
	row := m.DB.QueryRow(stmt, email)
	err := row.Scan(&id, &hashedPassword)

	// using email check if the record exist.
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}

	// check the record has the same password
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	// here return id as session data store for future request, same funtionality like JWT.
	return id, nil
}

// use Get method to fetch details for user based on user ID
func (m *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}
