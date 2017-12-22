package server

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"

	"github.com/s17-ahorowi2-jtc2/updated_code/internal"
)

/**
 * Validate that the first part of registration is valid (pre email verification code).
 *
 * Input:
 *   email string - the email from the client
 *   username string - the username from the client
 *   password string - the password from the client
 * Output: A register return type indicicating success or failure.
 *
 * Precondition: N/A
 * Postcondition: If successful, email, username, and password are all valid by the server check. The email and
 *   username have never been used before. An email has been sent to the user with the verification code.
 */
func emailRegisterHandler(email string, username string, password string) internal.RegisterReturn {
	// ensure can create a new user
	err := overUsersLimit()
	if err != nil {
		return internal.RegisterReturn{Success: false, Err: err.Error()}
	}

	email = strings.TrimSpace(email)
	// server side validate user email
	if !validateEmail(email) {
		return internal.RegisterReturn{Success: false, Err: "Invalid email address, canceling registration."}
	}

	username = strings.TrimSpace(username)
	// server side validate username
	if !validateUsername(username) {
		return internal.RegisterReturn{Success: false, Err: "Invalid username, canceling registration."}
	}

	//server side validate password
	passed, err := validatePassword(password)
	if !passed {
		return internal.RegisterReturn{Success: false, Err: err.Error()}
	}

	// make and salt the password
	salt := createRandomString(8)
	salted_pass := password + salt
	byte_pass := []byte(salted_pass)
	hashed_pass, err := bcrypt.GenerateFromPassword(byte_pass, bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return internal.RegisterReturn{Success: false, Err: err.Error()}
	}

	// create a verification code, hash it for storing in the database
	ver := createRandomString(6)
	ver_bytes := []byte(ver)

	hashed_ver, err := bcrypt.GenerateFromPassword(ver_bytes, bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return internal.RegisterReturn{Success: false, Err: err.Error()}
	}

	err = createNewUserPreEmail(email, username, hashed_pass, salt, hashed_ver)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return internal.RegisterReturn{Success: false, Err: err.Error()}
	}
	// send unencrypted verification code to the user.
	sendEmail(email, ver)

	return internal.RegisterReturn{Success: true}

}

/**
 * Validate that the econd part of registration is valid (meaning the email verification code is valid).
 *
 * Input:
 *   email string - the email from the client
 *   username string - the username from the client
 *   password string - the password from the client
 *   ver_code string - the client's entry for the email verification code
 * Output: A register return type indicicating success or failure.
 *
 * Precondition: The email and username are in the database, and the user has a verification code but not a
 *   session token.
 * Postcondition: If successful, email, username, and password are all valid by the server check. The given
 *   verification code matches the one in the database. The account is successfully created.
 *   Failures at any point will remove the row from the database for this username and email.
 */
