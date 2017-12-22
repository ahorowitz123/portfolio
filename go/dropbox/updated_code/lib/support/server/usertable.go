package server

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

/**
 * Returns amount of users with this username
 *
 * Input:
 *  username string - the username to be checked
 * Output: amount of users with this username, or database error
 *
 * Precondition: username and corresponding session token have been validated and sanitized
 * Postcondition: N/A
 */

func getNumUsersWithUsername(username string) (int, error) {
	// open a connection to the database
	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return -1, err
	}
	defer db.Close()

	// prepare statement for determining number of users
	stmt, err := db.Prepare("SELECT COUNT(*) FROM users WHERE username = ?")
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	// query the db to determine number of users
	var numUsers int
	err = stmt.QueryRow(username).Scan(&numUsers) // it will only return 1 row cause count
	if err != nil {
		return -1, err
	}

	return numUsers, nil
}

/**
 * Returns amount of users with this email
 *
 * Input:
 *  email string - the email address to be checked
 * Output: amount of users with this email, or database error
 *
 * Precondition: email and corresponding session token have been validated and sanitized
 * Postcondition: N/A
 */

func getNumUsersWithEmail(email string) (int, error) {
	// open a connection to the database
	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return -1, err
	}
	defer db.Close()

	// prepare statement for determining number of users with the given email
	stmt, err := db.Prepare("SELECT COUNT(*) FROM users WHERE email = ?")
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	// query the db to determine number of users with this email
	var numUsers int
	err = stmt.QueryRow(email).Scan(&numUsers) // it will only return 1 row cause count
	if err != nil {
		return -1, err
	}

	return numUsers, nil
}

/**
 * Returns amount of users with this email
 *
 * Input:
 *  email string - the email address to be checked
 * Output: amount of users with this email, or database error
 *
 * Precondition: email and corresponding session token have been validated and sanitized
 * Postcondition: N/A
 */
func getNumUsersWithUsernameAndEmail(username string, email string) (int, error) {
	// open a connection to the database
	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return -1, err
	}
	defer db.Close()

	// prepare statement for determining number of users with the given username and email
	stmt, err := db.Prepare("SELECT COUNT(*) FROM users WHERE username = ? AND email = ?")
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	// query the db to determine number of users with this email
	var numUsers int
	err = stmt.QueryRow(username, email).Scan(&numUsers) // it will only return 1 row cause count
	if err != nil {
		return -1, err
	}

	return numUsers, nil
}

/**
 * Create user in the database before they have verified their email. Set verification code timeout as well
 * It will be removed if they don't give the right verification code in a later function
 *
 * Input:
 *  email string - the email the user entered
 *  username string - the username the user entered
 * 	password []byte - the hash of the password the user entered
 *  salt string - the salt used to hash the password for the user entered
 * 	verification_code []byte - the hash of the verification code the user must enter
 * Output: if there is already a user with this username or email, or if the db fails, return error, else nil
 *
 * Precondition: email, username, and password have all been sanitized
 * Postcondition: N/A
 */

