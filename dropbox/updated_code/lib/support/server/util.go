package server

import (
	"crypto/sha512"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"net/smtp"
	"regexp"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

/**
 * Checks if the email is a valid email string.
 *
 * Input:
 *   email string - the email address to be checked
 * Output: true if a valid email address, else false
 *
 * Precondition: N/A
 * Postcondition: N/A
 */
func validateEmail(email string) bool {
	// Regex from https://www.socketloop.com/tutorials/golang-validate-email-address-with-regular-expression
	emailRegexp := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegexp.MatchString(email)
}

/**
 * Checks if the username is a valid username string.
 *
 * Input:
 *   username string - the username to be checked
 * Output: true if a valid username, else false
 *
 * Precondition: N/A
 * Postcondition: N/A
 */
func validateUsername(username string) bool {
	// Regex from https://www.socketloop.com/tutorials/golang-regular-expression-alphanumeric-underscore
	usernameRegexp := regexp.MustCompile("^[a-zA-Z0-9]+$")
	return len(username) <= 16 && usernameRegexp.MatchString(username)
}

/**
 * Check if the passwords are equal and are a valid password.
 *
 * Input:
 *   password1 string - the first password enetered during registration
 *   password2 string - the second password entered during registration
 * Output:
 *   bool - true if passwords equal and valid, else false
 *   error - error describing why passwords failed to validate, nil if valid
 *
 * Precondition: N/A
 * Postcondition: N/A
 */
func validatePassword(password string) (bool, error) {
	if len(password) < 8 || len(password) > 24 {
		return false, errors.New("Passwords must be 8-24 characters")
	}

	whitespaceRegexp := regexp.MustCompile(`[\s]+`)
	lowercaseRegexp := regexp.MustCompile(`[a-z]+`)
	uppercaseRegexp := regexp.MustCompile(`[A-Z]+`)
	numericRegexp := regexp.MustCompile(`[0-9]+`)
	punctuationRegexp := regexp.MustCompile(`[(!#$&*+,-.:;?@^_~)]+`)
	invalidRegexp := regexp.MustCompile(`[^(!#$&*+,-.:;?@^_~)0-9A-Za-z]+`)

	if whitespaceRegexp.MatchString(password) || invalidRegexp.MatchString(password) {
		return false, errors.New("Passwords have invalid characters")
	}

	counter := 0
	if lowercaseRegexp.MatchString(password) {
		counter++
	}
	if uppercaseRegexp.MatchString(password) {
		counter++
	}
	if numericRegexp.MatchString(password) {
		counter++
	}
	if punctuationRegexp.MatchString(password) {
		counter++
	}

	if counter >= 3 {
		return true, nil
	} else {
		return false, errors.New("Password did not meet constraints")
	}
}

/**
 * Check that the given dirname is a valid name for a directory.
 *
 * Input:
 *   dirname string - the name of the directory to be checked
 * Output: true if a valid dirname, else false
 *
 * Precondition: N/A
 * Postcondition: N/A
 */
func validateDirname(dirname string) bool {
	dirnameRegexp := regexp.MustCompile("^[a-zA-Z0-9]+$")
	return len(dirname) <= 32 && dirnameRegexp.MatchString(dirname)
}

/**
 * Check that the given filename is a valid name for a file.
 *
 * Input:
 *   filename string - the name of the file to be checked
 * Output: true if a valid filename, else false
 *
 * Precondition: N/A
 * Postcondition: N/A
 */
func validateFilename(filename string) bool {
	filenameRegexp := regexp.MustCompile("^[a-zA-Z0-9_\\-.]+$")
	return len(filename) <= 18 && filenameRegexp.MatchString(filename)
}

/**
 * Send an email with the given body to the given email address.
 *
 * Input:
 *   addr string - the email address to send an email to
 *   body string - the body of the email
 * Output: the error, if an error occurred, else nil
 *
 * Precondition: N/A
 * Postcondition: N/A
 */
func sendEmail(addr string, body string) error {
	from := "boxdrop162@gmail.com"
	pass := "thisshouldnotbeplaintext"
	to := addr

	msg := body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		return err
	}

	return nil
}

/**
 * Generate a random alphanumeric string of the given length, from
 * http://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
 *
 * Input:
 *   length int - the length of the string
 * Output: the random string that was created
 *
 * Precondition: N/A
 * Postcondition: N/A
 */