func validateRegisterHandler(email string, username string, password string, ver_code string) internal.RegisterReturn {
	email = strings.TrimSpace(email)
	username = strings.TrimSpace(username)
	password = strings.TrimSpace(password)
	ver_code = strings.TrimSpace(ver_code)

	email = strings.TrimSpace(email)
	// server side validate user email
	if !validateEmail(email) {
		return internal.RegisterReturn{Success: false, Err: "Invalid email address, canceling registration."}
	}

	username = strings.TrimSpace(username)
	// server side validate username
	if !validateUsername(username) {
		return internal.RegisterReturn{Success: false, Err: "Invalid username, canceling registration."}
	}

	//server side validate password
	passed, err := validatePassword(password)
	if !passed {
		return internal.RegisterReturn{Success: false, Err: err.Error()}
	}

	// ensure given password is correct
	stored_password, err := getUserPasswordByUsernameAndEmail(username, email)
	if err != nil {
		removerErr := removeUserRowByUsernameAndEmail(username, email)
		if removerErr != nil {
			return internal.RegisterReturn{Success: false, Err: removerErr.Error()}
		}
		return internal.RegisterReturn{Success: false, Err: err.Error()}
	}

	salted_password, err := getSaltedPass(username, password)
	if err != nil {
		return internal.RegisterReturn{Success: false, Err: err.Error()}
	}

	err = bcrypt.CompareHashAndPassword(stored_password, salted_password)
	if err != nil {
		return internal.RegisterReturn{Success: false, Err: "Incorrect password."}
	}

	// Get and check verification code along with its creation time
	stored_ver_code, stored_ver_code_creation_time, err := getVerCodeInfoByUsernameAndEmail(username, email)
	if err != nil {
		removerErr := removeUserRowByUsernameAndEmail(username, email)
		if removerErr != nil {
			return internal.RegisterReturn{Success: false, Err: removerErr.Error()}
		}
		return internal.RegisterReturn{Success: false, Err: err.Error()}
	}

	// before any verification, remove ver_code and its creation time
	err = removeVerCodeInfoByUsernameAndEmail(username, email)
	if err != nil {
		removerErr := removeUserRowByUsernameAndEmail(username, email)
		if removerErr != nil {
			return internal.RegisterReturn{Success: false, Err: removerErr.Error()}
		}
		return internal.RegisterReturn{Success: false, Err: err.Error()}
	}

	// check that verification code not expired
	cur_time := int(time.Now().Unix())
	time_diff := 60 * 5 // 5 minutes
	if stored_ver_code_creation_time+time_diff < cur_time {
		removerErr := removeUserRowByUsernameAndEmail(username, email)
		if removerErr != nil {
			return internal.RegisterReturn{Success: false, Err: removerErr.Error()}
		}
		return internal.RegisterReturn{Success: false, Err: "Verification code expired"}
	}

	ver_code_bytes := []byte(ver_code)

	// check that verification code is correct
	err = bcrypt.CompareHashAndPassword(stored_ver_code, ver_code_bytes)
	if err != nil {
		removerErr := removeUserRowByUsernameAndEmail(username, email)
		if removerErr != nil {
			return internal.RegisterReturn{Success: false, Err: removerErr.Error()}
		}
		return internal.RegisterReturn{Success: false, Err: "Incorrect password."}
	}

	// new user, give them root directory and shared folder
	err = createRootandShared(username)
	if err != nil {
		removerErr := removeUserRowByUsernameAndEmail(username, email)
		if removerErr != nil {
			return internal.RegisterReturn{Success: false, Err: removerErr.Error()}
		}
		return internal.RegisterReturn{Success: false, Err: err.Error()}
	}

	return internal.RegisterReturn{Success: true}
}

/**
 * Validate that the login is valid.
 *
 * Input:
 *   username string - the username from the client
 *   password string - the password from the client
 * Output: A login return type indicicating success or failure. Success will return a session token.
 *
 * Precondition: N/A
 * Postcondition: If successful, username and password are all valid by the server check. The database has the
 *   username, along with a salted version of the password. A session token is generated for the user and
 *   returned to the user for use in future requests.
 *   Failures will not generate a session token.
 */
func validateLoginHandler(username string, password string) internal.LoginReturn {
	// validate username and password
	username = strings.TrimSpace(username)
	// server side validate username
	if !validateUsername(username) {
		return internal.LoginReturn{SessionToken: []byte{}, Err: "Invalid username, canceling registration."}
	}

	//server side validate password
	passed, err := validatePassword(password)
	if !passed {
		return internal.LoginReturn{SessionToken: []byte{}, Err: err.Error()}
	}

	// ensure password is right
	stored_password, err := getUserPasswordByUsername(username)
	if err != nil {
		return internal.LoginReturn{SessionToken: []byte{}, Err: err.Error()}
	}

	salted_password, err := getSaltedPass(username, password)
	if err != nil {
		return internal.LoginReturn{SessionToken: []byte{}, Err: err.Error()}
	}

	err = bcrypt.CompareHashAndPassword(stored_password, salted_password)
	if err != nil {
		return internal.LoginReturn{SessionToken: []byte{}, Err: "Incorrect password."}
	}

	// if valid password, generate session token
	session_token := createRandomString(64)
	session_token_bytes := []byte(session_token)
	session_token_hash, err := bcrypt.GenerateFromPassword(session_token_bytes, bcrypt.DefaultCost)
	if err != nil {
		return internal.LoginReturn{SessionToken: []byte{}, Err: err.Error()}
	}

	// Put session token into database, along with its creation time
	err = createSessionTokenByUsername(username, session_token_hash)
	if err != nil {
		return internal.LoginReturn{SessionToken: []byte{}, Err: err.Error()}
	}

	//Set current direcotry to root
	err = setCurrDir(username, "/")
	if err != nil {
		return internal.LoginReturn{SessionToken: []byte{}, Err: err.Error()}
	}

	return internal.LoginReturn{SessionToken: session_token_bytes, Err: ""}
}
