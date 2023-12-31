package oracle

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"io"
	"log"
	"strconv"

	"test/models"

	"github.com/godror/godror"
	"github.com/pkg/errors"
)

func (s *Store) AddUser(user models.UserSignUpRequest) error {
	// Prepare the SQL statement with the stored procedure call
	var retCode int
	stmt, err := s.DB.Prepare(`BEGIN "EV_SIGNUP"(:firstname, :lastname, :email, :password, :retcode); END;`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the stored procedure with the provided parameters
	_, err = stmt.Exec(
		sql.Named("firstname", user.FirstName),
		sql.Named("lastname", user.LastName),
		sql.Named("email", user.Email),
		sql.Named("password", user.HashPassword),
		sql.Named("retcode", sql.Out{Dest: &retCode}),
	)
	if err != nil {
		fmt.Println(err)
		return err
	}

	log.Println("User inserted successfully", retCode)

	return nil
}

// GetUserByEmail retrieves user information based on the email using EV_GETUSERS stored procedure.
func (s *Store) GetUserByEmail(email string) (models.User, error) {
	var user models.User

	// Prepare the SQL statement with the stored procedure call
	stmt, err := s.DB.Prepare(`BEGIN EV_GETUSERS(:p_EMAIL, :p_RESULT, :p_RETURNCODE); END;`)
	if err != nil {
		return user, err
	}
	defer stmt.Close()

	// Execute the stored procedure with the provided parameters
	var result driver.Rows
	var retCode int

	_, err = stmt.Exec(
		sql.Named("p_EMAIL", email),
		sql.Named("p_RESULT", sql.Out{Dest: &result}),
		sql.Named("p_RETURNCODE", sql.Out{Dest: &retCode}),
	)
	if err != nil {
		return user, err
	}
	defer result.Close()

	columns := result.(driver.RowsColumnTypeScanType).Columns()

	// Create a slice to hold the destination values for scanning
	for {
		dest := make([]driver.Value, len(columns))

		// Scan the values from the row into the dest slice
		if err := result.Next(dest); err != nil {
			if err == io.EOF {
				break
			}
			if e := result.Close(); e != nil {
				fmt.Println("Error reading the file")
			}
			return user, err
		}
		id, err := strconv.Atoi(dest[0].(godror.Number).String())
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(id)
		// Populate the user model with the retrieved data
		user = models.User{
			Id:        id,
			FirstName: dest[1].(string),
			LastName:  dest[2].(string),
			Email:     dest[3].(string), // Assuming email is the first column, adjust accordingly
			// Add other fields as needed
		}
	}

	log.Printf("User retrieved successfully: %v\n", user)

	return user, nil
}

// AuthenticateUser authenticates against the Oracle DB by calling EV_LOGIN stored procedure.
// If the stored procedure returns 3001, then authentication is successful.
// Otherwise, it considers the authentication to have failed.
func (s *Store) AuthenticateUser(loginDetail models.UserSignInRequest) (bool, error) {
	// Prepare the SQL statement with the stored procedure call
	stmt, err := s.DB.Prepare(`BEGIN EV_LOGIN(:p_EMAIL, :p_PASSWORDHASH, :p_RETURNCODE); END;`)
	if err != nil {
		return false, errors.Wrap(err, "failed to prepare statement")
	}
	defer stmt.Close()

	// Execute the stored procedure with the provided parameters
	var retCode int
	_, err = stmt.Exec(
		sql.Named("p_EMAIL", loginDetail.Email),
		sql.Named("p_PASSWORDHASH", loginDetail.HashPassword),
		sql.Named("p_RETURNCODE", sql.Out{Dest: &retCode}),
	)
	if err != nil {
		return false, errors.Wrap(err, "failed to execute statement")
	}

	// Check the return code
	if retCode == 3001 {
		// Authentication successful
		return true, nil
	}

	// Authentication failed
	return false, nil
}