func createRandomString(length int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789")
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

/**
 * Return the given password concatenated with the users salt.
 *
 * Input:
 *   username string - the username of the user
 *   password string - the password of the user
 * Output:
 *   []byte - the password concatenated with the salt
 *   error - the error if one occurred
 *
 * Precondition: The user exists, already has a salt.
 * Postcondition: N/A
 */
func getSaltedPass(username string, password string) ([]byte, error) {
	salt, err := getSaltByUsername(username)
	if err != nil {
		return []byte{}, err
	}
	salted_pass := password + salt
	byte_pass := []byte(salted_pass)
	return byte_pass, nil
}

/**
 * Check if the given password matches the password stored in the database.
 *
 * Input:
 *   username string - the username of the user
 *   password string - the password of the user
 * Output:
 *   []byte - the password concatenated with the salt
 *   error - the error if one occurred
 *
 * Precondition: The user exists, already has a salt. The password is already salted if appropriate.
 * Postcondition: N/A
 */
func checkPassword(username string, password string) error {
	stored_pass, err := getUserPasswordByUsernameAndEmail(username, password)
	if err != nil {
		return errors.New("Incorrect password.")
	}
	pass_bytes := []byte(password)
	err = bcrypt.CompareHashAndPassword(stored_pass, pass_bytes)
	if err != nil {
		errors.New("Incorrect password.")
	}
	return nil
}

/**
 * Check if the given session_token matches the session_token stored in the database for the user.
 *
 * Input:
 *   username string - the username of the user
 *   session_token []byte - the session_token of the user
 * Output:
 *   error - the error if one occurred
 *
 * Precondition: The user exists, already has a session_token generated previously
 * Postcondition: N/A
 */
func checkSessionToken(username string, session_token []byte) error {
	stored_token, stored_token_last_access_time, err := getSessionTokenInfoByUsername(username)
	if err != nil {
		return err
	}

	cur_time := int(time.Now().Unix())
	time_diff := 60 * 5 // 5 minutes
	if stored_token_last_access_time+time_diff < cur_time {
		// invalid session token, remove session token info and indicate client should be logged out
		err = removeSessionTokenInfoByUsername(username)
		if err != nil {
			return err
		}
		return errors.New("Session token expired.")
	}

	// returns error if not equal, nil if equal
	err = bcrypt.CompareHashAndPassword(stored_token, session_token)
	if err != nil {
		return errors.New("Incorrect password.")
	}

	// if valid session token, udpate its last access time
	err = updateSessionTokenInfoByUsername(username)
	return err
}

/**
 * Check if path is absolute.
 *
 * Input:
 *   path string - the path to be checked
 * Output: true if absolute, else false
 *
 * Precondition: N/A
 * Postcondition: N/A
 */
func checkIfAbsolutePath(path string) bool {
	if path == "" {
		return false
	}

	if string(path[0]) == "/" {
		return true
	}

	return false
}

/**
 * Make the checksum of the given file using sha512. Note that none of these functions return errors.
 *
 * Input:
 *   file_bytes - the content of the file.
 * Output:
 *   []byte - the sha512 checksum of the file
 *
 * Precondition: N/A
 * Postcondition: N/A
 */
func makeFileChecksum(file_bytes []byte) []byte {
	sha_512 := sha512.New()
	sha_512.Write(file_bytes)
	checksum := sha_512.Sum(nil)
	return checksum
}

/**
 * Get the filename and path from the given filepath.
 *
 * Input:
 *   path string - the path to the file, including filename
 * Output:
 *   filepath string - the path to the file
 *   file string - the name of the file
 *
 * Precondition: N/A
 * Postcondition: N/A
 */
func getFileNameAndPath(path string) (filepath string, file string) {
	last_slash_index := strings.LastIndex(path, "/")
	filepath = path[:last_slash_index+1]
	file = path[last_slash_index+1:]

	return filepath, file
}

/**
 * Check if file exists and gets the file's checksum if so.
 *
 * Input:
 *   username string - the username of the user
 *   filename string - the name of the file
 *   filepath string - the path of the file
 * Output:
 *   checksum string - the checksum of the file, if the file exists
 *   err error - the error, if one occurred or if the file did not exist. Else nil
 *
 * Precondition: User is valid.
 * Postcondition: N/A
 */
func checkIfFileExistsAndgetChecksum(username string, filename string, filepath string) (checksum string, err error) {
	count, err := countNumFilesByUsernameAndFileInfo(username, filename, filepath)
	if err != nil {
		return "", err
	}

	if count != 1 {
		return "", errors.New("This file does not exist")
	}

	checksum, err = getChecksum(username, filename, filepath)
	if err != nil {
		return "", err
	}

	return checksum, nil
}

/**
 * Check if the given directory is in the user's shared directory.
 *
 * Input:
 *   curr_dir string - the user's current directory
 * Output: true if curr_dir in shared or a subdirectory of shared, else false
 *
 * Precondition: curr_dir is a valid directory path
 * Postcondition: N/A
 */
func checkIfCurrDirIsInSharedFolder(curr_dir string) bool {
	return strings.HasPrefix(curr_dir, "/shared/")
}

/**
 * Get the owner of a shared directory, meaning the sharer of the files
 *
 * Input:
 *   curr_dir string - the user's current directory
 * Output: the name of the owner
 *
 * Precondition: The curr_dir is in the shared directory and is in a subdirectory of it.
 * Postcondition: N/A
 */
func getSharedOwner(curr_dir string) string {
	dirs := strings.Split(curr_dir, "/")
	return dirs[2]
}