func createNewUserPreEmail(email string, username string, password []byte, salt string, verification_code []byte) error {
	// check username
	numUsers, err := getNumUsersWithUsername(username)
	if err != nil {
		return err
	}
	if numUsers != 0 {
		return errors.New("Invalid username")
	}

	// check email
	numUsers, err = getNumUsersWithEmail(email)
	if err != nil {
		return err
	}
	if numUsers != 0 {
		return errors.New("Invalid email")
	}

	// open a connection to the database
	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()

	// once verified, insert the new values into the table
	// row = username, email, ver_code, ver_code_creation_time, password, salt, session_token
	stmt, err := db.Prepare(`INSERT INTO users 
		(username, email, ver_code, ver_code_creation_time, password, salt)
		VALUES (?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// and execute the insert
	curr_time := time.Now().Unix()
	_, err = stmt.Exec(username, email, verification_code, curr_time, password, salt)
	if err != nil {
		return err
	}

	return nil
}

/**
 * Returns if the username and email given by the client is unique
 *
 * Input:
 *	username string - the username to be checked
 *  email string - the email address to be checked
 * Output: if there is a user with this username or email return error, else nil
 *
 * Precondition: username, email and corresponding session token have been validated and sanitized
 * Postcondition: N/A
 */

func validateUniqueUserByUsernameAndEmail(username string, email string) error {
	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()

	numUsers, err := getNumUsersWithUsername(username)
	if err != nil {
		return err
	}
	if numUsers != 1 {
		return errors.New("Invalid username")
	}

	numUsers, err = getNumUsersWithEmail(email)
	if err != nil {
		return err
	}
	if numUsers != 1 {
		return errors.New("Invalid email")
	}

	// now, check that username/password combo belong to only one user
	numUsers, err = getNumUsersWithUsernameAndEmail(username, email)
	if err != nil {
		return err
	}
	if numUsers != 1 {
		return errors.New("Invalid username/email")
	}

	return nil
}

/**
 * Returns if the username given by the client is unique
 *
 * Input:
 *	username string - the username to be checked
 * Output: if there is a user with this username return error, else nil
 *
 * Precondition: username and corresponding session token have been validated and sanitized
 * Postcondition: N/A
 */

func validateUniqueUserByUsername(username string) error {
	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()

	numUsers, err := getNumUsersWithUsername(username)
	if err != nil {
		return err
	}
	if numUsers != 1 {
		fmt.Println("NUMUSERS:", numUsers)
		return errors.New("Invalid username")
	}

	return nil
}

/**
 * Returns the hashed version of a users password given the username and email
 *
 * Input:
 *	username string - the username to be checked
 *  email string - the email address to be checked
 * Output: if the username and email is unique then return the hashed password, else error
 *
 * Precondition: username, email and corresponding session token have been validated and sanitized
 * Postcondition: N/A
 */

func getUserPasswordByUsernameAndEmail(username string, email string) ([]byte, error) {
	// validate only 1 user with this username/email combo
	err := validateUniqueUserByUsernameAndEmail(username, email)
	if err != nil {
		return []byte{}, err
	}

	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return []byte{}, err
	}
	defer db.Close()

	// get the password for this user
	stmt, err := db.Prepare("SELECT password FROM users WHERE username = ? AND email = ?")
	if err != nil {
		return []byte{}, err
	}
	defer stmt.Close()

	// query the db to get password of user
	var password []byte
	err = stmt.QueryRow(username, email).Scan(&password) // it will only return 1 row cause unique check
	// TODO: Does stmt.QueryRow close itself? idk
	if err != nil {
		return []byte{}, err
	}

	return password, nil
}

/**
 * Returns the hashed version of a users password given the username
 *
 * Input:
 *	username string - the username to be checked
 * Output: if the username is unique then return the hashed password, else error
 *
 * Precondition: username, email and corresponding session token have been validated and sanitized
 * Postcondition: N/A
 */

func getUserPasswordByUsername(username string) ([]byte, error) {
	// validate only 1 user with this username
	err := validateUniqueUserByUsername(username)
	if err != nil {
		return []byte{}, err
	}

	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return []byte{}, err
	}
	defer db.Close()

	// get the password for this user
	stmt, err := db.Prepare("SELECT password FROM users WHERE username = ?")
	if err != nil {
		return []byte{}, err
	}
	defer stmt.Close()

	// query the db to get password of user
	var password []byte
	err = stmt.QueryRow(username).Scan(&password) // it will only return 1 row cause unique check
	// TODO: Does stmt.QueryRow close itself? idk
	if err != nil {
		return []byte{}, err
	}

	return password, nil
}

/**
 * Returns the hashed version of a users verification code
 *
 * Input:
 *	username string - the username to be checked
 *  email string - the username to be checked
 * Output: if the username and email is unique then return the varification code and the verification timeout time, else error
 *
 * Precondition: username, email and corresponding session token have been validated and sanitized
 * Postcondition: N/A
 */

func getVerCodeInfoByUsernameAndEmail(username string, email string) ([]byte, int, error) {
	// validate only 1 user with this username/email combo
	err := validateUniqueUserByUsernameAndEmail(username, email)
	if err != nil {
		return []byte{}, -1, err
	}

	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return []byte{}, -1, err
	}
	defer db.Close()

	// get the ver_code and ver_code_creation_time for this user
	stmt, err := db.Prepare(`SELECT ver_code, ver_code_creation_time FROM users 
		WHERE username = ? AND email = ?`)
	if err != nil {
		return []byte{}, -1, err
	}
	defer stmt.Close()

	// query db to get the user info
	var ver_code []byte
	var ver_code_creation_time int
	err = stmt.QueryRow(username, email).Scan(&ver_code, &ver_code_creation_time)
	if err != nil {
		return []byte{}, -1, err
	}

	return ver_code, ver_code_creation_time, nil
}

/**
 * Removes a users verification code and verification code timeout from the database
 *
 * Input:
 *	username string - the username to be checked
 *  email string - the username to be checked
 * Output: if the username and email is unique then remove the verification code and timeout from the database, else error
 *
 * Precondition: username, email and corresponding session token have been validated and sanitized
 * Postcondition: Verification code and time should be removed from the users table
 */

func removeVerCodeInfoByUsernameAndEmail(username string, email string) error {
	// validate only 1 user with this username/email combo
	err := validateUniqueUserByUsernameAndEmail(username, email)
	if err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()

	// remove ver_code and ver_code_creation_time for this user
	stmt, err := db.Prepare(`UPDATE users
		SET ver_code = NULL, ver_code_creation_time = NULL
		WHERE username = ? AND email = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, email)
	if err != nil {
		return err
	}

	return nil
}

