package psql

import (
	"database/sql"
	"strings"

	"github.com/lib/pq"
	"github.com/ninomaj/snippetbox/pkg/models"
	"golang.org/x/crypto/bcrypt"
)

// UserModel type which wraps a sql.DB connection pool.
type UserModel struct {
	DB *sql.DB
}

// Insert method adds a new record to the users table.
func (m *UserModel) Insert(name, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, created)
    VALUES($1, $2, $3, NOW())`

	// Use the Exec() method to insert the user details and hashed password
	// into the users table. If this returns an error, we try to type assert
	// it to a *mysql.MySQLError object so we can check if the error number is
	// 1062 and, if it is, we also check whether or not the error relates to
	// our users_uc_email key by checking the contents of the message string.
	// If it does, we return an ErrDuplicateEmail error. Otherwise, we just
	// return the original error (or nil if everything worked).
	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		if psqlErr, ok := err.(*pq.Error); ok {
			if psqlErr.Code.Name() == "unique_violation" && strings.Contains(psqlErr.Detail, "Key (email)") {
				return models.ErrDuplicateEmail
			}
		}
	}
	return err
}

// Authenticate method verifies whether a user exists with
// the provided email address and password. This will return the relevant
// user ID if they do.
func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

// Get method fetchs details for a specific user based
// on their user ID.
func (m *UserModel) Get(id int) (*models.User, error) {
	return nil, nil
}