/**
 * Returns the salt of a user
 *
 * Input:
 *	username string - the username to be checked
 * Output: if the username is unique then return the corresponding salt, else error
 *
 * Precondition: username and corresponding session token have been validated and sanitized
 * Postcondition: N/A
 */

func getSaltByUsername(username string) (string, error) {
	// validate only 1 user with this username/email combo
	err := validateUniqueUserByUsername(username)
	if err != nil {
		return "", err
	}

	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return "", err
	}
	defer db.Close()

	// get the salt for this user
	stmt, err := db.Prepare(`SELECT salt FROM users 
		WHERE username = ?`)
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	// query db to get the user info
	var salt string
	err = stmt.QueryRow(username).Scan(&salt)
	if err != nil {
		return "", err
	}

	return salt, nil
}

/**
 * Removes the given user from the users table
 *
 * Input:
 *	username string - the username to be checked
 *  email string - the username to be checked
 * Output: if the username and email is unique then remove the given user from the user table, else error
 *
 * Precondition: username, email and corresponding session token have been validated and sanitized
 * Postcondition: The given user will be removed from the users table
 */

func removeUserRowByUsernameAndEmail(username string, email string) error {
	// validate only 1 user with this username/email combo
	err := validateUniqueUserByUsernameAndEmail(username, email)
	if err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(`DELETE FROM users
		WHERE username = ? AND email = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, email)
	if err != nil {
		return err
	}

	return nil
}

/**
 * Removes the given user from the users table
 *
 * Input:
 *	username string - the username to be checked
 * Output: if the username is unique then remove the given user from the user table, else error
 *
 * Precondition: username and corresponding session token have been validated and sanitized
 * Postcondition: The given user will be removed from the users table
 */

func removeUserRowByUsername(username string) error {
	// validate only 1 user with this username/email combo
	fmt.Printf("HERE!\n")
	err := validateUniqueUserByUsername(username)
	if err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(`DELETE FROM users
		WHERE username = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	fmt.Printf("Execing...\n")
	_, err = stmt.Exec(username)
	if err != nil {
		return err
	}
	fmt.Printf("Done, returning nil...\n")

	return nil
}

/**
 * Puts the given session token in the database to correspond with the given user. Also sets timeout for the session
 *
 * Input:
 *	username string - the username to be checked
 *  session_token []byte - the username to be checked
 * Output: if the username is unique then return nil, else error
 *
 * Precondition: username has been validated and sanitized
 * Postcondition: The given user will have a session token with a timeout in the users table
 */

func createSessionTokenByUsername(username string, session_token []byte) error {
	// validate only 1 user with this username
	err := validateUniqueUserByUsername(username)
	if err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()

	// set session_token and session_token_last_access_time for this user
	stmt, err := db.Prepare(`UPDATE users
		SET session_token = ?, session_token_last_access_time = ?
		WHERE username = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	curr_time := time.Now().Unix() // session token creation time
	_, err = stmt.Exec(session_token, curr_time, username)
	if err != nil {
		return err
	}

	return nil

}

/**
 * Removes the session token for the given user from the database. Also removes timeout for the session
 *
 * Input:
 *	username string - the username to be checked
 * Output: if the username is unique then return nil, else error
 *
 * Precondition: username and corresponding session token have been validated and sanitized
 * Postcondition: The given user will no longer have a session token and corresponding timeout
 */

func removeSessionTokenInfoByUsername(username string) error {
	// validate only 1 user with this username
	err := validateUniqueUserByUsername(username)
	if err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()

	// remove session_token and session_token_last_access_time for this user
	stmt, err := db.Prepare(`UPDATE users
		SET session_token = NULL, session_token_last_access_time = NULL
		WHERE username = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(username)
	if err != nil {
		return err
	}

	return nil
}

/**
 * Sets the current directory of the user. Used for cd, pwd, etc.
 *
 * Input:
 *	username string - the username to be checked
 *	path string - the path to be set for the given user
 * Output: if the username is unique and the path is a valid existing direcotry then nil, else error
 *
 * Precondition: username and corresponding session token have been validated and sanitized
 * Postcondition: The cuurent directory for that user is set in the users table
 */

func setCurrDir(username string, path string) error {
	//add backslash at end of path if they don't have it
	if path[len(path)-1:] != "/" {
		path = path + "/"
	}

	// validate only 1 user with this username
	err := validateUniqueUserByUsername(username)
	if err != nil {
		return err
	}

	path, err = checkValidDirectoryPath(username, path)
	if err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(`UPDATE users
		SET curr_dir = ? WHERE username = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(path, username)
	if err != nil {
		return err
	}

	return nil
}

/**
 * Gets the current directory of the user. Used for cd, pwd, etc.
 *
 * Input:
 *	username string - the username to be checked
 * Output: if the username is unique then it returns a string of the current directory for the user, else error
 *
 * Precondition: username and corresponding session token have been validated and sanitized
 * Postcondition: N/A
 */

func getCurrDir(username string) (string, error) {
	err := validateUniqueUserByUsername(username)
	if err != nil {
		return "", err
	}

	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return "", err
	}
	defer db.Close()

	stmt, err := db.Prepare(`SELECT curr_dir FROM users WHERE username = ?`)
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	var curr_dir string
	err = stmt.QueryRow(username).Scan(&curr_dir)
	if err != nil {
		return "", err
	}

	return curr_dir, nil
}

/**
 * Resets the current directory of the given user to the root (/)
 *
 * Input:
 *	username string - the username to be checked
 * Output: if the username is unique then nil, else error
 *
 * Precondition: username and corresponding session token have been validated and sanitized
 * Postcondition: The cuurent directory for that user is set to the root
 */

func resetCurrDir(username string) error {
	// validate only 1 user with this username
	err := validateUniqueUserByUsername(username)
	if err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(`UPDATE users
		SET curr_dir = ? WHERE username = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec("", username)
	if err != nil {
		return err
	}

	return nil
}

/**
 * Gets the session token and timeout time of the user.
 *
 * Input:
 *	username string - the username to be checked
 * Output: if the username is unique then it returns the session token and timeout for the user, else error
 *
 * Precondition: username and corresponding session token have been validated and sanitized
 * Postcondition: N/A
 */

func getSessionTokenInfoByUsername(username string) ([]byte, int, error) {
	// validate unique user
	err := validateUniqueUserByUsername(username)
	if err != nil {
		return []byte{}, -1, err
	}

	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return []byte{}, -1, err
	}
	defer db.Close()

	// prepare the query
	stmt, err := db.Prepare(`SELECT session_token, session_token_last_access_time
	 FROM users
	 WHERE username = ?`)
	if err != nil {
		return []byte{}, -1, err
	}
	defer stmt.Close()

	// query db to get session token info
	var session_token_hash []byte
	var session_token_last_access_time int
	err = stmt.QueryRow(username).Scan(&session_token_hash, &session_token_last_access_time)
	if err != nil {
		return []byte{}, -1, err
	}

	return session_token_hash, session_token_last_access_time, nil
}

/**
 * Resets the session timeout time of the user.
 *
 * Input:
 *	username string - the username to be checked
 * Output: if the username is unique then it returns nil, else error
 *
 * Precondition: username and corresponding session token have been validated and sanitized
 * Postcondition: The session timeout for the user is reset
 */

func updateSessionTokenInfoByUsername(username string) error {
	// validate unique user
	err := validateUniqueUserByUsername(username)
	if err != nil {
		return err
	}

	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()

	// update session_token_last_access_time for this user
	stmt, err := db.Prepare(`UPDATE users
		SET session_token_last_access_time = ?
		WHERE username = ?`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	curr_time := time.Now().Unix() // session token creation time
	_, err = stmt.Exec(curr_time, username)
	if err != nil {
		return err
	}

	return nil
}

/**
 * Check if we are over the user limit of 100.
 *
 * Input: N/A
 * Output: If we are over 100 users, an error, else nil.
 *
 * Precondition: username and session token are validated
 * Postcondition: N/A
 */

func overUsersLimit() error {
	// open a connection to the database
	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()

	// prepare statement for determining number of users
	stmt, err := db.Prepare("SELECT COUNT(*) FROM users")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// query the db to determine number of users
	var numUsers int
	err = stmt.QueryRow().Scan(&numUsers) // it will only return 1 row cause count
	if err != nil {
		return err
	}

	fmt.Println("NUMUSERS", numUsers)

	if numUsers >= 100 {
		return errors.New("No more users can be registered")
	}

	return nil
}

/**
 * Check if users has reached total bytes capacity for itself.
 * Calulated by vm size/max num of users
 *
 * Input: The username of the user we are checking
 * Output: If we there are over 28000 total bytes for username, an error, else nil.
 *
 * Precondition: username and session token are validated
 * Postcondition: N/A
 */

func overUserSizeLimit(username string) error {
	// open a connection to the database
	db, err := sql.Open("sqlite3", "boxdrop.sqlite")
	if err != nil {
		return err
	}
	defer db.Close()

	// prepare statement for determining number of users
	stmt, err := db.Prepare("SELECT checksum FROM files WHERE username = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	// query the db to determine number of users
	rows, err := stmt.Query(username) // it will only return 1 row cause count
	if err != nil {
		return err
	}

	total_bytes := int64(0)

	for rows.Next() {
		var filename string
		err := rows.Scan(&filename)
		if err != nil {
			return err
		}
		filepath := basedir + "/" + filename
		file, err := os.Open(filepath)
		if err != nil {
			return err
		}
		fi, err := file.Stat()
		if err != nil {
			return err
		}
		total_bytes += fi.Size()
	}

	fmt.Println("TOTALBYTES", total_bytes)

	if total_bytes >= 28000 {
		return errors.New("You cannot upload any more files")
	}

	return nil
}
